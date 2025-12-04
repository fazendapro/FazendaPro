import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
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

  const abortControllerRef = useRef<AbortController | null>(null);
  const isRequestInProgressRef = useRef(false);
  const requestIdRef = useRef<string | null>(null);
  const setSelectedFarmRef = useRef(setSelectedFarm);
  const navigateRef = useRef(navigate);

  useEffect(() => {
    setSelectedFarmRef.current = setSelectedFarm;
    navigateRef.current = navigate;
  }, [setSelectedFarm, navigate]);

  useEffect(() => {
    const currentRequestId = `farm-load-${Date.now()}-${Math.random()}`;

    if (isRequestInProgressRef.current && abortControllerRef.current) {
      abortControllerRef.current.abort();
    }

    if (isRequestInProgressRef.current && requestIdRef.current === currentRequestId) {
      return;
    }

    const loadFarms = async () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }

      const abortController = new AbortController();
      abortControllerRef.current = abortController;
      isRequestInProgressRef.current = true;
      requestIdRef.current = currentRequestId;

      try {
        setLoading(true);
        setError(null);

        const response = await farmSelectionService.getUserFarms();

        if (abortController.signal.aborted || requestIdRef.current !== currentRequestId) {
          return;
        }
        
        if (response.success) {
          setFarms(response.farms);
          setAutoSelect(response.auto_select);
          
          if (response.auto_select && response.selected_farm_id) {
            setSelectedFarmId(response.selected_farm_id);
            setTimeout(() => {
              if (!abortController.signal.aborted && requestIdRef.current === currentRequestId) {
                navigateRef.current('/');
              }
            }, 2000);
          }
        } else {
          setError('Erro ao carregar fazendas');
        }
      } catch (err) {
        if (err instanceof Error && err.name === 'AbortError') {
          return;
        }
        if (!abortController.signal.aborted && requestIdRef.current === currentRequestId) {
          setError(err instanceof Error ? err.message : 'Erro desconhecido');
        }
      } finally {
        if (!abortController.signal.aborted && requestIdRef.current === currentRequestId) {
          setLoading(false);
        }
        if (requestIdRef.current === currentRequestId) {
          isRequestInProgressRef.current = false;
        }
      }
    };

    loadFarms();

    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
        abortControllerRef.current = null;
      }
      if (requestIdRef.current === currentRequestId) {
        isRequestInProgressRef.current = false;
      }
    };
  }, []);

  const selectFarm = async (farmId: number) => {
    try {
      setError(null);

      const response = await farmSelectionService.selectFarm(farmId);

      if (response.success && response.access_token) {
        updateToken(response.access_token);

        setSelectedFarmId(farmId);

        const selectedFarm = farms.find((farm: Farm) => farm.ID === farmId);
        if (selectedFarm) {
          setSelectedFarmRef.current(selectedFarm);
        }

        navigate('/');
      } else {
        const errorMessage = 'Erro ao selecionar fazenda';
        setError(errorMessage);
        toast.error(errorMessage, {
          toastId: 'farm-selection-select-error-toast',
          autoClose: 5000,
        });
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Erro ao selecionar fazenda';
      setError(errorMessage);
      toast.error(errorMessage, {
        toastId: 'farm-selection-select-error-toast',
        autoClose: 5000,
      });
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
