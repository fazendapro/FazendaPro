import { useState, useEffect } from 'react'
import { jwtDecode } from 'jwt-decode'
import axios from 'axios'
import { useNavigate } from 'react-router';
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
    const savedToken = localStorage.getItem('token')
    if (savedToken) {
      setToken(savedToken)
    }
  }, [])

  useEffect(() => {
    if (token) {
      try {
        const decoded = jwtDecode<DecodedToken>(token)
        if (decoded.exp * 1000 < Date.now()) {
          logout()
        } else {
          setUser(decoded)
        }
      } catch (error) {
        logout()
      }
    }
  }, [token])

  const login = async (email: string, password: string) => {
    try {
      const response = await api.post('/auth/login', { email, password })
      const newToken = response.data.access_token

      if (!newToken) {
        return false
      }
      
      navigate('/home')
      localStorage.setItem('token', newToken)

      setToken(newToken)

      return true
    } catch (error) {
      return false
    }
  }

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