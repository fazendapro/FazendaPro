import { useState, useCallback, useEffect } from 'react';
import { useFarm } from '../../../../hooks/useFarm';
import { GetMonthlySalesAndPurchasesFactory } from '../factories/usecases/get-monthly-sales-and-purchases-factory';
import { MonthlySalesAndPurchasesData } from '../types/dashboard.types';

export const useMonthlySalesAndPurchases = (farmId?: number, months: number = 12) => {
  const { farm } = useFarm();
  const [data, setData] = useState<MonthlySalesAndPurchasesData>({ sales: [], purchases: [] });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getMonthlySalesAndPurchases = useCallback(async (): Promise<MonthlySalesAndPurchasesData> => {
    const targetFarmId = farmId || farm?.id;
    if (!targetFarmId) {
      return { sales: [], purchases: [] };
    }

    setLoading(true);
    setError(null);

    try {
      const getMonthlySalesAndPurchasesUseCase = GetMonthlySalesAndPurchasesFactory();
      const result = await getMonthlySalesAndPurchasesUseCase.getMonthlySalesAndPurchases({ 
        farm_id: targetFarmId,
        months 
      });
      const monthlyData = result.data || { sales: [], purchases: [] };
      setData(monthlyData);
      return monthlyData;
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Erro ao buscar dados mensais de vendas e compras';
      setError(errorMessage);
      return { sales: [], purchases: [] };
    } finally {
      setLoading(false);
    }
  }, [farmId, farm?.id, months]);

  useEffect(() => {
    if (farmId || farm?.id) {
      getMonthlySalesAndPurchases();
    }
  }, [farmId, farm?.id, months, getMonthlySalesAndPurchases]);

  return {
    data,
    getMonthlySalesAndPurchases,
    refetch: getMonthlySalesAndPurchases,
    loading,
    error,
    clearError: () => setError(null)
  };
};


