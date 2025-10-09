// @ts-ignore
import axios from 'axios'
import { apiConfig } from '../../../config/api'

// Tipos customizados para interceptors
interface RequestConfig {
  headers?: Record<string, string>
  _retry?: boolean
}

interface ResponseData {
  data: unknown
  status: number
  statusText: string
}

interface ErrorData {
  response?: {
    status: number
    data: unknown
  }
  config: RequestConfig
}

export function baseAxios(baseUrl: string) {
  const instance = axios.create({ 
    baseURL: `${apiConfig.baseUrl}/${baseUrl}`,
    timeout: apiConfig.timeout
  })

  instance.interceptors.request.use(
    (config: RequestConfig) => {
      const token = localStorage.getItem('token')
      if (token) {
        config.headers = config.headers || {}
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    },
    (error: ErrorData) => {
      return Promise.reject(error)
    }
  )

  instance.interceptors.response.use(
    (response: ResponseData) => response,
    async (error: ErrorData) => {
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
          } catch (refreshError) {
            console.error('Erro ao renovar token:', refreshError)
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