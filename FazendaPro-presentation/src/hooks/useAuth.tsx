import { useState, useEffect } from 'react'
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
  const [token, setToken] = useState<string | null>(() => {
    const savedToken = localStorage.getItem('token')
    return savedToken
  })
  const [user, setUser] = useState<DecodedToken | null>(null)

  useEffect(() => {
    if (!!token && !!user) {
      navigate('/');
    }
  }, [token, user, navigate]);

  useEffect(() => {
    const savedToken = localStorage.getItem('token')
    if (savedToken) {
      setToken(savedToken)
    }
  }, [])

  useEffect(() => {
    if (token) {
      try {
        const decoded = jwtDecode<DecodedToken>(token)
        console.log('Token decodificado:', decoded)
        if (decoded.exp * 1000 < Date.now()) {
          console.log('Token expirado')
          logout()
        } else {
          setUser(decoded)
        }
      } catch (error) {
        console.error('Erro ao decodificar token:', error)
        logout()
      }
    }
  }, [token])

  const login = async (email: string, password: string) => {
  try {
    const response = await api.post('/auth/login', { email, password });
    const newToken = response.data.access_token;

    if (!newToken) {
      return false;
    }
    
    localStorage.setItem('token', newToken);
    setToken(newToken);
    
    try {
      const decoded = jwtDecode<DecodedToken>(newToken);
      setUser(decoded);
      
      setTimeout(() => {
        navigate('/');
      }, 100); 
      
    } catch (error) {
      console.error('Erro ao decodificar token:', error);
      return false;
    }

    return true;
  } catch (error) {
    console.error('Erro durante login:', error);
    return false;
  }
};

  const logout = () => {
    localStorage.removeItem('token')
    setToken(null)
    setUser(null)
  }

  return {
    isAuthenticated: !!token && !!user,
    user,
    login,
    logout,
    token
  }
}