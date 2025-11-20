import React, { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react';
import { jwtDecode } from 'jwt-decode';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { useTranslation } from 'react-i18next';
import { LoginFactory, AuthFactory } from '../pages/Login/factories';
import { useCsrfTokenContext } from './csrf-token';
import { farmSelectionService } from '../pages/FarmSelection/services/farm-selection-service';
import { useFarm } from './FarmContext';
import { DecodedToken, AuthContextType } from './AuthContext.types';

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { csrfToken } = useCsrfTokenContext();
  const { setSelectedFarm } = useFarm();
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('token'));
  const [refreshToken, setRefreshToken] = useState<string | null>(() => localStorage.getItem('refreshToken'));
  const [user, setUser] = useState<DecodedToken | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [initialized, setInitialized] = useState(false);

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

  useEffect(() => {
    if (initialized) return;
    
    const initializeAuth = async () => {
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
          setInitialized(true);
        } else {
          if (savedRefreshToken) {
            try {
              const authUseCase = AuthFactory(csrfToken);
              const response = await authUseCase.refreshToken({ refresh_token: savedRefreshToken });

              if (response.success && response.access_token) {
                const decoded = validateToken(response.access_token);
                if (decoded) {
                  localStorage.setItem('token', response.access_token);
                  setToken(response.access_token);
                  setUser(decoded);
                } else {
                  clearAuth();
                }
              } else {
                clearAuth();
              }
            } catch {
              clearAuth();
            }
          } else {
            localStorage.removeItem('token');
            localStorage.removeItem('refreshToken');
            setToken(null);
            setRefreshToken(null);
            setUser(null);
          }
          setIsLoading(false);
          setInitialized(true);
        }
      } else {
        setIsLoading(false);
        setInitialized(true);
      }
    };

    initializeAuth();
  }, [clearAuth, csrfToken, initialized, validateToken]);

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
        
        try {
          const farmResponse = await farmSelectionService.getUserFarms();
          
          if (farmResponse.success) {
            if (farmResponse.auto_select) {
              if (farmResponse.farms && farmResponse.farms.length > 0) {
                setSelectedFarm(farmResponse.farms[0]);
              }
              navigate('/', { replace: true });
            } else {
              navigate('/farm-selection', { replace: true });
              
              setTimeout(() => {
                if (window.location.pathname === '/login') {
                  window.location.href = '/farm-selection';
                }
              }, 100);
            }
          } else {
            navigate('/', { replace: true });
          }
        } catch {
          navigate('/', { replace: true });
        }
        
        return true;
      } catch (error) {
        toast.error(`Erro no login: ${error instanceof Error ? error.message : 'Erro desconhecido'}`);
        return false;
      }
    },
    [navigate, validateToken, t, csrfToken, setSelectedFarm]
  );

  const logout = useCallback(async () => {
    try {
      if (refreshToken) {
        const authUseCase = AuthFactory(csrfToken);
        await authUseCase.logout({ refresh_token: refreshToken });
      }
    } catch {
      void 0;
    } finally {
      localStorage.removeItem('token');
      localStorage.removeItem('refreshToken');
      setToken(null);
      setRefreshToken(null);
      setUser(null);
      setIsLoading(false);
      setSelectedFarm(null);
      navigate('/login', { replace: true });
    }
  }, [navigate, refreshToken, csrfToken, setSelectedFarm]);

  const value: AuthContextType = {
    isAuthenticated: !!token && !!user,
    isLoading,
    user,
    login,
    logout,
    token,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth deve ser usado dentro de um AuthProvider');
  }
  return context;
};
