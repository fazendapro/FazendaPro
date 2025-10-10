import { useMemo, useState, useEffect, useCallback } from 'react';
import { useFarm } from '../contexts/FarmContext';
import { GetFarmFactory } from '../pages/contents/Settings/factories';
import { BackendFarmData } from '../pages/contents/Settings/types/farm-types';

export const useSelectedFarm = () => {
  const { selectedFarm, setSelectedFarm, clearSelectedFarm } = useFarm();
  const [farmLogo, setFarmLogo] = useState<string>('');
  
  const farmId = useMemo(() => selectedFarm?.ID || null, [selectedFarm?.ID]);
  const farmName = useMemo(() => {
    if (!selectedFarm) return null;
    return selectedFarm.Company?.CompanyName || `Fazenda ${selectedFarm.ID}`;
  }, [selectedFarm]);

  const loadFarmData = useCallback(async () => {
    if (!farmId) {
      setFarmLogo('');
      return;
    }
    
    try {
      const getFarmUseCase = GetFarmFactory.create();
      const response = await getFarmUseCase.get(farmId);
      
      if (response.data) {
        const backendData = response.data as BackendFarmData;
        setFarmLogo(backendData.Logo || '');
      }
    } catch {
      setFarmLogo('');
    }
  }, [farmId]);

  useEffect(() => {
    loadFarmData();
  }, [loadFarmData]);
  
  return {
    selectedFarm,
    setSelectedFarm,
    clearSelectedFarm,
    farmId,
    farmName,
    farmLogo,
  };
};
