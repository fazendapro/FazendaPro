import { useMemo, useState, useEffect, useCallback } from 'react';
import { useFarm } from '../contexts/FarmContext';
import { GetFarmFactory } from '../pages/contents/Settings/factories';
import { BackendFarmData } from '../pages/contents/Settings/types/farm-types';
import i18n from '../locale/i18n';
import dayjs from '../config/dayjs';

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
        
        if (backendData.Language) {
          const language = backendData.Language;
          await i18n.changeLanguage(language);
          
          const dayjsLocaleMap: Record<string, string> = {
            'pt': 'pt-br',
            'en': 'en',
            'es': 'es',
          };
          const dayjsLocale = dayjsLocaleMap[language] || 'pt-br';
          dayjs.locale(dayjsLocale);
        }
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
