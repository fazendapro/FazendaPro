import { useMemo } from 'react';
import { useFarm } from '../contexts/FarmContext';

export const useSelectedFarm = () => {
  const { selectedFarm, setSelectedFarm, clearSelectedFarm } = useFarm();
  
  const farmId = useMemo(() => selectedFarm?.ID || null, [selectedFarm?.ID]);
  const farmName = useMemo(() => {
    if (!selectedFarm) return null;
    return selectedFarm.Company?.CompanyName || `Fazenda ${selectedFarm.ID}`;
  }, [selectedFarm]);
  
  return {
    selectedFarm,
    setSelectedFarm,
    clearSelectedFarm,
    farmId,
    farmName,
  };
};
