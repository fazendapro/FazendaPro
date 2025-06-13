import { AxiosInstance } from 'axios'
import { baseAxios } from './base'

export function api(domain: string, prefix: string = 'o'): AxiosInstance {
  return baseAxios(domain + `/api/v1/${prefix}`)
}