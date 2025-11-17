import axios from 'axios'
import { apiConfig } from '../../../config/api'

export function baseAxios(baseUrl: string) {
  const instance = axios.create({ 
    baseURL: `${apiConfig.baseUrl}/${baseUrl}`,
    timeout: apiConfig.timeout
  })

  instance.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem('token')
      if (token) {
        config.headers = config.headers || {}
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    },
    (error) => {
      return Promise.reject(error)
    }
  )

  instance.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config

      if (error.response?.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true

        const refreshToken = localStorage.getItem('refreshToken')
        if (refreshToken) {
          try {
            const { AuthFactory } = await import('../../../pages/Login/factories')
            const authUseCase = AuthFactory()
            const response = await authUseCase.refreshToken({ refresh_token: refreshToken })

            if (response.success && response.access_token) {
              localStorage.setItem('token', response.access_token)
              originalRequest.headers = originalRequest.headers || {}
              originalRequest.headers.Authorization = `Bearer ${response.access_token}`
              return instance(originalRequest)
            }
          } catch {
            // Erro ao renovar token, será tratado abaixo
          }
        }

        localStorage.removeItem('token')
        localStorage.removeItem('refreshToken')
        window.location.href = '/login'
      }
      return Promise.reject(error)
    }
  )

  return instance
}