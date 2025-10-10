import React, { createContext, useContext, useState, ReactNode, useCallback, useMemo } from 'react';
import { Sale, CreateSaleRequest, UpdateSaleRequest, SaleFilters } from '../types/sale';
import { saleService } from '../components/services/saleService';
import { toast } from 'react-toastify';

interface SaleContextType {
  sales: Sale[];
  loading: boolean;
  error: string | null;
  createSale: (saleData: CreateSaleRequest) => Promise<void>;
  updateSale: (id: number, saleData: UpdateSaleRequest) => Promise<void>;
  deleteSale: (id: number) => Promise<void>;
  getSalesByFarm: () => Promise<void>;
  getSalesHistory: () => Promise<void>;
  getSalesByAnimal: (animalId: number) => Promise<Sale[]>;
  getSalesByDateRange: (filters: SaleFilters) => Promise<Sale[]>;
  clearError: () => void;
}

const SaleContext = createContext<SaleContextType | undefined>(undefined);

interface SaleProviderProps {
  children: ReactNode;
}

export const SaleProvider: React.FC<SaleProviderProps> = ({ children }) => {
  const [sales, setSales] = useState<Sale[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const clearError = useCallback(() => setError(null), []);

  const createSale = useCallback(async (saleData: CreateSaleRequest) => {
    try {
      setLoading(true);
      setError(null);
      const newSale = await saleService.createSale(saleData);
      setSales(prev => [newSale, ...prev]);
      toast.success('Venda registrada com sucesso!');
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Erro ao registrar venda';
      setError(errorMessage);
      toast.error(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const updateSale = useCallback(async (id: number, saleData: UpdateSaleRequest) => {
    try {
      setLoading(true);
      setError(null);
      const updatedSale = await saleService.updateSale(id, saleData);
      setSales(prev => prev.map(sale => sale.id === id ? updatedSale : sale));
      toast.success('Venda atualizada com sucesso!');
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Erro ao atualizar venda';
      setError(errorMessage);
      toast.error(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const deleteSale = useCallback(async (id: number) => {
    try {
      setLoading(true);
      setError(null);
      await saleService.deleteSale(id);
      setSales(prev => prev.filter(sale => sale.id !== id));
      toast.success('Venda excluída com sucesso!');
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Erro ao excluir venda';
      setError(errorMessage);
      toast.error(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const getSalesByFarm = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const farmSales = await saleService.getSalesByFarm();
      setSales(farmSales);
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Erro ao carregar vendas';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, []);

  const getSalesHistory = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const history = await saleService.getSalesHistory();
      setSales(history);
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Erro ao carregar histórico de vendas';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, []);

  const getSalesByAnimal = useCallback(async (animalId: number): Promise<Sale[]> => {
    try {
      setError(null);
      return await saleService.getSalesByAnimal(animalId);
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Erro ao carregar vendas do animal';
      setError(errorMessage);
      toast.error(errorMessage);
      return [];
    }
  }, []);

  const getSalesByDateRange = useCallback(async (filters: SaleFilters): Promise<Sale[]> => {
    try {
      setError(null);
      return await saleService.getSalesByDateRange(filters);
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Erro ao carregar vendas por período';
      setError(errorMessage);
      toast.error(errorMessage);
      return [];
    }
  }, []);

  const value: SaleContextType = useMemo(() => ({
    sales,
    loading,
    error,
    createSale,
    updateSale,
    deleteSale,
    getSalesByFarm,
    getSalesHistory,
    getSalesByAnimal,
    getSalesByDateRange,
    clearError,
  }), [
    sales,
    loading,
    error,
    createSale,
    updateSale,
    deleteSale,
    getSalesByFarm,
    getSalesHistory,
    getSalesByAnimal,
    getSalesByDateRange,
    clearError,
  ]);

  return (
    <SaleContext.Provider value={value}>
      {children}
    </SaleContext.Provider>
  );
};

export const useSaleContext = (): SaleContextType => {
  const context = useContext(SaleContext);
  if (context === undefined) {
    throw new Error('useSaleContext must be used within a SaleProvider');
  }
  return context;
};
