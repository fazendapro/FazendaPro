import axios, { AxiosInstance } from 'axios'

export function baseAxios(baseUrl: string): AxiosInstance {
  const apiUrl = import.meta.env.VITE_API_URL
  
  if (!apiUrl) {
    throw new Error('VITE_API_URL não está definida. Verifique as variáveis de ambiente.')
  }

  return axios.create({
    baseURL: `${apiUrl}/${baseUrl}` 
  })
}