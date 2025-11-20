
import axios from 'axios';

interface ApiConfig {
  baseUrl: string
  timeout: number
  retryAttempts: number
}

interface ImportMetaEnv {
  DEV?: boolean
  VITE_API_URL?: string
}

function getApiUrl(): string {
  const env = (import.meta as { env?: ImportMetaEnv }).env
  if (env?.DEV) {
    return env.VITE_API_URL || 'http://localhost:8080'
  }

  return env?.VITE_API_URL || ''
}

export const apiConfig: ApiConfig = {
  baseUrl: getApiUrl(),
  timeout: 10000,
  retryAttempts: 3
}

export const api = axios.create({
  baseURL: `${apiConfig.baseUrl}/api/v1`,
  timeout: apiConfig.timeout,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401 && 
        (error.response?.data?.message?.includes('token') || 
         error.response?.data?.message?.includes('unauthorized') ||
         error.response?.data?.message?.includes('authentication'))) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

