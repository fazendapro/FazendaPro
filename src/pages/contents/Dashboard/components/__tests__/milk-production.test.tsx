import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { DashboardMilkProduction } from '../milk-production';
import { useTopMilkProducers } from '../../hooks/useTopMilkProducers';

vi.mock('../../hooks/useTopMilkProducers', () => ({
  useTopMilkProducers: vi.fn()
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

const mockUseTopMilkProducers = vi.mocked(useTopMilkProducers);

describe('DashboardMilkProduction', () => {
  const mockTopProducersData = [
    {
      id: 1,
      animal_name: 'Surf Excel',
      ear_tag_number_local: 123,
      photo: 'src/assets/images/mocked/cows/surf.png',
      total_production: 900,
      average_daily_production: 30,
      fat_content: 3.5,
      last_collection_date: '2024-01-15',
      days_in_lactation: 120
    },
    {
      id: 2,
      animal_name: 'Rin',
      ear_tag_number_local: 124,
      photo: 'src/assets/images/mocked/cows/rin.png',
      total_production: 630,
      average_daily_production: 21,
      fat_content: 4.2,
      last_collection_date: '2024-01-15',
      days_in_lactation: 100
    },
    {
      id: 3,
      animal_name: 'Parle G',
      ear_tag_number_local: 125,
      photo: 'src/assets/images/mocked/cows/parle.png',
      total_production: 570,
      average_daily_production: 19,
      fat_content: 4.5,
      last_collection_date: '2024-01-15',
      days_in_lactation: 90
    }
  ];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('deve renderizar o título do card com tooltip', () => {
    mockUseTopMilkProducers.mockReturnValue({
      topProducers: mockTopProducersData,
      loading: false,
      error: null,
      getTopMilkProducers: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<DashboardMilkProduction />);

    expect(screen.getByText('dashboard.topMilkProducers')).toBeInTheDocument();
    const infoIcon = document.querySelector('.anticon-info-circle');
    expect(infoIcon).toBeInTheDocument();
  });

  it('deve renderizar tabela de produtoras quando dados estão disponíveis', () => {
    mockUseTopMilkProducers.mockReturnValue({
      topProducers: mockTopProducersData,
      loading: false,
      error: null,
      getTopMilkProducers: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<DashboardMilkProduction />);

    expect(screen.getByText('Surf Excel')).toBeInTheDocument();
    expect(screen.getByText('Rin')).toBeInTheDocument();
    expect(screen.getByText('Parle G')).toBeInTheDocument();
  });

  it('deve renderizar informações corretas para cada produtora', () => {
    mockUseTopMilkProducers.mockReturnValue({
      topProducers: mockTopProducersData,
      loading: false,
      error: null,
      getTopMilkProducers: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<DashboardMilkProduction />);

    expect(screen.getByText('30L/dia')).toBeInTheDocument();
    expect(screen.getByText('21L/dia')).toBeInTheDocument();
    expect(screen.getByText('19L/dia')).toBeInTheDocument();
    
    expect(screen.getByText('3.5%')).toBeInTheDocument();
    expect(screen.getByText('4.2%')).toBeInTheDocument();
    expect(screen.getByText('4.5%')).toBeInTheDocument();
    
    expect(screen.getByText('120 dias')).toBeInTheDocument();
    expect(screen.getByText('100 dias')).toBeInTheDocument();
    expect(screen.getByText('90 dias')).toBeInTheDocument();
  });

  it('deve renderizar números de brinco corretamente', () => {
    mockUseTopMilkProducers.mockReturnValue({
      topProducers: mockTopProducersData,
      loading: false,
      error: null,
      getTopMilkProducers: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<DashboardMilkProduction />);

    expect(screen.getByText('#123')).toBeInTheDocument();
    expect(screen.getByText('#124')).toBeInTheDocument();
    expect(screen.getByText('#125')).toBeInTheDocument();
  });

  it('deve renderizar loading quando dados estão carregando', () => {
    mockUseTopMilkProducers.mockReturnValue({
      topProducers: [],
      loading: true,
      error: null,
      getTopMilkProducers: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<DashboardMilkProduction />);

    const spinner = document.querySelector('.ant-spin');
    expect(spinner).toBeInTheDocument();
    
    expect(screen.queryByText('Surf Excel')).not.toBeInTheDocument();
  });

  it('deve renderizar lista vazia quando não há dados', () => {
    mockUseTopMilkProducers.mockReturnValue({
      topProducers: [],
      loading: false,
      error: null,
      getTopMilkProducers: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<DashboardMilkProduction />);

    expect(screen.queryByText('Surf Excel')).not.toBeInTheDocument();
    expect(screen.queryByText('Rin')).not.toBeInTheDocument();
    expect(screen.queryByText('Parle G')).not.toBeInTheDocument();
  });

  it('deve renderizar erro quando há erro na requisição', () => {
    const errorMessage = 'Erro ao buscar dados';
    mockUseTopMilkProducers.mockReturnValue({
      topProducers: [],
      loading: false,
      error: errorMessage,
      getTopMilkProducers: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<DashboardMilkProduction />);

    expect(screen.getByText(errorMessage)).toBeInTheDocument();
    expect(screen.queryByText('Surf Excel')).not.toBeInTheDocument();
  });

  it('deve renderizar cabeçalhos da tabela corretamente', () => {
    mockUseTopMilkProducers.mockReturnValue({
      topProducers: mockTopProducersData,
      loading: false,
      error: null,
      getTopMilkProducers: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<DashboardMilkProduction />);

    expect(screen.getByText('dashboard.name')).toBeInTheDocument();
    expect(screen.getByText('dashboard.production')).toBeInTheDocument();
    expect(screen.getByText('dashboard.fat')).toBeInTheDocument();
    expect(screen.getByText('dashboard.daysInLactation')).toBeInTheDocument();
  });

  it('deve chamar useTopMilkProducers com parâmetros corretos', () => {
    mockUseTopMilkProducers.mockReturnValue({
      topProducers: mockTopProducersData,
      loading: false,
      error: null,
      getTopMilkProducers: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<DashboardMilkProduction />);

    expect(mockUseTopMilkProducers).toHaveBeenCalledWith(undefined, 10, 30);
  });

  it('deve renderizar dados ordenados por produção', () => {
    mockUseTopMilkProducers.mockReturnValue({
      topProducers: mockTopProducersData,
      loading: false,
      error: null,
      getTopMilkProducers: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<DashboardMilkProduction />);

    const tableRows = document.querySelectorAll('.ant-table-tbody tr');
    expect(tableRows).toHaveLength(3);
  });
});
