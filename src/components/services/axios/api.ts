import { AxiosInstance } from 'axios'
import { baseAxios } from './base'

export function api(domain: string, prefix: string = 'api'): AxiosInstance {
  return baseAxios(domain + `/${prefix}`)
}