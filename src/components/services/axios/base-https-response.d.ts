export interface BaseHttpResponse<T> {
  message: string
  status: number
  data?: T
}