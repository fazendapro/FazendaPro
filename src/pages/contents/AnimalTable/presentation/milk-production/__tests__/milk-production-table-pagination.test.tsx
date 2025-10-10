import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import { MilkProductionTable } from '../milk-production-table';
import { useMilkProduction } from '../../../hooks/useMilkProduction';
import { useFarm } from '../../../../../../hooks/useFarm';
import { useResponsive } from '../../../../../../hooks';

vi.mock('../../../hooks/useMilkProduction');
vi.mock('../../../../../../hooks/useFarm');
vi.mock('../../../../../../hooks');

const mockUseMilkProduction = useMilkProduction as any;
const mockUseFarm = useFarm as any;
const mockUseResponsive = useResponsive as any;

describe('MilkProductionTable Pagination Integration', () => {
  const mockMilkProductions = Array.from({ length: 30 }, (_, i) => ({
    id: i + 1,
    animal: {
      id: i + 1,
      animal_name: `Animal ${i + 1}`,
      ear_tag_number_local: i + 1
    },
    liters: 10 + (i % 5),
    date: `2024-01-${String(i % 28 + 1).padStart(2, '0')}`
  }));

  beforeEach(() => {
    vi.clearAllMocks();
    
    mockUseFarm.mockReturnValue({
      farm: { id: 1, name: 'Test Farm' }
    });
    
    mockUseMilkProduction.mockReturnValue({
      milkProductions: mockMilkProductions,
      loading: false,
      error: null,
      refetch: vi.fn()
    });
    
    mockUseResponsive.mockReturnValue({
      isMobile: false,
      isTablet: false,
      isDesktop: true
    });
  });

  describe('Paginação básica', () => {
    it('deve renderizar primeira página com 10 itens por padrão', () => {
      render(<MilkProductionTable />);
      
      expect(screen.getByText('Animal 1')).toBeInTheDocument();
      expect(screen.getByText('Animal 10')).toBeInTheDocument();
      expect(screen.queryByText('Animal 11')).not.toBeInTheDocument();
      
      expect(screen.getByText('1-10 de 30 registros')).toBeInTheDocument();
    });

    it('deve navegar para próxima página', async () => {
      render(<MilkProductionTable />);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      fireEvent.click(nextButton);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 11')).toBeInTheDocument();
        expect(screen.getByText('Animal 20')).toBeInTheDocument();
        expect(screen.queryByText('Animal 1')).not.toBeInTheDocument();
      });
    });

    it('deve navegar para última página', async () => {
      render(<MilkProductionTable />);
      
      const page3Button = screen.getByRole('button', { name: /3/i });
      fireEvent.click(page3Button);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 21')).toBeInTheDocument();
        expect(screen.getByText('Animal 30')).toBeInTheDocument();
        expect(screen.queryByText('Animal 20')).not.toBeInTheDocument();
      });
    });
  });

  describe('Filtros de período', () => {
    it('deve manter paginação com filtros aplicados', async () => {
      const { rerender } = render(<MilkProductionTable />);
      
      const periodSelect = screen.getByRole('combobox');
      fireEvent.click(periodSelect);
      
      const weekOption = screen.getByText('Semana');
      fireEvent.click(weekOption);
      
      mockUseMilkProduction.mockReturnValue({
        milkProductions: mockMilkProductions.slice(0, 7), // Apenas 7 registros
        loading: false,
        error: null,
        refetch: vi.fn()
      });
      
      rerender(<MilkProductionTable />);
      
      await waitFor(() => {
        // Deve mostrar apenas 7 registros
        expect(screen.getByText('1-7 de 7 registros')).toBeInTheDocument();
      });
    });

    it('deve resetar paginação ao alterar filtros', async () => {
      render(<MilkProductionTable />);
      
      // Navegar para página 2
      const nextButton = screen.getByRole('button', { name: /next/i });
      fireEvent.click(nextButton);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 11')).toBeInTheDocument();
      });
      
      // Alterar filtro
      const periodSelect = screen.getByRole('combobox');
      fireEvent.click(periodSelect);
      
      const monthOption = screen.getByText('Mês');
      fireEvent.click(monthOption);
      
      // Simular dados filtrados
      mockUseMilkProduction.mockReturnValue({
        milkProductions: mockMilkProductions.slice(0, 15), // 15 registros
        loading: false,
        error: null,
        refetch: vi.fn()
      });
      
      const { rerender } = render(<MilkProductionTable />);
      
      await waitFor(() => {
        // Deve voltar para primeira página
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
        expect(screen.getByText('1-10 de 15 registros')).toBeInTheDocument();
      });
    });
  });

  describe('Filtros de data', () => {
    it('deve aplicar filtro de intervalo de datas', async () => {
      render(<MilkProductionTable />);
      
      // Simular seleção de intervalo de datas
      const dateRangePicker = screen.getByPlaceholderText('Data inicial');
      fireEvent.click(dateRangePicker);
      
      // Simular dados filtrados por data
      mockUseMilkProduction.mockReturnValue({
        milkProductions: mockMilkProductions.slice(0, 12), // 12 registros
        loading: false,
        error: null,
        refetch: vi.fn()
      });
      
      const { rerender } = render(<MilkProductionTable />);
      
      await waitFor(() => {
        expect(screen.getByText('1-10 de 12 registros')).toBeInTheDocument();
      });
    });
  });

  describe('Comportamento responsivo', () => {
    it('deve usar paginação simplificada em mobile', () => {
      mockUseResponsive.mockReturnValue({
        isMobile: true,
        isTablet: false,
        isDesktop: false
      });

      render(<MilkProductionTable />);
      
      // Em mobile, deve mostrar apenas 5 itens por página
      expect(screen.getByText('Animal 1')).toBeInTheDocument();
      expect(screen.getByText('Animal 5')).toBeInTheDocument();
      expect(screen.queryByText('Animal 6')).not.toBeInTheDocument();
      
      // Não deve mostrar controles de tamanho
      expect(screen.queryByRole('combobox')).not.toBeInTheDocument();
    });

    it('deve usar paginação intermediária em tablet', () => {
      mockUseResponsive.mockReturnValue({
        isMobile: false,
        isTablet: true,
        isDesktop: false
      });

      render(<MilkProductionTable />);
      
      // Em tablet, deve mostrar 8 itens por página
      expect(screen.getByText('Animal 1')).toBeInTheDocument();
      expect(screen.getByText('Animal 8')).toBeInTheDocument();
      expect(screen.queryByText('Animal 9')).not.toBeInTheDocument();
    });
  });

  describe('Estados de carregamento', () => {
    it('deve mostrar loading sem paginação', () => {
      mockUseMilkProduction.mockReturnValue({
        milkProductions: [],
        loading: true,
        error: null,
        refetch: vi.fn()
      });

      render(<MilkProductionTable />);
      
      expect(screen.getByRole('img', { name: /loading/i })).toBeInTheDocument();
      expect(screen.queryByText(/registros/)).not.toBeInTheDocument();
    });

    it('deve mostrar erro sem paginação', () => {
      mockUseMilkProduction.mockReturnValue({
        milkProductions: [],
        loading: false,
        error: 'Erro ao carregar dados',
        refetch: vi.fn()
      });

      render(<MilkProductionTable />);
      
      expect(screen.getByText('Erro ao carregar dados')).toBeInTheDocument();
      expect(screen.queryByText(/registros/)).not.toBeInTheDocument();
    });
  });

  describe('Dados vazios', () => {
    it('deve lidar com lista vazia', () => {
      mockUseMilkProduction.mockReturnValue({
        milkProductions: [],
        loading: false,
        error: null,
        refetch: vi.fn()
      });

      render(<MilkProductionTable />);
      
      expect(screen.getByText('0-0 de 0 registros')).toBeInTheDocument();
    });
  });

  describe('Interação com botões de ação', () => {
    it('deve manter paginação ao clicar em editar', async () => {
      const onEditProduction = vi.fn();
      render(<MilkProductionTable onEditProduction={onEditProduction} />);
      
      // Navegar para página 2
      const nextButton = screen.getByRole('button', { name: /next/i });
      fireEvent.click(nextButton);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 11')).toBeInTheDocument();
      });
      
      // Clicar em editar
      const editButton = screen.getAllByRole('button', { name: /edit/i })[0];
      fireEvent.click(editButton);
      
      // Deve manter a página atual
      expect(screen.getByText('Animal 11')).toBeInTheDocument();
      expect(onEditProduction).toHaveBeenCalled();
    });
  });
});
