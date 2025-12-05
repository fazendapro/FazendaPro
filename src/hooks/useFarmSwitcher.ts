import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { farmSelectionService, Farm } from '../pages/FarmSelection/services/farm-selection-service';
import { useFarm } from '../contexts/FarmContext';
import { useAuth } from '../contexts/AuthContext';

export interface UseFarmSwitcherReturn {
  farms: Farm[];
  loading: boolean;
  error: string | null;
  switchFarm: (farmId: number) => Promise<void>;
  loadFarms: () => Promise<void>;
}

export const useFarmSwitcher = (): UseFarmSwitcherReturn => {
  const [farms, setFarms] = useState<Farm[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const navigate = useNavigate();
  const { setSelectedFarm, selectedFarm } = useFarm();
  const { updateToken } = useAuth();

  const loadFarms = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await farmSelectionService.getUserFarms();

      if (response.success) {
        setFarms(response.farms || []);
      } else {
        const errorMessage = response.message || 'Erro ao carregar fazendas';
        setError(errorMessage);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Erro ao carregar fazendas';
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  }, []);

  const switchFarm = useCallback(async (farmId: number) => {
    // Não fazer nada se já estiver na fazenda selecionada
    if (selectedFarm?.ID === farmId) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const response = await farmSelectionService.selectFarm(farmId);

      if (response.success && response.access_token) {
        updateToken(response.access_token);

        let selectedFarmData = farms.find((farm: Farm) => farm.ID === farmId);

        if (!selectedFarmData) {
          const farmsResponse = await farmSelectionService.getUserFarms();
          if (farmsResponse.success && farmsResponse.farms) {
            setFarms(farmsResponse.farms);
            selectedFarmData = farmsResponse.farms.find((farm: Farm) => farm.ID === farmId);
          }
        }

        if (selectedFarmData) {
          setSelectedFarm(selectedFarmData);
        }

        toast.success('Fazenda alterada com sucesso!', {
          toastId: 'farm-switch-success-toast',
          autoClose: 3000,
        });

        const currentPath = window.location.pathname;
        if (currentPath === '/farm-selection') {
          navigate('/', { replace: true });
        }
      } else {
        const errorMessage = response.message || 'Erro ao trocar de fazenda';
        setError(errorMessage);
        toast.error(errorMessage, {
          toastId: 'farm-switch-error-toast',
          autoClose: 5000,
        });
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Erro ao trocar de fazenda';
      setError(errorMessage);
      toast.error(errorMessage, {
        toastId: 'farm-switch-error-toast',
        autoClose: 5000,
      });
    } finally {
      setLoading(false);
    }
  }, [selectedFarm, farms, setSelectedFarm, updateToken, navigate]);

  return {
    farms,
    loading,
    error,
    switchFarm,
    loadFarms,
  };
};

