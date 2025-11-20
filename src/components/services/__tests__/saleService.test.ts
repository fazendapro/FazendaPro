import { describe, it, expect, vi, beforeEach } from 'vitest';

// Mock api antes de importar saleService
vi.mock('../../../config/api', () => ({
  api: {
    post: vi.fn(),
    get: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

import { saleService } from '../saleService';
import { api as mockApi } from '../../../config/api';

describe('saleService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('createSale', () => {
    it('deve criar uma venda com sucesso', async () => {
      const saleData = {
        animal_id: 1,
        sale_date: '2024-01-01',
        price: 1000,
        buyer_name: 'Test Buyer',
        notes: 'Test notes',
      };

      const mockResponse = {
        data: {
          id: 1,
          ...saleData,
        },
      };

      mockApi.post.mockResolvedValue(mockResponse);

      const result = await saleService.createSale(saleData);

      expect(mockApi.post).toHaveBeenCalledWith('/sales', saleData);
      expect(result).toEqual(mockResponse.data);
    });
  });

  describe('getSalesByFarm', () => {
    it('deve buscar vendas da fazenda', async () => {
      const mockResponse = {
        data: [
          { id: 1, animal_id: 1, price: 1000 },
          { id: 2, animal_id: 2, price: 2000 },
        ],
      };

      mockApi.get.mockResolvedValue(mockResponse);

      const result = await saleService.getSalesByFarm();

      expect(mockApi.get).toHaveBeenCalledWith('/sales');
      expect(result).toEqual(mockResponse.data);
    });
  });

  describe('getSalesHistory', () => {
    it('deve buscar histórico de vendas', async () => {
      const mockResponse = {
        data: [
          { id: 1, animal_id: 1, price: 1000 },
        ],
      };

      mockApi.get.mockResolvedValue(mockResponse);

      const result = await saleService.getSalesHistory();

      expect(mockApi.get).toHaveBeenCalledWith('/sales/history');
      expect(result).toEqual(mockResponse.data);
    });
  });

  describe('getSalesByAnimal', () => {
    it('deve buscar vendas de um animal específico', async () => {
      const animalId = 1;
      const mockResponse = {
        data: [
          { id: 1, animal_id: animalId, price: 1000 },
        ],
      };

      mockApi.get.mockResolvedValue(mockResponse);

      const result = await saleService.getSalesByAnimal(animalId);

      expect(mockApi.get).toHaveBeenCalledWith(`/animals/${animalId}/sales`);
      expect(result).toEqual(mockResponse.data);
    });
  });

  describe('getSalesByDateRange', () => {
    it('deve buscar vendas por intervalo de datas', async () => {
      const filters = {
        start_date: '2024-01-01',
        end_date: '2024-01-31',
      };

      const mockResponse = {
        data: [
          { id: 1, sale_date: '2024-01-15', price: 1000 },
        ],
      };

      mockApi.get.mockResolvedValue(mockResponse);

      const result = await saleService.getSalesByDateRange(filters);

      expect(mockApi.get).toHaveBeenCalledWith(
        '/sales/date-range?start_date=2024-01-01&end_date=2024-01-31'
      );
      expect(result).toEqual(mockResponse.data);
    });

    it('deve buscar vendas apenas com start_date', async () => {
      const filters = {
        start_date: '2024-01-01',
      };

      mockApi.get.mockResolvedValue({ data: [] });

      await saleService.getSalesByDateRange(filters);

      expect(mockApi.get).toHaveBeenCalledWith(
        '/sales/date-range?start_date=2024-01-01'
      );
    });

    it('deve buscar vendas apenas com end_date', async () => {
      const filters = {
        end_date: '2024-01-31',
      };

      mockApi.get.mockResolvedValue({ data: [] });

      await saleService.getSalesByDateRange(filters);

      expect(mockApi.get).toHaveBeenCalledWith(
        '/sales/date-range?end_date=2024-01-31'
      );
    });
  });

  describe('getSaleById', () => {
    it('deve buscar venda por ID', async () => {
      const saleId = 1;
      const mockResponse = {
        data: {
          id: saleId,
          animal_id: 1,
          price: 1000,
        },
      };

      mockApi.get.mockResolvedValue(mockResponse);

      const result = await saleService.getSaleById(saleId);

      expect(mockApi.get).toHaveBeenCalledWith(`/sales/${saleId}`);
      expect(result).toEqual(mockResponse.data);
    });
  });

  describe('updateSale', () => {
    it('deve atualizar uma venda', async () => {
      const saleId = 1;
      const updateData = {
        price: 1500,
        notes: 'Updated notes',
      };

      const mockResponse = {
        data: {
          id: saleId,
          ...updateData,
        },
      };

      mockApi.put.mockResolvedValue(mockResponse);

      const result = await saleService.updateSale(saleId, updateData);

      expect(mockApi.put).toHaveBeenCalledWith(`/sales/${saleId}`, updateData);
      expect(result).toEqual(mockResponse.data);
    });
  });

  describe('deleteSale', () => {
    it('deve deletar uma venda', async () => {
      const saleId = 1;

      mockApi.delete.mockResolvedValue(undefined);

      await saleService.deleteSale(saleId);

      expect(mockApi.delete).toHaveBeenCalledWith(`/sales/${saleId}`);
    });
  });
});

