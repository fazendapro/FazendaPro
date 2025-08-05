import { AxiosInstance } from 'axios'
import { baseAxios } from './base'

export function api(prefix: string = 'api/v1'): AxiosInstance {
  return baseAxios(prefix)
}