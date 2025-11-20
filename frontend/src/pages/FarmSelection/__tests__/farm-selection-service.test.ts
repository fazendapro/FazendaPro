import { describe, it, expect, vi, beforeEach } from 'vitest';
import { farmSelectionService } from '../services/farm-selection-service';
import { api } from '../../../components/services/axios/api';

vi.mock('../../../components/services/axios/api');
const mockApi = api as ReturnType<typeof vi.fn>;

const mockApiInstance = {
  get: vi.fn(),
  post: vi.fn()
};

mockApi.mockReturnValue(mockApiInstance as unknown as ReturnType<typeof api>);

describe('farmSelectionService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getUserFarms', () => {
    it('deve retornar fazendas do usuário com sucesso', async () => {
      const mockResponse = {
        success: true,
        farms: [
          { ID: 1, CompanyID: 1, Logo: 'logo1.png' },
          { ID: 2, CompanyID: 2, Logo: 'logo2.png' },
        ],
        auto_select: false,
      };

      mockApiInstance.get.mockResolvedValue({ data: mockResponse });

      const result = await farmSelectionService.getUserFarms();

      expect(mockApiInstance.get).toHaveBeenCalledWith('/farms/user');
      expect(result).toEqual(mockResponse);
    });

    it('deve retornar auto-seleção quando há apenas uma fazenda', async () => {
      const mockResponse = {
        success: true,
        farms: [
          { ID: 1, CompanyID: 1, Logo: 'logo1.png' },
        ],
        auto_select: true,
        selected_farm_id: 1,
      };

      mockApiInstance.get.mockResolvedValue({ data: mockResponse });

      const result = await farmSelectionService.getUserFarms();

      expect(result.auto_select).toBe(true);
      expect(result.selected_farm_id).toBe(1);
    });

    it('deve tratar erro quando a API retorna erro', async () => {
      const mockError = new Error('Erro ao buscar fazendas');
      mockApiInstance.get.mockRejectedValue(mockError);

      await expect(farmSelectionService.getUserFarms()).rejects.toThrow('Erro ao buscar fazendas');
    });
  });

  describe('selectFarm', () => {
    it('deve selecionar uma fazenda com sucesso', async () => {
      const mockResponse = {
        success: true,
        farm_id: 2,
        message: 'Fazenda selecionada com sucesso',
      };

      mockApiInstance.post.mockResolvedValue({ data: mockResponse });

      const result = await farmSelectionService.selectFarm(2);

      expect(mockApiInstance.post).toHaveBeenCalledWith('/farms/select', { farm_id: 2 });
      expect(result).toEqual(mockResponse);
    });

    it('deve tratar erro quando a fazenda não é encontrada', async () => {
      const mockError = new Error('Fazenda não encontrada');
      mockApiInstance.post.mockRejectedValue(mockError);

      await expect(farmSelectionService.selectFarm(999)).rejects.toThrow('Erro ao selecionar fazenda');
    });
  });
});
