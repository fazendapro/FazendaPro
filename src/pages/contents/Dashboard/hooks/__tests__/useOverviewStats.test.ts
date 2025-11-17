import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useOverviewStats } from '../useOverviewStats';
import { useFarm } from '../../../../../hooks/useFarm';

vi.mock('../../../../../hooks/useFarm', () => ({
  useFarm: vi.fn()
}));

const mockGetOverviewStats = vi.fn();
vi.mock('../../factories/usecases/get-overview-stats-factory', () => ({
  GetOverviewStatsFactory: vi.fn(() => ({
    getOverviewStats: mockGetOverviewStats
  }))
}));

const mockUseFarm = vi.mocked(useFarm);

describe('useOverviewStats', () => {
  const mockFarm = { id: 1, name: 'Test Farm', location: 'Test Location', created_at: '2021-01-01', updated_at: '2021-01-01' };
  const mockStats = {
    males_count: 10,
    females_count: 15,
    total_sold: 5,
    total_revenue: 15000.50
  };

  beforeEach(() => {
    vi.clearAllMocks();
    mockGetOverviewStats.mockReset();
    mockUseFarm.mockReturnValue({
      farm: mockFarm,
      loading: false,
      error: null,
      refetch: vi.fn()
    });
  });

  it('deve retornar dados iniciais corretos', async () => {
    mockGetOverviewStats.mockResolvedValue({
      data: { males_count: 0, females_count: 0, total_sold: 0, total_revenue: 0 },
      success: true,
      message: 'Success',
      status: 200
    });

    const { result } = renderHook(() => useOverviewStats());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.stats).toEqual({ males_count: 0, females_count: 0, total_sold: 0, total_revenue: 0 });
    expect(result.current.error).toBe(null);
  });

  it('deve buscar estatísticas quando farmId for fornecido', async () => {
    mockGetOverviewStats.mockResolvedValue({
      data: mockStats,
      success: true,
      message: 'Success',
      status: 200
    });

    const { result } = renderHook(() => useOverviewStats(1));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(mockGetOverviewStats).toHaveBeenCalledWith({
      farm_id: 1
    });
    expect(result.current.stats).toEqual(mockStats);
  });

  it('deve usar farm do contexto quando farmId não for fornecido', async () => {
    mockUseFarm.mockReturnValue({
      farm: mockFarm,
      loading: false,
      error: null,
      refetch: vi.fn()
    });

    mockGetOverviewStats.mockResolvedValue({
      data: mockStats,
      success: true,
      message: 'Success',
      status: 200
    });

    const { result } = renderHook(() => useOverviewStats());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    }, { timeout: 3000 });

    expect(mockGetOverviewStats).toHaveBeenCalledWith({
      farm_id: 1
    });
  });

  it('deve lidar com erro quando farm não for encontrada', async () => {
    mockUseFarm.mockReturnValue({
      farm: null,
      loading: false,
      error: null,
      refetch: vi.fn()
    });

    const { result } = renderHook(() => useOverviewStats());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.stats).toEqual({ males_count: 0, females_count: 0, total_sold: 0, total_revenue: 0 });
    expect(result.current.error).toBe(null);
  });

  it('deve lidar com erro na requisição', async () => {
    mockGetOverviewStats.mockRejectedValue(new Error('Network error'));

    const { result } = renderHook(() => useOverviewStats(1));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toBe('Network error');
    expect(result.current.stats).toEqual({ males_count: 0, females_count: 0, total_sold: 0, total_revenue: 0 });
  });

  it('deve permitir refetch dos dados', async () => {
    mockGetOverviewStats.mockResolvedValue({
      data: mockStats,
      success: true,
      message: 'Success',
      status: 200
    });

    const { result } = renderHook(() => useOverviewStats(1));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await result.current.refetch();

    expect(mockGetOverviewStats).toHaveBeenCalledTimes(2);
  });

  it('deve permitir limpar erro', () => {
    const { result } = renderHook(() => useOverviewStats());

    result.current.clearError();

    expect(result.current.error).toBe(null);
  });
});

