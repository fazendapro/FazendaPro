import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { FarmSelection } from '../farm-selection';
import { useFarmSelection } from '../hooks/useFarmSelection';

vi.mock('../hooks/useFarmSelection');
const mockUseFarmSelection = useFarmSelection as ReturnType<typeof vi.fn>;

vi.mock('../../hooks/use-responsive', () => ({
  useResponsive: () => ({
    isMobile: false,
    isTablet: false,
    isDesktop: true,
    isLargeDesktop: false,
    screenWidth: 1024,
  }),
}));

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  );
};

describe('FarmSelection', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('deve renderizar a tela de seleção quando há múltiplas fazendas', () => {
    const mockFarms = [
      { ID: 1, CompanyID: 1, Logo: 'logo1.png' },
      { ID: 2, CompanyID: 2, Logo: 'logo2.png' },
    ];

    mockUseFarmSelection.mockReturnValue({
      farms: mockFarms,
      loading: false,
      error: null,
      autoSelect: false,
      selectedFarmId: null,
      selectFarm: vi.fn(),
    });

    renderWithRouter(<FarmSelection />);

    expect(screen.getByText('Selecione uma Fazenda')).toBeInTheDocument();
    expect(screen.getByText('Fazenda 1')).toBeInTheDocument();
    expect(screen.getByText('Fazenda 2')).toBeInTheDocument();
  });

  it('deve mostrar loading quando está carregando', () => {
    mockUseFarmSelection.mockReturnValue({
      farms: [],
      loading: true,
      error: null,
      autoSelect: false,
      selectedFarmId: null,
      selectFarm: vi.fn(),
    });

    renderWithRouter(<FarmSelection />);

    expect(screen.getByText('Carregando fazendas...')).toBeInTheDocument();
  });

  it('deve mostrar erro quando há erro ao carregar fazendas', () => {
    mockUseFarmSelection.mockReturnValue({
      farms: [],
      loading: false,
      error: 'Erro ao carregar fazendas',
      autoSelect: false,
      selectedFarmId: null,
      selectFarm: vi.fn(),
    });

    renderWithRouter(<FarmSelection />);

    expect(screen.getByText('Erro ao carregar fazendas')).toBeInTheDocument();
  });

  it('deve chamar selectFarm quando uma fazenda é selecionada', async () => {
    const mockSelectFarm = vi.fn();
    const mockFarms = [
      { ID: 1, CompanyID: 1, Logo: 'logo1.png' },
      { ID: 2, CompanyID: 2, Logo: 'logo2.png' },
    ];

    mockUseFarmSelection.mockReturnValue({
      farms: mockFarms,
      loading: false,
      error: null,
      autoSelect: false,
      selectedFarmId: null,
      selectFarm: mockSelectFarm,
    });

    renderWithRouter(<FarmSelection />);

    // O Card tem onClick que chama selectFarm diretamente
    // Vamos encontrar o Card que contém "Fazenda 1" e clicar nele
    const farmTitle = screen.getByText('Fazenda 1');
    // Encontrar o Card pai (pode estar em vários níveis acima)
    const cardElement = farmTitle.closest('[class*="ant-card"]') || farmTitle.parentElement?.parentElement;
    
    if (cardElement) {
      fireEvent.click(cardElement as HTMLElement);
    } else {
      // Fallback: clicar diretamente no título
      fireEvent.click(farmTitle);
    }

    await waitFor(() => {
      expect(mockSelectFarm).toHaveBeenCalledWith(1);
    });
  });

  it('deve redirecionar automaticamente quando há apenas uma fazenda', () => {
    const mockFarms = [
      { ID: 1, CompanyID: 1, Logo: 'logo1.png' },
    ];

    mockUseFarmSelection.mockReturnValue({
      farms: mockFarms,
      loading: false,
      error: null,
      autoSelect: true,
      selectedFarmId: 1,
      selectFarm: vi.fn(),
    });

    renderWithRouter(<FarmSelection />);

    expect(screen.getByText('Redirecionando para sua fazenda...')).toBeInTheDocument();
  });
});
