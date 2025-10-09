import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios'
import { apiConfig } from '../../../config/api'

export function baseAxios(baseUrl: string): AxiosInstance {

  const instance = axios.create({ 
    baseURL: `${apiConfig.baseUrl}/${baseUrl}`,
    timeout: apiConfig.timeout
  })

  instance.interceptors.request.use(
    (config: AxiosRequestConfig) => {
      const token = localStorage.getItem('token')
      if (token) {
        config.headers = config.headers || {}
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    },
    (error: AxiosError) => {
      return Promise.reject(error)
    }
  )

  instance.interceptors.response.use(
    (response: AxiosResponse) => response,
    async (error: AxiosError) => {
      const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean }

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