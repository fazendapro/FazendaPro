import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { FarmProvider, useFarm } from '../FarmContext';

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <FarmProvider>{children}</FarmProvider>
);

describe('FarmContext', () => {
  beforeEach(() => {
    localStorage.clear();
    vi.clearAllMocks();
  });

  afterEach(() => {
    localStorage.clear();
  });

  describe('Inicialização', () => {
    it('deve inicializar sem fazenda selecionada quando não há fazenda salva', () => {
      const { result } = renderHook(() => useFarm(), { wrapper });

      expect(result.current.selectedFarm).toBeNull();
    });

    it('deve inicializar com fazenda salva no localStorage', () => {
      const mockFarm = {
        ID: 1,
        CompanyID: 1,
        Logo: 'logo.png',
        Company: {
          ID: 1,
          CompanyName: 'Test Farm',
          Location: 'Test Location',
          FarmCNPJ: '12345678901234',
        },
      };

      localStorage.setItem('selectedFarm', JSON.stringify(mockFarm));

      const { result } = renderHook(() => useFarm(), { wrapper });

      expect(result.current.selectedFarm).toEqual(mockFarm);
    });

    it('deve limpar localStorage quando dados são inválidos', () => {
      localStorage.setItem('selectedFarm', 'invalid-json');

      const { result } = renderHook(() => useFarm(), { wrapper });

      expect(result.current.selectedFarm).toBeNull();
      expect(localStorage.getItem('selectedFarm')).toBeNull();
    });
  });

  describe('setSelectedFarm', () => {
    it('deve definir fazenda selecionada e salvar no localStorage', () => {
      const mockFarm = {
        ID: 1,
        CompanyID: 1,
        Logo: 'logo.png',
        Company: {
          ID: 1,
          CompanyName: 'Test Farm',
          Location: 'Test Location',
          FarmCNPJ: '12345678901234',
        },
      };

      const { result } = renderHook(() => useFarm(), { wrapper });

      act(() => {
        result.current.setSelectedFarm(mockFarm);
      });

      expect(result.current.selectedFarm).toEqual(mockFarm);
      expect(localStorage.getItem('selectedFarm')).toBe(JSON.stringify(mockFarm));
    });

    it('deve atualizar fazenda selecionada', () => {
      const mockFarm1 = {
        ID: 1,
        CompanyID: 1,
        Logo: 'logo1.png',
      };

      const mockFarm2 = {
        ID: 2,
        CompanyID: 2,
        Logo: 'logo2.png',
      };

      const { result } = renderHook(() => useFarm(), { wrapper });

      act(() => {
        result.current.setSelectedFarm(mockFarm1);
      });

      expect(result.current.selectedFarm).toEqual(mockFarm1);

      act(() => {
        result.current.setSelectedFarm(mockFarm2);
      });

      expect(result.current.selectedFarm).toEqual(mockFarm2);
      expect(localStorage.getItem('selectedFarm')).toBe(JSON.stringify(mockFarm2));
    });

    it('deve remover fazenda do localStorage quando null é passado', () => {
      const mockFarm = {
        ID: 1,
        CompanyID: 1,
        Logo: 'logo.png',
      };

      const { result } = renderHook(() => useFarm(), { wrapper });

      act(() => {
        result.current.setSelectedFarm(mockFarm);
      });

      expect(result.current.selectedFarm).toEqual(mockFarm);

      act(() => {
        result.current.setSelectedFarm(null);
      });

      expect(result.current.selectedFarm).toBeNull();
      expect(localStorage.getItem('selectedFarm')).toBeNull();
    });
  });

  describe('clearSelectedFarm', () => {
    it('deve limpar fazenda selecionada e remover do localStorage', () => {
      const mockFarm = {
        ID: 1,
        CompanyID: 1,
        Logo: 'logo.png',
      };

      const { result } = renderHook(() => useFarm(), { wrapper });

      act(() => {
        result.current.setSelectedFarm(mockFarm);
      });

      expect(result.current.selectedFarm).toEqual(mockFarm);

      act(() => {
        result.current.clearSelectedFarm();
      });

      expect(result.current.selectedFarm).toBeNull();
      expect(localStorage.getItem('selectedFarm')).toBeNull();
    });

    it('deve funcionar mesmo quando não há fazenda selecionada', () => {
      const { result } = renderHook(() => useFarm(), { wrapper });

      expect(result.current.selectedFarm).toBeNull();

      act(() => {
        result.current.clearSelectedFarm();
      });

      expect(result.current.selectedFarm).toBeNull();
      expect(localStorage.getItem('selectedFarm')).toBeNull();
    });
  });

  describe('useFarm hook', () => {
    it('deve lançar erro quando usado fora do FarmProvider', () => {
      // Suprime o erro esperado no console
      const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {});

      expect(() => {
        renderHook(() => useFarm());
      }).toThrow('useFarm deve ser usado dentro de um FarmProvider');

      consoleError.mockRestore();
    });
  });

  describe('Persistência', () => {
    it('deve persistir fazenda entre renderizações', () => {
      const mockFarm = {
        ID: 1,
        CompanyID: 1,
        Logo: 'logo.png',
      };

      const { result, rerender } = renderHook(() => useFarm(), { wrapper });

      act(() => {
        result.current.setSelectedFarm(mockFarm);
      });

      expect(result.current.selectedFarm).toEqual(mockFarm);

      rerender();

      expect(result.current.selectedFarm).toEqual(mockFarm);
      expect(localStorage.getItem('selectedFarm')).toBe(JSON.stringify(mockFarm));
    });

    it('deve carregar fazenda do localStorage em nova instância', () => {
      const mockFarm = {
        ID: 1,
        CompanyID: 1,
        Logo: 'logo.png',
      };

      localStorage.setItem('selectedFarm', JSON.stringify(mockFarm));

      const { result } = renderHook(() => useFarm(), { wrapper });

      expect(result.current.selectedFarm).toEqual(mockFarm);
    });
  });
});




