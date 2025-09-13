import { AxiosInstance } from 'axios'
import { baseAxios } from './base'

export function api(prefix: string = ''): AxiosInstance {
  return baseAxios(prefix)
}