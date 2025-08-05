// import { useState, useEffect } from 'react';
// import { GetFarmFactory } from '../factories/usecases/get-farm-factory';
// import { Farm } from '../types/farm';

export const useFarm = () => {
  // const [farm, setFarm] = useState<Farm | null>(null);
  // const [loading, setLoading] = useState(false);
  // const [error, setError] = useState<string | null>(null);

  // const fetchFarm = async () => {
  //   setLoading(true);
  //   setError(null);
    
  //   try {
  //     const getFarm = GetFarmFactory();
  //     const response = await getFarm.getFarm({ farm_id: farmId });
      
  //     if (response.success) {
  //       setFarm(response.data || null);
  //     } else {
  //       setError(response.message || 'Erro ao buscar informações da fazenda');
  //     }
  //   } catch (err) {
  //     setError(err instanceof Error ? err.message : 'Erro desconhecido');
  //   } finally {
  //     setLoading(false);
  //   }
  // };

  // useEffect(() => {
  //   fetchFarm();
  // }, [farmId]);

  const farm = {
    id: 1,
    name: 'Fazenda 1',
    location: 'Rua 1, 123',
    created_at: '2021-01-01',
    updated_at: '2021-01-01'
  }

  return {
    farm,
    // loading,
    // error,
    // refetch: fetchFarm
  };
}; 