import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useMonthlySalesAndPurchases } from '../useMonthlySalesAndPurchases';
import { useFarm } from '../../../../../hooks/useFarm';

vi.mock('../../../../../hooks/useFarm', () => ({
  useFarm: vi.fn()
}));

const mockGetMonthlySalesAndPurchases = vi.fn();
vi.mock('../../factories/usecases/get-monthly-sales-and-purchases-factory', () => ({
  GetMonthlySalesAndPurchasesFactory: vi.fn(() => ({
    getMonthlySalesAndPurchases: mockGetMonthlySalesAndPurchases
  }))
}));

const mockUseFarm = vi.mocked(useFarm);

describe('useMonthlySalesAndPurchases', () => {
  const mockFarm = { id: 1, name: 'Test Farm', location: 'Test Location', created_at: '2021-01-01', updated_at: '2021-01-01' };
  const mockMonthlyData = {
    sales: [
      { month: 'Jan', year: 2024, sales: 15000, count: 5 },
      { month: 'Fev', year: 2024, sales: 20000, count: 7 },
    ],
    purchases: [
      { month: 'Jan', year: 2024, total: 0 },
      { month: 'Fev', year: 2024, total: 0 },
    ],
  };

  beforeEach(() => {
    vi.clearAllMocks();
    mockGetMonthlySalesAndPurchases.mockReset();
    mockUseFarm.mockReturnValue({
      farm: mockFarm,
      loading: false,
      error: null,
      refetch: vi.fn()
    });
  });

  it('deve retornar dados iniciais corretos', async () => {
    mockGetMonthlySalesAndPurchases.mockResolvedValue({
      data: { sales: [], purchases: [] },
      success: true,
      message: 'Success',
      status: 200
    });


    const { result } = renderHook(() => useMonthlySalesAndPurchases());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.data).toEqual({ sales: [], purchases: [] });
    expect(result.current.error).toBe(null);
  });

  it('deve buscar dados mensais quando farmId for fornecido', async () => {
    mockGetMonthlySalesAndPurchases.mockResolvedValue({
      data: mockMonthlyData,
      success: true,
      message: 'Success',
      status: 200
    });


    const { result } = renderHook(() => useMonthlySalesAndPurchases(1));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(mockGetMonthlySalesAndPurchases).toHaveBeenCalledWith({
      farm_id: 1,
      months: 12
    });
    expect(result.current.data).toEqual(mockMonthlyData);
  });

  it('deve usar farm do contexto quando farmId não for fornecido', async () => {
    mockUseFarm.mockReturnValue({
      farm: mockFarm,
      loading: false,
      error: null,
      refetch: vi.fn()
    });

    mockGetMonthlySalesAndPurchases.mockResolvedValue({
      data: mockMonthlyData,
      success: true,
      message: 'Success',
      status: 200
    });


    const { result } = renderHook(() => useMonthlySalesAndPurchases());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    }, { timeout: 3000 });

    expect(mockGetMonthlySalesAndPurchases).toHaveBeenCalledWith({
      farm_id: 1,
      months: 12
    });
  });

  it('deve usar parâmetros customizados quando fornecidos', async () => {
    mockGetMonthlySalesAndPurchases.mockResolvedValue({
      data: mockMonthlyData,
      success: true,
      message: 'Success',
      status: 200
    });


    const { result } = renderHook(() => useMonthlySalesAndPurchases(1, 6));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(mockGetMonthlySalesAndPurchases).toHaveBeenCalledWith({
      farm_id: 1,
      months: 6
    });
  });

  it('deve lidar com erro quando farm não for encontrada', async () => {
    mockUseFarm.mockReturnValue({
      farm: null,
      loading: false,
      error: null,
      refetch: vi.fn()
    });

    const { result } = renderHook(() => useMonthlySalesAndPurchases());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.data).toEqual({ sales: [], purchases: [] });
    expect(result.current.error).toBe(null);
  });

  it('deve lidar com erro na requisição', async () => {
    mockGetMonthlySalesAndPurchases.mockRejectedValue(new Error('Network error'));


    const { result } = renderHook(() => useMonthlySalesAndPurchases(1));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toBe('Network error');
    expect(result.current.data).toEqual({ sales: [], purchases: [] });
  });

  it('deve permitir refetch dos dados', async () => {
    mockGetMonthlySalesAndPurchases.mockResolvedValue({
      data: mockMonthlyData,
      success: true,
      message: 'Success',
      status: 200
    });


    const { result } = renderHook(() => useMonthlySalesAndPurchases(1));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await result.current.refetch();

    expect(mockGetMonthlySalesAndPurchases).toHaveBeenCalledTimes(2);
  });

  it('deve permitir limpar erro', () => {
    const { result } = renderHook(() => useMonthlySalesAndPurchases());

    result.current.clearError();

    expect(result.current.error).toBe(null);
  });

  it('deve usar meses padrão (12) quando não fornecido', async () => {
    mockGetMonthlySalesAndPurchases.mockResolvedValue({
      data: mockMonthlyData,
      success: true,
      message: 'Success',
      status: 200
    });


    const { result } = renderHook(() => useMonthlySalesAndPurchases());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(mockGetMonthlySalesAndPurchases).toHaveBeenCalledWith({
      farm_id: 1,
      months: 12
    });
  });
});

