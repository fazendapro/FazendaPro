import { farmSelectionService } from '../services/farm-selection-service';
import { apiClient } from '../../../components/services/axios/api-client';

jest.mock('../../../components/services/axios/api-client');
const mockApiClient = apiClient as jest.Mocked<typeof apiClient>;

describe('farmSelectionService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('getUserFarms', () => {
    it('deve retornar fazendas do usuário com sucesso', async () => {
      const mockResponse = {
        success: true,
        farms: [
          { id: 1, name: 'Fazenda 1', logo: 'logo1.png' },
          { id: 2, name: 'Fazenda 2', logo: 'logo2.png' },
        ],
        auto_select: false,
      };

      mockApiClient.get.mockResolvedValue({ data: mockResponse });

      const result = await farmSelectionService.getUserFarms();

      expect(mockApiClient.get).toHaveBeenCalledWith('/api/v1/farms/user');
      expect(result).toEqual(mockResponse);
    });

    it('deve retornar auto-seleção quando há apenas uma fazenda', async () => {
      const mockResponse = {
        success: true,
        farms: [
          { id: 1, name: 'Fazenda 1', logo: 'logo1.png' },
        ],
        auto_select: true,
        selected_farm_id: 1,
      };

      mockApiClient.get.mockResolvedValue({ data: mockResponse });

      const result = await farmSelectionService.getUserFarms();

      expect(result.auto_select).toBe(true);
      expect(result.selected_farm_id).toBe(1);
    });

    it('deve tratar erro quando a API retorna erro', async () => {
      const mockError = new Error('Erro ao buscar fazendas');
      mockApiClient.get.mockRejectedValue(mockError);

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

      mockApiClient.post.mockResolvedValue({ data: mockResponse });

      const result = await farmSelectionService.selectFarm(2);

      expect(mockApiClient.post).toHaveBeenCalledWith('/api/v1/farms/select', { farm_id: 2 });
      expect(result).toEqual(mockResponse);
    });

    it('deve tratar erro quando a fazenda não é encontrada', async () => {
      const mockError = new Error('Fazenda não encontrada');
      mockApiClient.post.mockRejectedValue(mockError);

      await expect(farmSelectionService.selectFarm(999)).rejects.toThrow('Fazenda não encontrada');
    });
  });
});
