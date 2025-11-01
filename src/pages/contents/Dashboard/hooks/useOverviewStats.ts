import { useState, useCallback, useEffect } from 'react';
import { useFarm } from '../../../../hooks/useFarm';
import { GetOverviewStatsFactory } from '../factories/usecases/get-overview-stats-factory';
import { OverviewStats } from '../types/dashboard.types';

export const useOverviewStats = (farmId?: number) => {
  const { farm } = useFarm();
  const [stats, setStats] = useState<OverviewStats>({ 
    males_count: 0, 
    females_count: 0, 
    total_sold: 0, 
    total_revenue: 0 
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getOverviewStats = useCallback(async (): Promise<OverviewStats> => {
    const targetFarmId = farmId || farm?.id;
    if (!targetFarmId) {
      return { males_count: 0, females_count: 0, total_sold: 0, total_revenue: 0 };
    }

    setLoading(true);
    setError(null);

    try {
      const getOverviewStatsUseCase = GetOverviewStatsFactory();
      const result = await getOverviewStatsUseCase.getOverviewStats({
        farm_id: targetFarmId
      });
      const statsData = result.data || { males_count: 0, females_count: 0, total_sold: 0, total_revenue: 0 };
      setStats(statsData);
      return statsData;
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Erro ao buscar estatísticas gerais';
      setError(errorMessage);
      return { males_count: 0, females_count: 0, total_sold: 0, total_revenue: 0 };
    } finally {
      setLoading(false);
    }
  }, [farmId, farm?.id]);

  useEffect(() => {
    if (farmId || farm?.id) {
      getOverviewStats();
    }
  }, [farmId, farm?.id, getOverviewStats]);

  return {
    stats,
    getOverviewStats,
    refetch: getOverviewStats,
    loading,
    error,
    clearError: () => setError(null)
  };
};

