import { describe, it, expect, vi, beforeEach } from 'vitest';
import { api } from '../api';
import { baseAxios } from '../base';

// Mock baseAxios
const mockBaseAxiosInstance = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
};

vi.mock('../base', () => ({
  baseAxios: vi.fn(() => mockBaseAxiosInstance),
}));

// Mock central api
const mockCentralApi = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
};

vi.mock('../../../config/api', () => ({
  api: mockCentralApi,
}));

describe('api', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Factory de instâncias', () => {
    it('deve retornar api central quando prefix é vazio', () => {
      vi.clearAllMocks();
      api('');

      expect(baseAxios).not.toHaveBeenCalled();
    });

    it('deve retornar api central quando prefix não é fornecido', () => {
      vi.clearAllMocks();
      api();

      expect(baseAxios).not.toHaveBeenCalled();
    });

    it('deve criar nova instância com baseAxios quando prefix é fornecido', () => {
      vi.clearAllMocks();
      const instance = api('test-prefix');

      expect(baseAxios).toHaveBeenCalledWith('test-prefix');
      expect(instance).toBe(mockBaseAxiosInstance);
    });

    it('deve criar instâncias diferentes para prefixs diferentes', () => {
      vi.clearAllMocks();
      const instance1 = api('prefix1');
      const instance2 = api('prefix2');

      expect(baseAxios).toHaveBeenCalledWith('prefix1');
      expect(baseAxios).toHaveBeenCalledWith('prefix2');
      expect(instance1).toBe(mockBaseAxiosInstance);
      expect(instance2).toBe(mockBaseAxiosInstance);
    });
  });
});

