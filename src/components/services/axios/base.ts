import axios, { AxiosInstance } from 'axios'

export function baseAxios(baseUrl: string): AxiosInstance {
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  
  if (import.meta.env.DEV) {
    console.log('API URL:', apiUrl)
  }
  
  return axios.create({ 
    baseURL: `${apiUrl}/${baseUrl}` 
  })
}