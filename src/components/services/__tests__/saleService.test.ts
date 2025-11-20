import { describe, it, expect, vi, beforeEach } from 'vitest';
import type { Mock } from 'vitest';

const mockPost = vi.fn();
const mockGet = vi.fn();
const mockPut = vi.fn();
const mockDelete = vi.fn();

vi.mock('../../../config/api', () => ({
  api: {
    post: mockPost,
    get: mockGet,
    put: mockPut,
    delete: mockDelete,
  },
}));

import { saleService } from '../saleService';

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

      (mockPost as Mock).mockResolvedValue(mockResponse);

      const result = await saleService.createSale(saleData);

      expect(mockPost).toHaveBeenCalledWith('/sales', saleData);
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

      (mockGet as Mock).mockResolvedValue(mockResponse);

      const result = await saleService.getSalesByFarm();

      expect(mockGet).toHaveBeenCalledWith('/sales');
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

      (mockGet as Mock).mockResolvedValue(mockResponse);

      const result = await saleService.getSalesHistory();

      expect(mockGet).toHaveBeenCalledWith('/sales/history');
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

      (mockGet as Mock).mockResolvedValue(mockResponse);

      const result = await saleService.getSalesByAnimal(animalId);

      expect(mockGet).toHaveBeenCalledWith(`/animals/${animalId}/sales`);
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

      (mockGet as Mock).mockResolvedValue(mockResponse);

      const result = await saleService.getSalesByDateRange(filters);

      expect(mockGet).toHaveBeenCalledWith(
        '/sales/date-range?start_date=2024-01-01&end_date=2024-01-31'
      );
      expect(result).toEqual(mockResponse.data);
    });

    it('deve buscar vendas apenas com start_date', async () => {
      const filters = {
        start_date: '2024-01-01',
      };

      (mockGet as Mock).mockResolvedValue({ data: [] });

      await saleService.getSalesByDateRange(filters);

      expect(mockGet).toHaveBeenCalledWith(
        '/sales/date-range?start_date=2024-01-01'
      );
    });

    it('deve buscar vendas apenas com end_date', async () => {
      const filters = {
        end_date: '2024-01-31',
      };

      (mockGet as Mock).mockResolvedValue({ data: [] });

      await saleService.getSalesByDateRange(filters);

      expect(mockGet).toHaveBeenCalledWith(
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

      (mockGet as Mock).mockResolvedValue(mockResponse);

      const result = await saleService.getSaleById(saleId);

      expect(mockGet).toHaveBeenCalledWith(`/sales/${saleId}`);
      expect(result).toEqual(mockResponse.data);
    });
  });

  describe('updateSale', () => {
    it('deve atualizar uma venda', async () => {
      const saleId = 1;
      const updateData = {
        buyer_name: 'Updated Buyer',
        price: 1500,
        sale_date: '2024-01-01',
        notes: 'Updated notes',
      };

      const mockResponse = {
        data: {
          id: saleId,
          ...updateData,
        },
      };

      (mockPut as Mock).mockResolvedValue(mockResponse);

      const result = await saleService.updateSale(saleId, updateData);

      expect(mockPut).toHaveBeenCalledWith(`/sales/${saleId}`, updateData);
      expect(result).toEqual(mockResponse.data);
    });
  });

  describe('deleteSale', () => {
    it('deve deletar uma venda', async () => {
      const saleId = 1;

      (mockDelete as Mock).mockResolvedValue(undefined);

      await saleService.deleteSale(saleId);

      expect(mockDelete).toHaveBeenCalledWith(`/sales/${saleId}`);
    });
  });
});

