import { useSaleContext } from '../contexts/SaleContext';
import { CreateSaleRequest, UpdateSaleRequest } from '../types/sale';

export const useSale = () => {
  const context = useSaleContext();
  return context;
};

export const useSaleForm = () => {
  const { createSale, updateSale, loading } = useSaleContext();

  const handleCreateSale = async (saleData: CreateSaleRequest) => {
    return await createSale(saleData);
  };

  const handleUpdateSale = async (id: number, saleData: UpdateSaleRequest) => {
    return await updateSale(id, saleData);
  };

  return {
    createSale: handleCreateSale,
    updateSale: handleUpdateSale,
    loading,
  };
};

export const useSaleList = () => {
  const { 
    sales, 
    loading, 
    error, 
    getSalesByFarm, 
    getSalesHistory, 
    getSalesByAnimal, 
    getSalesByDateRange,
    deleteSale,
    clearError 
  } = useSaleContext();

  return {
    sales,
    loading,
    error,
    getSalesByFarm,
    getSalesHistory,
    getSalesByAnimal,
    getSalesByDateRange,
    deleteSale,
    clearError,
  };
};
