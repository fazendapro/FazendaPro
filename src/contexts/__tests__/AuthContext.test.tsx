import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { AuthProvider, useAuth } from '../AuthContext';
import { FarmProvider } from '../FarmContext';
import { toast } from 'react-toastify';
import { jwtDecode } from 'jwt-decode';

// Mock jwt-decode
vi.mock('jwt-decode', () => ({
  jwtDecode: vi.fn(),
}));

// Mock react-router-dom
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

// Mock react-toastify
vi.mock('react-toastify', () => ({
  toast: {
    error: vi.fn(),
    success: vi.fn(),
  },
}));

// Mock react-i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

// Mock factories
const mockLoginUseCase = {
  authenticate: vi.fn(),
};

const mockAuthUseCase = {
  refreshToken: vi.fn(),
  logout: vi.fn(),
};

vi.mock('../../pages/Login/factories', () => ({
  LoginFactory: vi.fn(() => mockLoginUseCase),
  AuthFactory: vi.fn(() => mockAuthUseCase),
}));

// Mock farm selection service
vi.mock('../../pages/FarmSelection/services/farm-selection-service', () => ({
  farmSelectionService: {
    getUserFarms: vi.fn(),
  },
}));

import { farmSelectionService as mockFarmSelectionService } from '../../pages/FarmSelection/services/farm-selection-service';

// Mock CSRF token context
vi.mock('../csrf-token', () => ({
  useCsrfTokenContext: () => ({
    csrfToken: 'test-csrf-token',
  }),
}));

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <BrowserRouter>
    <FarmProvider>
      <AuthProvider>{children}</AuthProvider>
    </FarmProvider>
  </BrowserRouter>
);

describe('AuthContext', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
    mockNavigate.mockClear();
  });

  afterEach(() => {
    localStorage.clear();
  });

  describe('Inicialização', () => {
    it('deve inicializar sem token quando não há token salvo', async () => {
      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.isAuthenticated).toBe(false);
      expect(result.current.user).toBeNull();
      expect(result.current.token).toBeNull();
    });

    it('deve inicializar com token válido quando há token salvo', async () => {
      const mockToken = 'valid-token';
      const mockDecodedToken = {
        exp: Math.floor(Date.now() / 1000) + 3600,
        iat: Math.floor(Date.now() / 1000),
        sub: '1',
        email: 'test@example.com',
      };

      localStorage.setItem('token', mockToken);
      vi.mocked(jwtDecode).mockReturnValue(mockDecodedToken as never);

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.isAuthenticated).toBe(true);
      expect(result.current.user).toEqual(mockDecodedToken);
      expect(result.current.token).toBe(mockToken);
    });

    it('deve tentar refresh token quando token está expirado', async () => {
      const expiredToken = 'expired-token';
      const refreshToken = 'refresh-token';
      const newToken = 'new-token';
      const mockDecodedToken = {
        exp: Math.floor(Date.now() / 1000) + 3600,
        iat: Math.floor(Date.now() / 1000),
        sub: '1',
        email: 'test@example.com',
      };

      localStorage.setItem('token', expiredToken);
      localStorage.setItem('refreshToken', refreshToken);

      // Primeira chamada retorna token expirado
      vi.mocked(jwtDecode).mockReturnValueOnce({
        exp: Math.floor(Date.now() / 1000) - 3600,
        iat: Math.floor(Date.now() / 1000) - 7200,
        sub: '1',
        email: 'test@example.com',
      } as never);

      // Mock do refresh token bem-sucedido
      mockAuthUseCase.refreshToken.mockResolvedValue({
        success: true,
        access_token: newToken,
      });

      // Segunda chamada retorna token válido
      vi.mocked(jwtDecode).mockReturnValueOnce(mockDecodedToken as never);

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      }, { timeout: 3000 });

      expect(mockAuthUseCase.refreshToken).toHaveBeenCalledWith({
        refresh_token: refreshToken,
      });
      expect(localStorage.getItem('token')).toBe(newToken);
    });

    it('deve limpar tokens quando refresh token falha', async () => {
      const expiredToken = 'expired-token';
      const refreshToken = 'refresh-token';

      localStorage.setItem('token', expiredToken);
      localStorage.setItem('refreshToken', refreshToken);

      vi.mocked(jwtDecode).mockReturnValueOnce({
        exp: Math.floor(Date.now() / 1000) - 3600,
        iat: Math.floor(Date.now() / 1000) - 7200,
        sub: '1',
        email: 'test@example.com',
      } as never);

      mockAuthUseCase.refreshToken.mockRejectedValue(new Error('Refresh failed'));

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      }, { timeout: 3000 });

      expect(localStorage.getItem('token')).toBeNull();
      expect(localStorage.getItem('refreshToken')).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
    });
  });

  describe('Login', () => {
    it('deve fazer login com sucesso', async () => {
      const mockToken = 'access-token';
      const mockRefreshToken = 'refresh-token';
      const mockDecodedToken = {
        exp: Math.floor(Date.now() / 1000) + 3600,
        iat: Math.floor(Date.now() / 1000),
        sub: '1',
        email: 'test@example.com',
      };

      mockLoginUseCase.authenticate.mockResolvedValue({
        data: {
          access_token: mockToken,
        },
        refresh_token: mockRefreshToken,
      });

      mockFarmSelectionService.getUserFarms.mockResolvedValue({
        success: true,
        auto_select: true,
        farms: [{ ID: 1, CompanyID: 1, Logo: '' }],
      });

      vi.mocked(jwtDecode).mockReturnValue(mockDecodedToken as never);

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      let loginResult = false;
      await act(async () => {
        loginResult = await result.current.login('test@example.com', 'password');
      });

      expect(loginResult).toBe(true);
      expect(mockLoginUseCase.authenticate).toHaveBeenCalledWith({
        email: 'test@example.com',
        password: 'password',
      });
      expect(localStorage.getItem('token')).toBe(mockToken);
      expect(localStorage.getItem('refreshToken')).toBe(mockRefreshToken);
      expect(result.current.isAuthenticated).toBe(true);
      expect(result.current.user).toEqual(mockDecodedToken);
      expect(toast.success).toHaveBeenCalled();
    });

    it('deve falhar login quando não há access_token', async () => {
      mockLoginUseCase.authenticate.mockResolvedValue({
        data: null,
      });

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      let loginResult = false;
      await act(async () => {
        loginResult = await result.current.login('test@example.com', 'password');
      });

      expect(loginResult).toBe(false);
      expect(toast.error).toHaveBeenCalledWith('loginError');
      expect(result.current.isAuthenticated).toBe(false);
    });

    it('deve falhar login quando token é inválido', async () => {
      mockLoginUseCase.authenticate.mockResolvedValue({
        data: {
          access_token: 'invalid-token',
        },
      });

      vi.mocked(jwtDecode).mockReturnValue(null as never);

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      let loginResult = false;
      await act(async () => {
        loginResult = await result.current.login('test@example.com', 'password');
      });

      expect(loginResult).toBe(false);
      expect(toast.error).toHaveBeenCalledWith('invalidToken');
    });

    it('deve navegar para farm-selection quando auto_select é false', async () => {
      const mockToken = 'access-token';
      const mockDecodedToken = {
        exp: Math.floor(Date.now() / 1000) + 3600,
        iat: Math.floor(Date.now() / 1000),
        sub: '1',
        email: 'test@example.com',
      };

      mockLoginUseCase.authenticate.mockResolvedValue({
        data: {
          access_token: mockToken,
        },
      });

      mockFarmSelectionService.getUserFarms.mockResolvedValue({
        success: true,
        auto_select: false,
        farms: [{ ID: 1, CompanyID: 1, Logo: '' }],
      });

      vi.mocked(jwtDecode).mockReturnValue(mockDecodedToken as never);

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      await act(async () => {
        await result.current.login('test@example.com', 'password');
      });

      expect(mockNavigate).toHaveBeenCalledWith('/farm-selection', { replace: true });
    });

    it('deve navegar para home quando auto_select é true', async () => {
      const mockToken = 'access-token';
      const mockDecodedToken = {
        exp: Math.floor(Date.now() / 1000) + 3600,
        iat: Math.floor(Date.now() / 1000),
        sub: '1',
        email: 'test@example.com',
      };

      mockLoginUseCase.authenticate.mockResolvedValue({
        data: {
          access_token: mockToken,
        },
      });

      mockFarmSelectionService.getUserFarms.mockResolvedValue({
        success: true,
        auto_select: true,
        farms: [{ ID: 1, CompanyID: 1, Logo: '' }],
      });

      vi.mocked(jwtDecode).mockReturnValue(mockDecodedToken as never);

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      await act(async () => {
        await result.current.login('test@example.com', 'password');
      });

      expect(mockNavigate).toHaveBeenCalledWith('/', { replace: true });
    });

    it('deve tratar erro durante login', async () => {
      mockLoginUseCase.authenticate.mockRejectedValue(new Error('Network error'));

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      let loginResult = false;
      await act(async () => {
        loginResult = await result.current.login('test@example.com', 'password');
      });

      expect(loginResult).toBe(false);
      expect(toast.error).toHaveBeenCalled();
    });
  });

  describe('Logout', () => {
    it('deve fazer logout com sucesso', async () => {
      const mockToken = 'access-token';
      const mockRefreshToken = 'refresh-token';
      const mockDecodedToken = {
        exp: Math.floor(Date.now() / 1000) + 3600,
        iat: Math.floor(Date.now() / 1000),
        sub: '1',
        email: 'test@example.com',
      };

      localStorage.setItem('token', mockToken);
      localStorage.setItem('refreshToken', mockRefreshToken);

      vi.mocked(jwtDecode).mockReturnValue(mockDecodedToken as never);
      mockAuthUseCase.logout.mockResolvedValue(undefined);

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      await act(async () => {
        await result.current.logout();
      });

      expect(mockAuthUseCase.logout).toHaveBeenCalledWith({
        refresh_token: mockRefreshToken,
      });
      expect(localStorage.getItem('token')).toBeNull();
      expect(localStorage.getItem('refreshToken')).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
      expect(result.current.user).toBeNull();
      expect(mockNavigate).toHaveBeenCalledWith('/login', { replace: true });
    });

    it('deve fazer logout mesmo quando refresh token não existe', async () => {
      const mockToken = 'access-token';
      const mockDecodedToken = {
        exp: Math.floor(Date.now() / 1000) + 3600,
        iat: Math.floor(Date.now() / 1000),
        sub: '1',
        email: 'test@example.com',
      };

      localStorage.setItem('token', mockToken);

      vi.mocked(jwtDecode).mockReturnValue(mockDecodedToken as never);

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      await act(async () => {
        await result.current.logout();
      });

      expect(localStorage.getItem('token')).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
      expect(mockNavigate).toHaveBeenCalledWith('/login', { replace: true });
    });

    it('deve fazer logout mesmo quando chamada de logout falha', async () => {
      const mockToken = 'access-token';
      const mockRefreshToken = 'refresh-token';
      const mockDecodedToken = {
        exp: Math.floor(Date.now() / 1000) + 3600,
        iat: Math.floor(Date.now() / 1000),
        sub: '1',
        email: 'test@example.com',
      };

      localStorage.setItem('token', mockToken);
      localStorage.setItem('refreshToken', mockRefreshToken);

      vi.mocked(jwtDecode).mockReturnValue(mockDecodedToken as never);
      mockAuthUseCase.logout.mockRejectedValue(new Error('Logout failed'));

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      await act(async () => {
        await result.current.logout();
      });

      expect(localStorage.getItem('token')).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
      expect(mockNavigate).toHaveBeenCalledWith('/login', { replace: true });
    });
  });

  describe('Validação de Token', () => {
    it('deve validar token válido', async () => {
      const mockToken = 'valid-token';
      const mockDecodedToken = {
        exp: Math.floor(Date.now() / 1000) + 3600,
        iat: Math.floor(Date.now() / 1000),
        sub: '1',
        email: 'test@example.com',
      };

      localStorage.setItem('token', mockToken);
      vi.mocked(jwtDecode).mockReturnValue(mockDecodedToken as never);

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.isAuthenticated).toBe(true);
      expect(result.current.user).toEqual(mockDecodedToken);
    });

    it('deve rejeitar token expirado', async () => {
      const expiredToken = 'expired-token';

      localStorage.setItem('token', expiredToken);
      vi.mocked(jwtDecode).mockReturnValue({
        exp: Math.floor(Date.now() / 1000) - 3600,
        iat: Math.floor(Date.now() / 1000) - 7200,
        sub: '1',
        email: 'test@example.com',
      } as never);

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      }, { timeout: 3000 });

      expect(result.current.isAuthenticated).toBe(false);
    });
  });
});

