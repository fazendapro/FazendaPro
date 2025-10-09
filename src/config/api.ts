
interface ApiConfig {
  baseUrl: string
  timeout: number
  retryAttempts: number
}

function getApiUrl(): string {
  if ((import.meta as any).env?.DEV) {
    return (import.meta as any).env?.VITE_API_URL || 'http://localhost:8080'
  }

  return (import.meta as any).env?.VITE_API_URL
}

export const apiConfig: ApiConfig = {
  baseUrl: getApiUrl(),
  timeout: 10000,
  retryAttempts: 3
}

if ((import.meta as any).env?.DEV) {
  console.log('🔧 API Configuration:', {
    baseUrl: apiConfig.baseUrl,
    environment: (import.meta as any).env?.MODE,
    isDev: (import.meta as any).env?.DEV
  })
}
