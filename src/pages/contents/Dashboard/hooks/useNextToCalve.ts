import { useState, useCallback, useEffect } from 'react';
import { useFarm } from '../../../../hooks/useFarm';
import { GetNextToCalveFactory } from '../factories/usecases/get-next-to-calve-factory';
import { NextToCalveAnimal } from '../types/dashboard.types';

export const useNextToCalve = (farmId?: number) => {
  const { farm } = useFarm();
  const [nextToCalve, setNextToCalve] = useState<NextToCalveAnimal[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getNextToCalve = useCallback(async (): Promise<NextToCalveAnimal[]> => {
    const targetFarmId = farmId || farm?.id;
    if (!targetFarmId) {
      return [];
    }

    setLoading(true);
    setError(null);

    try {
      const getNextToCalveUseCase = GetNextToCalveFactory();
      const result = await getNextToCalveUseCase.getNextToCalve({ farm_id: targetFarmId });
      const nextToCalveData = result.data || [];
      setNextToCalve(nextToCalveData);
      return nextToCalveData;
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Erro ao buscar próximas vacas a parir';
      setError(errorMessage);
      return [];
    } finally {
      setLoading(false);
    }
  }, [farmId, farm?.id]);

  useEffect(() => {
    if (farmId || farm?.id) {
      getNextToCalve();
    }
  }, [farmId, farm?.id, getNextToCalve]);

  return {
    nextToCalve,
    getNextToCalve,
    refetch: getNextToCalve,
    loading,
    error,
    clearError: () => setError(null)
  };
};
