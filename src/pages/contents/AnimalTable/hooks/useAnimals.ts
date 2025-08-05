import { useState, useEffect } from 'react';
import { GetAnimalsByFarmFactory } from '../factories/usecases/get-animals-by-farm-factory';
import { Animal } from '../types/type';

export const useAnimals = (farmId: number) => {
  const [animals, setAnimals] = useState<Animal[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAnimals = async () => {
    setLoading(true);
    setError(null);

    try {
      const getAnimalsByFarm = GetAnimalsByFarmFactory();
      const response = await getAnimalsByFarm.getAnimalsByFarm({ farm_id: farmId });

      if (response.success) {
        setAnimals(response.data || []);
      } else {
        setError(response.message || 'Erro ao buscar animais');
      }
    } catch (err) {
      console.error('Error fetching animals:', err);
      setError(err instanceof Error ? err.message : 'Erro desconhecido');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAnimals();
  }, []);

  return {
    animals,
    loading,
    error,
    refetch: fetchAnimals
  };
}; 