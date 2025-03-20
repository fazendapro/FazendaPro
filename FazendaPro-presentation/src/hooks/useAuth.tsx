import { useState, useEffect } from 'react'
import { jwtDecode } from 'jwt-decode'
import axios from 'axios'

interface DecodedToken {
  exp: number
  iat: number
  sub: string
}

export const useAuth = () => {
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'))
  const [user, setUser] = useState<DecodedToken | null>(null)

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
        console.error(error)
        logout()
      }
    }
  }, [token])

  const login = async (email: string, password: string) => {
    try {
      const response = await axios.post('/api/auth/login', { email, password })
      const newToken = response.data.token
      localStorage.setItem('token', newToken)
      setToken(newToken)
      return true
    } catch (error) {
      console.error('Login failed:', error)
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