import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useTopMilkProducers } from '../useTopMilkProducers';
import { useFarm } from '../../../../../hooks/useFarm';

vi.mock('../../../../../hooks/useFarm', () => ({
  useFarm: vi.fn()
}));

const mockGetTopMilkProducers = vi.fn();
vi.mock('../../factories/usecases/get-top-milk-producers-factory', () => ({
  GetTopMilkProducersFactory: vi.fn(() => ({
    getTopMilkProducers: mockGetTopMilkProducers
  }))
}));

const mockUseFarm = vi.mocked(useFarm);

describe('useTopMilkProducers', () => {
  const mockFarm = { id: 1, name: 'Test Farm', location: 'Test Location', created_at: '2021-01-01', updated_at: '2021-01-01' };
  const mockTopProducersData = [
    {
      id: 1,
      animal_name: 'Surf Excel',
      ear_tag_number_local: 123,
      photo: 'src/assets/images/mocked/cows/surf.png',
      total_production: 900,
      average_daily_production: 30,
      fat_content: 3.5,
      last_collection_date: '2024-01-15',
      days_in_lactation: 120
    },
    {
      id: 2,
      animal_name: 'Rin',
      ear_tag_number_local: 124,
      photo: 'src/assets/images/mocked/cows/rin.png',
      total_production: 630,
      average_daily_production: 21,
      fat_content: 4.2,
      last_collection_date: '2024-01-15',
      days_in_lactation: 100
    }
  ];

  beforeEach(() => {
    vi.clearAllMocks();
    mockGetTopMilkProducers.mockReset();
    mockUseFarm.mockReturnValue({
      farm: mockFarm,
      loading: false,
      error: null,
      refetch: vi.fn()
    });
  });

  it('deve retornar dados iniciais corretos', async () => {
    mockGetTopMilkProducers.mockResolvedValue({
      data: [],
      success: true,
      message: 'Success',
      status: 200
    });


    const { result } = renderHook(() => useTopMilkProducers());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.topProducers).toEqual([]);
    expect(result.current.error).toBe(null);
  });

  it('deve buscar maiores produtoras quando farmId for fornecido', async () => {
    mockGetTopMilkProducers.mockResolvedValue({
      data: mockTopProducersData,
      success: true,
      message: 'Success',
      status: 200
    });


    const { result } = renderHook(() => useTopMilkProducers(1));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(mockGetTopMilkProducers).toHaveBeenCalledWith({
      farm_id: 1,
      limit: undefined,
      period_days: undefined
    });
    expect(result.current.topProducers).toEqual(mockTopProducersData);
  });

  it('deve usar farm do contexto quando farmId não for fornecido', async () => {
    mockUseFarm.mockReturnValue({
      farm: mockFarm,
      loading: false,
      error: null,
      refetch: vi.fn()
    });

    mockGetTopMilkProducers.mockResolvedValue({
      data: mockTopProducersData,
      success: true,
      message: 'Success',
      status: 200
    });


    const { result } = renderHook(() => useTopMilkProducers());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    }, { timeout: 3000 });

    expect(mockGetTopMilkProducers).toHaveBeenCalledWith({
      farm_id: 1,
      limit: undefined,
      period_days: undefined
    });
  });

  it('deve usar parâmetros customizados quando fornecidos', async () => {
    mockGetTopMilkProducers.mockResolvedValue({
      data: mockTopProducersData,
      success: true,
      message: 'Success',
      status: 200
    });


    const { result } = renderHook(() => useTopMilkProducers(1, 5, 15));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(mockGetTopMilkProducers).toHaveBeenCalledWith({
      farm_id: 1,
      limit: 5,
      period_days: 15
    });
  });

  it('deve lidar com erro quando farm não for encontrada', async () => {
    mockUseFarm.mockReturnValue({
      farm: null,
      loading: false,
      error: null,
      refetch: vi.fn()
    });

    const { result } = renderHook(() => useTopMilkProducers());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.topProducers).toEqual([]);
    expect(result.current.error).toBe(null);
  });

  it('deve lidar com erro na requisição', async () => {
    mockGetTopMilkProducers.mockRejectedValue(new Error('Network error'));


    const { result } = renderHook(() => useTopMilkProducers(1));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toBe('Network error');
    expect(result.current.topProducers).toEqual([]);
  });

  it('deve permitir refetch dos dados', async () => {
    mockGetTopMilkProducers.mockResolvedValue({
      data: mockTopProducersData,
      success: true,
      message: 'Success',
      status: 200
    });


    const { result } = renderHook(() => useTopMilkProducers(1));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await result.current.refetch();

    expect(mockGetTopMilkProducers).toHaveBeenCalledTimes(2);
  });

  it('deve permitir limpar erro', () => {
    const { result } = renderHook(() => useTopMilkProducers());

    result.current.clearError();

    expect(result.current.error).toBe(null);
  });
});
