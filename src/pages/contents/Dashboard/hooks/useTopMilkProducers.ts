import { useState, useCallback, useEffect } from 'react';
import { useFarm } from '../../../../hooks/useFarm';
import { GetTopMilkProducersFactory } from '../factories/usecases/get-top-milk-producers-factory';
import { TopMilkProducer } from '../types/dashboard.types';

export const useTopMilkProducers = (farmId?: number, limit?: number, periodDays?: number) => {
  const { farm } = useFarm();
  const [topProducers, setTopProducers] = useState<TopMilkProducer[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getTopMilkProducers = useCallback(async (): Promise<TopMilkProducer[]> => {
    const targetFarmId = farmId || farm?.id;
    if (!targetFarmId) {
      return [];
    }

    setLoading(true);
    setError(null);

    try {
      const getTopMilkProducersUseCase = GetTopMilkProducersFactory();
      const result = await getTopMilkProducersUseCase.getTopMilkProducers({
        farm_id: targetFarmId,
        limit,
        period_days: periodDays
      });
      const topProducersData = result.data || [];
      setTopProducers(topProducersData);
      return topProducersData;
    } catch (err: any) {
      const errorMessage = err.message || 'Erro ao buscar maiores produtoras de leite';
      setError(errorMessage);
      return [];
    } finally {
      setLoading(false);
    }
  }, [farmId, farm?.id, limit, periodDays]);

  useEffect(() => {
    if (farmId || farm?.id) {
      getTopMilkProducers();
    }
  }, [farmId, farm?.id, getTopMilkProducers]);

  return {
    topProducers,
    getTopMilkProducers,
    refetch: getTopMilkProducers,
    loading,
    error,
    clearError: () => setError(null)
  };
};
