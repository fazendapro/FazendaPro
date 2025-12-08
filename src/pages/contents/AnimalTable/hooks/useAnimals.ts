import { useState, useCallback, useEffect } from 'react';
import { useFarm } from '../../../../hooks/useFarm';
import { GetAnimalsByFarmFactory } from '../factories/usecases/get-animals-by-farm-factory';
import { Animal } from '../types/type';

interface UseAnimalsOptions {
  page?: number;
  limit?: number;
  searchTerm?: string;
}

export const useAnimals = (farmId?: number, options?: UseAnimalsOptions) => {
  const { farm } = useFarm();
  const [animals, setAnimals] = useState<Animal[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const shouldPaginate = !options?.searchTerm?.trim();
  const page = shouldPaginate ? options?.page ?? 1 : undefined;
  const limit = shouldPaginate ? options?.limit ?? 10 : undefined;
  const searchTerm = options?.searchTerm ?? '';

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
        limit,
        searchTerm: searchTerm.trim() || undefined
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
  }, [farmId, farm?.id, page, limit, searchTerm]);

  useEffect(() => {
    if (farmId || farm?.id) {
      getAnimalsByFarm();
    }
  }, [farmId, farm?.id, page, limit, searchTerm, getAnimalsByFarm]);

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