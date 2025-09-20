import { useState, useEffect, useCallback } from 'react';
import { jwtDecode } from 'jwt-decode';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { useTranslation } from 'react-i18next';
import { LoginFactory, AuthFactory } from '../factories';
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
  const [refreshToken, setRefreshToken] = useState<string | null>(() => localStorage.getItem('refreshToken'));
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

  const clearAuth = useCallback(() => {
    localStorage.removeItem('token');
    localStorage.removeItem('refreshToken');
    setToken(null);
    setRefreshToken(null);
    setUser(null);
    setIsLoading(false);
  }, []);

  const refreshAccessToken = useCallback(async (refreshTokenValue: string) => {
    try {
      const authUseCase = AuthFactory(csrfToken);
      const response = await authUseCase.refreshToken({ refresh_token: refreshTokenValue });

      if (response.success && response.access_token) {
        localStorage.setItem('token', response.access_token);
        setToken(response.access_token);
        
        const decoded = validateToken(response.access_token);
        if (decoded) {
          setUser(decoded);
        }
      } else {
        clearAuth();
      }
    } catch (error) {
      console.error('Erro ao renovar token:', error);
      clearAuth();
    }
    setIsLoading(false);
  }, [validateToken, clearAuth, csrfToken]);

  const initializeAuth = useCallback(() => {
    const savedToken = localStorage.getItem('token');
    const savedRefreshToken = localStorage.getItem('refreshToken');
    
    if (savedToken) {
      const decoded = validateToken(savedToken);
      if (decoded) {
        setToken(savedToken);
        setUser(decoded);
        if (savedRefreshToken) {
          setRefreshToken(savedRefreshToken);
        }
        setIsLoading(false);
      } else {
        if (savedRefreshToken) return refreshAccessToken(savedRefreshToken);
        return clearAuth();
      }
    } else {
      setIsLoading(false);
    }
  }, [validateToken, refreshAccessToken, clearAuth]);

  useEffect(() => {
    initializeAuth();
  }, [initializeAuth]);

  const login = useCallback(
    async (email: string, password: string) => {
      try {
        const loginUseCase = LoginFactory(csrfToken);
        const response = await loginUseCase.authenticate({ email, password });

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

        if (response.refresh_token) {
          localStorage.setItem('refreshToken', response.refresh_token);
          setRefreshToken(response.refresh_token);
        }

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

  const logout = useCallback(async () => {
    // Tentar invalidar refresh token no backend
    if (refreshToken) {
      try {
        const authUseCase = AuthFactory(csrfToken);
        await authUseCase.logout({ refresh_token: refreshToken });
      } catch (error) {
        console.error('Erro ao fazer logout no backend:', error);
      }
    }
    
    clearAuth();
    navigate('/login', { replace: true });
  }, [navigate, refreshToken, clearAuth, csrfToken]);

  return {
    isAuthenticated: !!token && !!user,
    isLoading,
    user,
    login,
    logout,
    token,
  };
};