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

      console.log('API Response:', response);

      if (response.success) {
        // Garantir que response.data seja sempre um array
        const animalsData = Array.isArray(response.data) ? response.data : [];
        console.log('Animals data after processing:', animalsData);
        setAnimals(animalsData);
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
  }, [farmId]);

  return {
    animals,
    loading,
    error,
    refetch: fetchAnimals
  };
}; 