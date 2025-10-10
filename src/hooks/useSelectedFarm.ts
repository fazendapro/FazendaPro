import { useMemo } from 'react';
import { useFarm } from '../contexts/FarmContext';

export const useSelectedFarm = () => {
  const { selectedFarm, setSelectedFarm, clearSelectedFarm } = useFarm();
  
  const farmId = useMemo(() => selectedFarm?.ID || null, [selectedFarm?.ID]);
  const farmName = useMemo(() => selectedFarm ? `Fazenda ${selectedFarm.ID}` : null, [selectedFarm]);
  
  return {
    selectedFarm,
    setSelectedFarm,
    clearSelectedFarm,
    farmId,
    farmName,
  };
};
