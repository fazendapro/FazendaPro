import { useMemo } from 'react';
import { useSelectedFarm } from './useSelectedFarm';

export const useFarm = () => {
  const { selectedFarm, farmName } = useSelectedFarm();

  const farm = useMemo(() => {
    return selectedFarm ? {
      id: selectedFarm.ID,
      name: farmName || `Fazenda ${selectedFarm.ID}`,
      location: 'Localização não disponível',
      created_at: '2021-01-01',
      updated_at: '2021-01-01'
    } : null;
  }, [selectedFarm, farmName]);

  return {
    farm,
    loading: false,
    error: null,
    refetch: () => {}
  };
}; 