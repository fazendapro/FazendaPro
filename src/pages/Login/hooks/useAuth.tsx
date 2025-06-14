import { useState, useEffect, useCallback } from 'react';
import { jwtDecode } from 'jwt-decode';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { useTranslation } from 'react-i18next';
import { LoginFactory } from '../factories';
import { useCsrfTokenContext } from '../../../contexts';

interface DecodedToken {
  exp: number;
  iat: number;
  sub: string;
}

export const useAuth = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { csrfToken } = useCsrfTokenContext();
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('token'));
  const [user, setUser] = useState<DecodedToken | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const validateToken = useCallback((tokenToValidate: string) => {
    try {
      const decoded = jwtDecode<DecodedToken>(tokenToValidate);
      if (decoded.exp * 1000 < Date.now()) {
        return false;
      }
      return decoded;
    } catch {
      return false;
    }
  }, []);

  const initializeAuth = useCallback(() => {
    const savedToken = localStorage.getItem('token');
    if (savedToken) {
      const decoded = validateToken(savedToken);
      if (decoded) {
        setToken(savedToken);
        setUser(decoded);
      } else {
        localStorage.removeItem('token');
        setToken(null);
        setUser(null);
      }
    }
    setIsLoading(false);
  }, [validateToken]);

  useEffect(() => {
    initializeAuth();
  }, [initializeAuth]);

  const login = useCallback(
    async (email: string, password: string) => {
      try {
        const loginUseCase = LoginFactory(csrfToken);
        const response = await loginUseCase.authenticate({ email, password });

        console.log(response);
        if (!response.data?.access_token) {
          toast.error(t('loginError'));
          return false;
        }

        const decoded = validateToken(response.data.access_token);
        if (!decoded) {
          toast.error(t('invalidToken'));
          return false;
        }

        localStorage.setItem('token', response.data.access_token);
        setToken(response.data.access_token);
        setUser(decoded);
        toast.success(t('loginSuccess'));
        navigate('/', { replace: true });
        return true;
      } catch (error) {
        toast.error(`Erro no login: ${error instanceof Error ? error.message : 'Erro desconhecido'}`);
        return false;
      }
    },
    [navigate, validateToken, t, csrfToken]
  );

  const logout = useCallback(() => {
    localStorage.removeItem('token');
    setToken(null);
    setUser(null);
    navigate('/login', { replace: true });
  }, [navigate]);

  return {
    isAuthenticated: !!token && !!user,
    isLoading,
    user,
    login,
    logout,
    token,
  };
};