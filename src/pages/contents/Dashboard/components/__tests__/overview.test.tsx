import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { Overview } from '../overview';
import { useOverviewStats } from '../../hooks/useOverviewStats';

vi.mock('../../hooks/useOverviewStats', () => ({
  useOverviewStats: vi.fn()
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

const mockUseOverviewStats = vi.mocked(useOverviewStats);

describe('Overview', () => {
  const mockStats = {
    males_count: 10,
    females_count: 15,
    total_sold: 5,
    total_revenue: 15000.50
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('deve renderizar o título do card', () => {
    mockUseOverviewStats.mockReturnValue({
      stats: mockStats,
      loading: false,
      error: null,
      getOverviewStats: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<Overview />);

    expect(screen.getByText('dashboard.salesOverview')).toBeInTheDocument();
  });

  it('deve renderizar todas as estatísticas quando dados estão disponíveis', () => {
    mockUseOverviewStats.mockReturnValue({
      stats: mockStats,
      loading: false,
      error: null,
      getOverviewStats: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<Overview />);

    expect(screen.getByText('dashboard.machos')).toBeInTheDocument();
    expect(screen.getByText('10')).toBeInTheDocument();
    expect(screen.getByText('dashboard.femeas')).toBeInTheDocument();
    expect(screen.getByText('15')).toBeInTheDocument();
    expect(screen.getByText('dashboard.sales')).toBeInTheDocument();
    expect(screen.getByText('5')).toBeInTheDocument();
    expect(screen.getByText('dashboard.revenue')).toBeInTheDocument();
  });

  it('deve renderizar loading quando dados estão carregando', () => {
    mockUseOverviewStats.mockReturnValue({
      stats: { males_count: 0, females_count: 0, total_sold: 0, total_revenue: 0 },
      loading: true,
      error: null,
      getOverviewStats: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<Overview />);

    const spinner = document.querySelector('.ant-spin');
    expect(spinner).toBeInTheDocument();
  });

  it('deve renderizar erro quando há erro na requisição', () => {
    const errorMessage = 'Erro ao buscar dados';
    mockUseOverviewStats.mockReturnValue({
      stats: { males_count: 0, females_count: 0, total_sold: 0, total_revenue: 0 },
      loading: false,
      error: errorMessage,
      getOverviewStats: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<Overview />);

    expect(screen.getByText(errorMessage)).toBeInTheDocument();
  });

  it('deve chamar useOverviewStats sem parâmetros', () => {
    mockUseOverviewStats.mockReturnValue({
      stats: mockStats,
      loading: false,
      error: null,
      getOverviewStats: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<Overview />);

    expect(mockUseOverviewStats).toHaveBeenCalledWith();
  });

  it('deve formatar receita como moeda', () => {
    mockUseOverviewStats.mockReturnValue({
      stats: mockStats,
      loading: false,
      error: null,
      getOverviewStats: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<Overview />);

    expect(screen.getByText('R$ 15.001')).toBeInTheDocument();
  });

  it('deve renderizar valores zero quando não há dados', () => {
    mockUseOverviewStats.mockReturnValue({
      stats: { males_count: 0, females_count: 0, total_sold: 0, total_revenue: 0 },
      loading: false,
      error: null,
      getOverviewStats: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<Overview />);

    expect(screen.getAllByText('0')).toHaveLength(3);
  });
});

