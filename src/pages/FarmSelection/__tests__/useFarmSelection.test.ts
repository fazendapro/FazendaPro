import { renderHook, act } from '@testing-library/react';
import { useFarmSelection } from '../hooks/useFarmSelection';
import { farmSelectionService } from '../services/farm-selection-service';

jest.mock('../services/farm-selection-service');
const mockFarmSelectionService = farmSelectionService as jest.Mocked<typeof farmSelectionService>;

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  useNavigate: () => mockNavigate,
}));

describe('useFarmSelection', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('deve carregar fazendas do usuário ao inicializar', async () => {
    const mockFarms = [
      { ID: 1, CompanyID: 1, Logo: 'logo1.png' },
      { ID: 2, CompanyID: 2, Logo: 'logo2.png' },
    ];

    mockFarmSelectionService.getUserFarms.mockResolvedValue({
      success: true,
      message: 'Fazendas carregadas com sucesso',
      farms: mockFarms,
      auto_select: false,
    });

    const { result } = renderHook(() => useFarmSelection());

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0));
    });

    expect(result.current.farms).toEqual(mockFarms);
    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
  });

  it('deve auto-selecionar quando há apenas uma fazenda', async () => {
    const mockFarms = [
      { ID: 1, CompanyID: 1, Logo: 'logo1.png' },
    ];

    mockFarmSelectionService.getUserFarms.mockResolvedValue({
      success: true,
      message: 'Fazendas carregadas com sucesso',
      farms: mockFarms,
      auto_select: true,
      selected_farm_id: 1,
    });

    const { result } = renderHook(() => useFarmSelection());

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0));
    });

    expect(result.current.autoSelect).toBe(true);
    expect(result.current.selectedFarmId).toBe(1);
  });

  it('deve selecionar uma fazenda quando selectFarm é chamado', async () => {
    const mockFarms = [
      { ID: 1, CompanyID: 1, Logo: 'logo1.png' },
      { ID: 2, CompanyID: 2, Logo: 'logo2.png' },
    ];

    mockFarmSelectionService.getUserFarms.mockResolvedValue({
      success: true,
      message: 'Fazendas carregadas com sucesso',
      farms: mockFarms,
      auto_select: false,
    });

    mockFarmSelectionService.selectFarm.mockResolvedValue({
      success: true,
      message: 'Fazenda selecionada com sucesso',
      farm_id: 2,
    });

    const { result } = renderHook(() => useFarmSelection());

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0));
    });

    await act(async () => {
      result.current.selectFarm(2);
    });

    expect(mockFarmSelectionService.selectFarm).toHaveBeenCalledWith(2);
    expect(mockNavigate).toHaveBeenCalledWith('/');
  });

  it('deve tratar erro ao carregar fazendas', async () => {
    mockFarmSelectionService.getUserFarms.mockRejectedValue(new Error('Erro ao carregar fazendas'));

    const { result } = renderHook(() => useFarmSelection());

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0));
    });

    expect(result.current.error).toBe('Erro ao carregar fazendas');
    expect(result.current.loading).toBe(false);
  });

  it('deve tratar erro ao selecionar fazenda', async () => {
    const mockFarms = [
      { ID: 1, CompanyID: 1, Logo: 'logo1.png' },
    ];

    mockFarmSelectionService.getUserFarms.mockResolvedValue({
      success: true,
      message: 'Fazendas carregadas com sucesso',
      farms: mockFarms,
      auto_select: false,
    });

    mockFarmSelectionService.selectFarm.mockRejectedValue(new Error('Erro ao selecionar fazenda'));

    const { result } = renderHook(() => useFarmSelection());

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0));
    });

    await act(async () => {
      result.current.selectFarm(1);
    });

    expect(result.current.error).toBe('Erro ao selecionar fazenda');
  });
});
