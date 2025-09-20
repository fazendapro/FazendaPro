import axios, { AxiosInstance } from 'axios'

export function baseAxios(baseUrl: string): AxiosInstance {
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  
  if (import.meta.env.DEV) {
    console.log('API URL:', apiUrl)
  }

  const instance = axios.create({ 
    baseURL: `${apiUrl}/${baseUrl}` 
  })

  instance.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem('token')
      if (token) {
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