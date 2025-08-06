import axios, { AxiosInstance } from 'axios'

export function baseAxios(baseUrl: string): AxiosInstance {
  const apiUrl = import.meta.env.VITE_API_URL

  return axios.create({
    baseURL: `${apiUrl}/${baseUrl}` 
  })
}