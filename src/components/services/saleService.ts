import { api } from '../../config/api';
import { Sale, CreateSaleRequest, UpdateSaleRequest, SaleFilters } from '../../types/sale';

export const saleService = {
  async createSale(saleData: CreateSaleRequest): Promise<Sale> {
    const response = await api.post('/sales', saleData);
    return response.data;
  },

  async getSalesByFarm(farmId?: number): Promise<Sale[]> {
    const params = farmId ? `?farmId=${farmId}` : '';
    const response = await api.get(`/sales${params}`);
    return response.data;
  },

  async getSalesHistory(farmId?: number): Promise<Sale[]> {
    const params = farmId ? `?farmId=${farmId}` : '';
    const response = await api.get(`/sales/history${params}`);
    return response.data;
  },

  async getSalesByAnimal(animalId: number, farmId?: number): Promise<Sale[]> {
    const params = farmId ? `?farmId=${farmId}` : '';
    const response = await api.get(`/animals/${animalId}/sales${params}`);
    return response.data;
  },

  async getSalesByDateRange(filters: SaleFilters, farmId?: number): Promise<Sale[]> {
    const params = new URLSearchParams();
    if (filters.start_date) params.append('start_date', filters.start_date);
    if (filters.end_date) params.append('end_date', filters.end_date);
    if (farmId) params.append('farmId', farmId.toString());

    const url = `/sales/date-range?${params.toString()}`;
    const response = await api.get(url);
    return response.data;
  },

  async getSaleById(id: number): Promise<Sale> {
    const response = await api.get(`/sales/${id}`);
    return response.data;
  },

  async updateSale(id: number, saleData: UpdateSaleRequest): Promise<Sale> {
    const response = await api.put(`/sales/${id}`, saleData);
    return response.data;
  },

  async deleteSale(id: number): Promise<void> {
    await api.delete(`/sales/${id}`);
  }
};
