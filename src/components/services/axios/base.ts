import axios, { AxiosInstance } from 'axios'

export function baseAxios(baseUrl: string): AxiosInstance {
  return axios.create({ baseURL: baseUrl })
}