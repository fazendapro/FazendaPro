import { describe, it, expect, beforeEach } from 'vitest'
import { api, apiConfig } from '../api'

describe('api config', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('deve ter configuração base correta', () => {
    expect(apiConfig.timeout).toBe(10000)
    expect(apiConfig.retryAttempts).toBe(3)
  })

  it('deve ter configuração de baseURL e timeout', () => {
    expect(api.defaults).toBeDefined()
    expect(api.defaults.timeout).toBe(10000)
    expect(api.defaults.baseURL).toBeDefined()
  })

  it('deve ter interceptors configurados', () => {
    expect(api.interceptors.request).toBeDefined()
    expect(api.interceptors.response).toBeDefined()
  })
})

