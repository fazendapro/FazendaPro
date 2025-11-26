import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { farmSelectionService, Farm } from '../services/farm-selection-service';
import { useFarm } from '../../../contexts/FarmContext';
import { useAuth } from '../../../contexts/AuthContext';

export interface UseFarmSelectionReturn {
  farms: Farm[];
  loading: boolean;
  error: string | null;
  autoSelect: boolean;
  selectedFarmId: number | null;
  selectFarm: (farmId: number) => Promise<void>;
}

export const useFarmSelection = (): UseFarmSelectionReturn => {
  const [farms, setFarms] = useState<Farm[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [autoSelect, setAutoSelect] = useState(false);
  const [selectedFarmId, setSelectedFarmId] = useState<number | null>(null);
  
  const navigate = useNavigate();
  const { setSelectedFarm } = useFarm();
  const { updateToken } = useAuth();

  useEffect(() => {
    const loadFarms = async () => {
      try {
        setLoading(true);
        setError(null);
        
        const response = await farmSelectionService.getUserFarms();
        
        if (response.success) {
          setFarms(response.farms);
          setAutoSelect(response.auto_select);
          
          if (response.auto_select && response.selected_farm_id) {
            setSelectedFarmId(response.selected_farm_id);
            setTimeout(() => {
              navigate('/');
            }, 2000);
          }
        } else {
          setError('Erro ao carregar fazendas');
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Erro desconhecido');
      } finally {
        setLoading(false);
      }
    };

    loadFarms();
  }, [navigate, setSelectedFarm]);

  const selectFarm = async (farmId: number) => {
    try {
      setError(null);

      const response = await farmSelectionService.selectFarm(farmId);

      if (response.success && response.access_token) {
        updateToken(response.access_token);

        setSelectedFarmId(farmId);

        const selectedFarm = farms.find((farm: Farm) => farm.ID === farmId);
        if (selectedFarm) {
          setSelectedFarm(selectedFarm);
        }

        navigate('/');
      } else {
        setError('Erro ao selecionar fazenda');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao selecionar fazenda');
    }
  };

  return {
    farms,
    loading,
    error,
    autoSelect,
    selectedFarmId,
    selectFarm,
  };
};
