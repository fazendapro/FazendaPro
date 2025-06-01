import { useState, useEffect, useCallback } from 'react'
import { jwtDecode } from 'jwt-decode'
import axios from 'axios'
import { useNavigate } from 'react-router-dom';

interface DecodedToken {
  exp: number
  iat: number
  sub: string
}

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL
});

export const useAuth = () => {
  const navigate = useNavigate();
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('token'));
  const [user, setUser] = useState<DecodedToken | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const validateToken = useCallback((tokenToValidate: string) => {
    try {
      const decoded = jwtDecode<DecodedToken>(tokenToValidate);
      if (decoded.exp * 500000 < Date.now()) {
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

  const login = async (email: string, password: string) => {
    try {
      const response = await api.post('/auth/login', { email, password });
      const newToken = response.data.access_token;

      if (!newToken) {
        return false;
      }

      const decoded = validateToken(newToken);
      if (!decoded) {
        return false;
      }

      localStorage.setItem('token', newToken);
      setToken(newToken);
      setUser(decoded);
      navigate('/', { replace: true });
      return true;
    } catch {
      return false;
    }
  };

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
    token
  };
};