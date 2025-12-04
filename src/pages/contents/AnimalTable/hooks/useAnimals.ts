import { useState, useCallback, useEffect } from 'react';
import { useFarm } from '../../../../hooks/useFarm';
import { GetAnimalsByFarmFactory } from '../factories/usecases/get-animals-by-farm-factory';
import { Animal } from '../types/type';

interface UseAnimalsOptions {
  page?: number;
  limit?: number;
}

export const useAnimals = (farmId?: number, options?: UseAnimalsOptions) => {
  const { farm } = useFarm();
  const [animals, setAnimals] = useState<Animal[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const page = options?.page ?? 1;
  const limit = options?.limit ?? 10;

  const getAnimalsByFarm = useCallback(async (): Promise<Animal[]> => {
    const targetFarmId = farmId || farm?.id;
    if (!targetFarmId) {
      setError('Fazenda não encontrada');
      return [];
    }

    setLoading(true);
    setError(null);

    try {
      const getAnimalsByFarmUseCase = GetAnimalsByFarmFactory();
      const result = await getAnimalsByFarmUseCase.getAnimalsByFarm({ 
        farm_id: targetFarmId,
        page,
        limit
      });
      
      const paginatedData = result.data;
      const animalsData = paginatedData?.animals || [];
      const totalCount = paginatedData?.total || 0;
      
      setAnimals(animalsData);
      setTotal(totalCount);
      return animalsData;
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Erro ao buscar animais';
      setError(errorMessage);
      return [];
    } finally {
      setLoading(false);
    }
  }, [farmId, farm?.id, page, limit]);

  useEffect(() => {
    if (farmId || farm?.id) {
      getAnimalsByFarm();
    }
  }, [farmId, farm?.id, page, limit, getAnimalsByFarm]);

  return {
    animals,
    total,
    getAnimalsByFarm,
    refetch: getAnimalsByFarm,
    loading,
    error,
    clearError: () => setError(null)
  };
};