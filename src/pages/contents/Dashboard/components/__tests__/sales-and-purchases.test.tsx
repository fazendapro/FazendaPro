import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { SalesAndPurchases } from '../sales-and-purchases';
import { useMonthlySalesAndPurchases } from '../../hooks/useMonthlySalesAndPurchases';

vi.mock('../../hooks/useMonthlySalesAndPurchases', () => ({
  useMonthlySalesAndPurchases: vi.fn()
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('react-chartjs-2', () => ({
  Bar: ({ data }: { data: unknown }) => <div data-testid="bar-chart">{JSON.stringify(data)}</div>,
}));

vi.mock('chart.js', () => ({
  Chart: {
    register: vi.fn(),
  },
  CategoryScale: {},
  LinearScale: {},
  BarElement: {},
  Title: {},
  Tooltip: {},
  Legend: {},
}));

const mockUseMonthlySalesAndPurchases = vi.mocked(useMonthlySalesAndPurchases);

describe('SalesAndPurchases', () => {
  const mockSalesData = {
    sales: [
      { month: 'Jan', year: 2024, sales: 15000, count: 5 },
      { month: 'Fev', year: 2024, sales: 20000, count: 7 },
      { month: 'Mar', year: 2024, sales: 18000, count: 6 },
    ],
    purchases: [
      { month: 'Jan', year: 2024, total: 0 },
      { month: 'Fev', year: 2024, total: 0 },
      { month: 'Mar', year: 2024, total: 0 },
    ],
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('deve renderizar o título do card com tooltip', () => {
    mockUseMonthlySalesAndPurchases.mockReturnValue({
      data: mockSalesData,
      loading: false,
      error: null,
      getMonthlySalesAndPurchases: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<SalesAndPurchases />);

    expect(screen.getByText('dashboard.salesAndPurchases')).toBeInTheDocument();
    const infoIcon = document.querySelector('.anticon-info-circle');
    expect(infoIcon).toBeInTheDocument();
  });

  it('deve renderizar gráfico de barras quando dados estão disponíveis', () => {
    mockUseMonthlySalesAndPurchases.mockReturnValue({
      data: mockSalesData,
      loading: false,
      error: null,
      getMonthlySalesAndPurchases: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<SalesAndPurchases />);

    const chart = screen.getByTestId('bar-chart');
    expect(chart).toBeInTheDocument();
  });

  it('deve renderizar loading quando dados estão carregando', () => {
    mockUseMonthlySalesAndPurchases.mockReturnValue({
      data: { sales: [], purchases: [] },
      loading: true,
      error: null,
      getMonthlySalesAndPurchases: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<SalesAndPurchases />);

    const spinner = document.querySelector('.ant-spin');
    expect(spinner).toBeInTheDocument();
    
    expect(screen.queryByTestId('bar-chart')).not.toBeInTheDocument();
  });

  it('deve renderizar erro quando há erro na requisição', () => {
    const errorMessage = 'Erro ao buscar dados';
    mockUseMonthlySalesAndPurchases.mockReturnValue({
      data: { sales: [], purchases: [] },
      loading: false,
      error: errorMessage,
      getMonthlySalesAndPurchases: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<SalesAndPurchases />);

    expect(screen.getByText(errorMessage)).toBeInTheDocument();
    expect(screen.queryByTestId('bar-chart')).not.toBeInTheDocument();
  });

  it('deve chamar useMonthlySalesAndPurchases com parâmetros corretos', () => {
    mockUseMonthlySalesAndPurchases.mockReturnValue({
      data: mockSalesData,
      loading: false,
      error: null,
      getMonthlySalesAndPurchases: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<SalesAndPurchases />);

    expect(mockUseMonthlySalesAndPurchases).toHaveBeenCalledWith(undefined, 12);
  });

  it('deve renderizar gráfico com dados vazios quando não há dados', () => {
    mockUseMonthlySalesAndPurchases.mockReturnValue({
      data: { sales: [], purchases: [] },
      loading: false,
      error: null,
      getMonthlySalesAndPurchases: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<SalesAndPurchases />);

    const chart = screen.getByTestId('bar-chart');
    expect(chart).toBeInTheDocument();
  });

  it('deve renderizar tooltip com mensagem explicativa', () => {
    mockUseMonthlySalesAndPurchases.mockReturnValue({
      data: mockSalesData,
      loading: false,
      error: null,
      getMonthlySalesAndPurchases: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<SalesAndPurchases />);

    const infoIcon = document.querySelector('.anticon-info-circle');
    expect(infoIcon).toBeInTheDocument();
  });
});

