import { useState, useCallback, useEffect } from 'react';
import { useFarm } from '../../../../hooks/useFarm';
import { GetAnimalsByFarmFactory } from '../factories/usecases/get-animals-by-farm-factory';
import { Animal } from '../types/type';

export const useAnimals = (farmId?: number) => {
  const { farm } = useFarm();
  const [animals, setAnimals] = useState<Animal[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

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
      const result = await getAnimalsByFarmUseCase.getAnimalsByFarm({ farm_id: targetFarmId });
      const animalsData = result.data || [];
      setAnimals(animalsData);
      return animalsData;
    } catch (err: any) {
      const errorMessage = err.message || 'Erro ao buscar animais';
      setError(errorMessage);
      return [];
    } finally {
      setLoading(false);
    }
  }, [farmId, farm?.id]);

  useEffect(() => {
    if (farmId || farm?.id) {
      getAnimalsByFarm();
    }
  }, [farmId, farm?.id, getAnimalsByFarm]);

  return {
    animals,
    getAnimalsByFarm,
    refetch: getAnimalsByFarm,
    loading,
    error,
    clearError: () => setError(null)
  };
};