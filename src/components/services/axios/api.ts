import { AxiosInstance } from 'axios'
import { baseAxios } from './base'
import { api as centralApi } from '../../../config/api'

export function api(prefix: string = ''): AxiosInstance {
  if (prefix === '') {
    return centralApi
  }
  return baseAxios(prefix)
}