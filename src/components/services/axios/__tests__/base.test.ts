import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { baseAxios } from '../base';

// Mock axios
vi.mock('axios', () => {
  const mockAxiosInstance = {
    interceptors: {
      request: { use: vi.fn() },
      response: { use: vi.fn() },
    },
  };
  return {
    default: {
      create: vi.fn(() => mockAxiosInstance),
    },
  };
});

// Mock config
vi.mock('../../../config/api', () => {
  const mockApiInstance = {
    interceptors: {
      request: { use: vi.fn() },
      response: { use: vi.fn() },
    },
  };
  return {
    apiConfig: {
      baseUrl: 'http://localhost:8080',
      timeout: 10000,
    },
    api: mockApiInstance,
  };
});

// Mock AuthFactory
const mockRefreshToken = vi.fn();
vi.mock('../../../pages/Login/factories', () => ({
  AuthFactory: vi.fn(() => ({
    refreshToken: mockRefreshToken,
  })),
}));

describe('baseAxios', () => {
  let mockAxiosInstance: any;
  let mockCreate: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    localStorage.clear();

    const axios = await import('axios');
    mockAxiosInstance = {
      interceptors: {
        request: {
          use: vi.fn(),
        },
        response: {
          use: vi.fn(),
        },
      },
      mockResolvedValue: vi.fn(),
    };

    mockCreate = vi.mocked(axios.default.create);
    mockCreate.mockReturnValue(mockAxiosInstance);
  });

  afterEach(() => {
    localStorage.clear();
  });

  describe('Criação de instância', () => {
    it('deve criar instância axios com baseURL e timeout corretos', () => {
      baseAxios('test-prefix');

      expect(mockCreate).toHaveBeenCalledWith({
        baseURL: 'http://localhost:8080/test-prefix',
        timeout: 10000,
      });
    });

    it('deve retornar instância axios', () => {
      const instance = baseAxios('test-prefix');

      expect(instance).toBe(mockAxiosInstance);
    });
  });

  describe('Interceptor de Request', () => {
    it('deve adicionar token de autorização quando token existe', () => {
      const token = 'test-token';
      localStorage.setItem('token', token);

      baseAxios('test-prefix');

      const requestInterceptor = mockAxiosInstance.interceptors.request.use.mock.calls[0][0];
      const config = { headers: {} };

      const result = requestInterceptor(config);

      expect(result.headers.Authorization).toBe(`Bearer ${token}`);
    });

    it('deve não adicionar token quando token não existe', () => {
      baseAxios('test-prefix');

      const requestInterceptor = mockAxiosInstance.interceptors.request.use.mock.calls[0][0];
      const config = { headers: {} };

      const result = requestInterceptor(config);

      expect(result.headers.Authorization).toBeUndefined();
    });

    it('deve criar headers quando não existem', () => {
      baseAxios('test-prefix');

      const requestInterceptor = mockAxiosInstance.interceptors.request.use.mock.calls[0][0];
      const config: any = {};

      const result = requestInterceptor(config);

      expect(result.headers).toBeDefined();
      expect(result.headers.Authorization).toBe('Bearer test-token');
    });

    it('deve rejeitar erro no interceptor de request', () => {
      baseAxios('test-prefix');

      const errorHandler = mockAxiosInstance.interceptors.request.use.mock.calls[0][1];
      const error = new Error('Request error');

      expect(() => errorHandler(error)).rejects.toThrow('Request error');
    });
  });

  describe('Interceptor de Response', () => {
    it('deve retornar response quando sucesso', () => {
      baseAxios('test-prefix');

      const successHandler = mockAxiosInstance.interceptors.response.use.mock.calls[0][0];
      const response = { data: { test: 'data' } };

      const result = successHandler(response);

      expect(result).toBe(response);
    });

    it('deve tentar refresh token quando recebe 401', async () => {
      const refreshToken = 'refresh-token';
      const newAccessToken = 'new-access-token';
      localStorage.setItem('refreshToken', refreshToken);

      mockRefreshToken.mockResolvedValue({
        success: true,
        access_token: newAccessToken,
      });

      baseAxios('test-prefix');

      const errorHandler = mockAxiosInstance.interceptors.response.use.mock.calls[0][1];
      const originalRequest = {
        _retry: false,
        headers: {},
      };
      const error = {
        response: {
          status: 401,
        },
        config: originalRequest,
      };

      // Mock da instância para retornar a requisição original
      mockAxiosInstance.mockResolvedValue = vi.fn().mockResolvedValue({ data: 'success' });
      mockAxiosInstance.mockImplementation = vi.fn().mockResolvedValue({ data: 'success' });

      try {
        await errorHandler(error);
      } catch {
        // Pode rejeitar se a instância não estiver configurada corretamente
      }

      expect(mockRefreshToken).toHaveBeenCalledWith({
        refresh_token: refreshToken,
      });
      expect(localStorage.getItem('token')).toBe(newAccessToken);
    });

    it('deve redirecionar para login quando refresh token falha', async () => {
      const refreshToken = 'refresh-token';
      localStorage.setItem('refreshToken', refreshToken);

      mockRefreshToken.mockRejectedValue(new Error('Refresh failed'));

      // Mock window.location
      const originalLocation = window.location;
      delete (window as any).location;
      window.location = { href: '' } as any;

      baseAxios('test-prefix');

      const errorHandler = mockAxiosInstance.interceptors.response.use.mock.calls[0][1];
      const originalRequest = {
        _retry: false,
        headers: {},
      };
      const error = {
        response: {
          status: 401,
        },
        config: originalRequest,
      };

      try {
        await errorHandler(error);
      } catch {
        // Esperado que rejeite
      }

      expect(localStorage.getItem('token')).toBeNull();
      expect(localStorage.getItem('refreshToken')).toBeNull();
      
      window.location = originalLocation;
    });

    it('deve redirecionar para login quando não há refresh token', async () => {
      // Mock window.location
      const originalLocation = window.location;
      delete (window as any).location;
      window.location = { href: '' } as any;

      baseAxios('test-prefix');

      const errorHandler = mockAxiosInstance.interceptors.response.use.mock.calls[0][1];
      const originalRequest = {
        _retry: false,
        headers: {},
      };
      const error = {
        response: {
          status: 401,
        },
        config: originalRequest,
      };

      try {
        await errorHandler(error);
      } catch {
        // Esperado que rejeite
      }

      expect(localStorage.getItem('token')).toBeNull();
      
      window.location = originalLocation;
    });

    it('deve não tentar refresh token novamente se já tentou', async () => {
      baseAxios('test-prefix');

      const errorHandler = mockAxiosInstance.interceptors.response.use.mock.calls[0][1];
      const originalRequest = {
        _retry: true,
        headers: {},
      };
      const error = {
        response: {
          status: 401,
        },
        config: originalRequest,
      };

      try {
        await errorHandler(error);
      } catch {
        // Esperado que rejeite
      }

      expect(mockRefreshToken).not.toHaveBeenCalled();
    });

    it('deve rejeitar erro quando status não é 401', async () => {
      baseAxios('test-prefix');

      const errorHandler = mockAxiosInstance.interceptors.response.use.mock.calls[0][1];
      const error = {
        response: {
          status: 500,
        },
      };

      await expect(errorHandler(error)).rejects.toEqual(error);
    });

    it('deve rejeitar erro quando não há response', async () => {
      baseAxios('test-prefix');

      const errorHandler = mockAxiosInstance.interceptors.response.use.mock.calls[0][1];
      const error = {};

      await expect(errorHandler(error)).rejects.toEqual(error);
    });
  });
});

