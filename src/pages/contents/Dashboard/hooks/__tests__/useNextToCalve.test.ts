import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useNextToCalve } from '../useNextToCalve';
import { useFarm } from '../../../../../hooks/useFarm';

// Mock do useFarm
vi.mock('../../../../../hooks/useFarm', () => ({
  useFarm: vi.fn()
}));

// Mock da factory
vi.mock('../../factories/usecases/get-next-to-calve-factory', () => ({
  GetNextToCalveFactory: vi.fn(() => ({
    getNextToCalve: vi.fn()
  }))
}));

const mockUseFarm = vi.mocked(useFarm);

describe('useNextToCalve', () => {
  const mockFarm = { id: 1, name: 'Test Farm' };
  const mockNextToCalveData = [
    {
      id: 1,
      animal_name: 'Tata Salt',
      ear_tag_number_local: 123,
      photo: 'src/assets/images/mocked/cows/tata.png',
      pregnancy_date: '2024-01-01',
      expected_birth_date: '2024-10-01',
      days_until_birth: 15,
      status: 'Alto' as const
    },
    {
      id: 2,
      animal_name: 'Lays',
      ear_tag_number_local: 124,
      photo: 'src/assets/images/mocked/cows/lays.png',
      pregnancy_date: '2024-01-15',
      expected_birth_date: '2024-10-15',
      days_until_birth: 40,
      status: 'Baixo' as const
    }
  ];

  beforeEach(() => {
    vi.clearAllMocks();
    mockUseFarm.mockReturnValue({
      farm: mockFarm,
      setFarm: vi.fn(),
      clearFarm: vi.fn()
    });
  });

  it('deve retornar dados iniciais corretos', async () => {
    const { result } = renderHook(() => useNextToCalve());

    // Inicialmente loading deve ser true porque o useEffect dispara automaticamente
    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.nextToCalve).toEqual([]);
    expect(result.current.error).toBe(null);
  });

  it('deve buscar próximas vacas a parir quando farmId for fornecido', async () => {
    const mockGetNextToCalve = vi.fn().mockResolvedValue({
      data: mockNextToCalveData,
      success: true,
      message: 'Success',
      status: 200
    });

    const { GetNextToCalveFactory } = await import('../../factories/usecases/get-next-to-calve-factory');
    vi.mocked(GetNextToCalveFactory).mockReturnValue({
      getNextToCalve: mockGetNextToCalve
    });

    const { result } = renderHook(() => useNextToCalve(1));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(mockGetNextToCalve).toHaveBeenCalledWith({ farm_id: 1 });
    expect(result.current.nextToCalve).toEqual(mockNextToCalveData);
  });

  it('deve usar farm do contexto quando farmId não for fornecido', async () => {
    const mockGetNextToCalve = vi.fn().mockResolvedValue({
      data: mockNextToCalveData,
      success: true,
      message: 'Success',
      status: 200
    });

    const { GetNextToCalveFactory } = await import('../../factories/usecases/get-next-to-calve-factory');
    vi.mocked(GetNextToCalveFactory).mockReturnValue({
      getNextToCalve: mockGetNextToCalve
    });

    const { result } = renderHook(() => useNextToCalve());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(mockGetNextToCalve).toHaveBeenCalledWith({ farm_id: 1 });
  });

  it('deve lidar com erro quando farm não for encontrada', async () => {
    mockUseFarm.mockReturnValue({
      farm: null,
      setFarm: vi.fn(),
      clearFarm: vi.fn()
    });

    const { result } = renderHook(() => useNextToCalve());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.nextToCalve).toEqual([]);
    expect(result.current.error).toBe(null);
  });

  it('deve lidar com erro na requisição', async () => {
    const mockGetNextToCalve = vi.fn().mockRejectedValue(new Error('Network error'));

    const { GetNextToCalveFactory } = await import('../../factories/usecases/get-next-to-calve-factory');
    vi.mocked(GetNextToCalveFactory).mockReturnValue({
      getNextToCalve: mockGetNextToCalve
    });

    const { result } = renderHook(() => useNextToCalve(1));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toBe('Network error');
    expect(result.current.nextToCalve).toEqual([]);
  });

  it('deve permitir refetch dos dados', async () => {
    const mockGetNextToCalve = vi.fn().mockResolvedValue({
      data: mockNextToCalveData,
      success: true,
      message: 'Success',
      status: 200
    });

    const { GetNextToCalveFactory } = await import('../../factories/usecases/get-next-to-calve-factory');
    vi.mocked(GetNextToCalveFactory).mockReturnValue({
      getNextToCalve: mockGetNextToCalve
    });

    const { result } = renderHook(() => useNextToCalve(1));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // Test refetch
    await result.current.refetch();

    expect(mockGetNextToCalve).toHaveBeenCalledTimes(2);
  });

  it('deve permitir limpar erro', () => {
    const { result } = renderHook(() => useNextToCalve());

    // Simular erro
    result.current.clearError();

    expect(result.current.error).toBe(null);
  });
});
