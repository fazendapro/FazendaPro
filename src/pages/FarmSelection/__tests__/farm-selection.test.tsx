import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { FarmSelection } from '../farm-selection';
import { useFarmSelection } from '../hooks/useFarmSelection';

jest.mock('../hooks/useFarmSelection');
const mockUseFarmSelection = useFarmSelection as jest.MockedFunction<typeof useFarmSelection>;

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  );
};

describe('FarmSelection', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('deve renderizar a tela de seleção quando há múltiplas fazendas', () => {
    const mockFarms = [
      { id: 1, name: 'Fazenda 1', logo: 'logo1.png' },
      { id: 2, name: 'Fazenda 2', logo: 'logo2.png' },
    ];

    mockUseFarmSelection.mockReturnValue({
      farms: mockFarms,
      loading: false,
      error: null,
      autoSelect: false,
      selectedFarmId: null,
      selectFarm: jest.fn(),
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
      selectFarm: jest.fn(),
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
      selectFarm: jest.fn(),
    });

    renderWithRouter(<FarmSelection />);

    expect(screen.getByText('Erro ao carregar fazendas')).toBeInTheDocument();
  });

  it('deve chamar selectFarm quando uma fazenda é selecionada', async () => {
    const mockSelectFarm = jest.fn();
    const mockFarms = [
      { id: 1, name: 'Fazenda 1', logo: 'logo1.png' },
      { id: 2, name: 'Fazenda 2', logo: 'logo2.png' },
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

    const farmButton = screen.getByText('Fazenda 1');
    fireEvent.click(farmButton);

    await waitFor(() => {
      expect(mockSelectFarm).toHaveBeenCalledWith(1);
    });
  });

  it('deve redirecionar automaticamente quando há apenas uma fazenda', () => {
    const mockFarms = [
      { id: 1, name: 'Fazenda 1', logo: 'logo1.png' },
    ];

    mockUseFarmSelection.mockReturnValue({
      farms: mockFarms,
      loading: false,
      error: null,
      autoSelect: true,
      selectedFarmId: 1,
      selectFarm: jest.fn(),
    });

    renderWithRouter(<FarmSelection />);

    expect(screen.getByText('Redirecionando para sua fazenda...')).toBeInTheDocument();
  });
});
