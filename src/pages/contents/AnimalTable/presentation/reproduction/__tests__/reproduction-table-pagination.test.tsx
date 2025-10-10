import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import { ReproductionTable } from '../reproduction-table';
import { useReproduction } from '../../../hooks/useReproduction';
import { useFarm } from '../../../../../../hooks/useFarm';
import { useResponsive } from '../../../../../../hooks';

// Mocks
vi.mock('../../../hooks/useReproduction');
vi.mock('../../../../../../hooks/useFarm');
vi.mock('../../../../../../hooks');

const mockUseReproduction = useReproduction as any;
const mockUseFarm = useFarm as any;
const mockUseResponsive = useResponsive as any;

describe('ReproductionTable Pagination Integration', () => {
  const mockReproductions = Array.from({ length: 35 }, (_, i) => ({
    id: i + 1,
    animal_id: i + 1,
    animal_name: `Animal ${i + 1}`,
    ear_tag: i + 1,
    current_phase: (i % 4) + 1,
    insemination_date: `2024-01-${String(i % 28 + 1).padStart(2, '0')}`,
    pregnancy_date: `2024-02-${String(i % 28 + 1).padStart(2, '0')}`,
    expected_birth_date: `2024-10-${String(i % 28 + 1).padStart(2, '0')}`,
    veterinary_confirmation: i % 2 === 0
  }));

  beforeEach(() => {
    vi.clearAllMocks();
    
    mockUseFarm.mockReturnValue({
      farm: { id: 1, name: 'Test Farm' }
    });
    
    mockUseReproduction.mockReturnValue({
      getReproductionsByFarm: vi.fn().mockResolvedValue(mockReproductions),
      deleteReproduction: vi.fn().mockResolvedValue(true),
      loading: false,
      error: null
    });
    
    mockUseResponsive.mockReturnValue({
      isMobile: false,
      isTablet: false,
      isDesktop: true
    });
  });

  describe('Paginação básica', () => {
    it('deve renderizar primeira página com 10 itens por padrão', async () => {
      render(<ReproductionTable onAddReproduction={vi.fn()} onEditReproduction={vi.fn()} />);
      
      await waitFor(() => {
        // Deve mostrar apenas os primeiros 10 registros
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
        expect(screen.getByText('Animal 10')).toBeInTheDocument();
        expect(screen.queryByText('Animal 11')).not.toBeInTheDocument();
        
        // Deve mostrar contador correto
        expect(screen.getByText('1-10 de 35 registros')).toBeInTheDocument();
      });
    });

    it('deve navegar para próxima página', async () => {
      render(<ReproductionTable onAddReproduction={vi.fn()} onEditReproduction={vi.fn()} />);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
      });
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      fireEvent.click(nextButton);
      
      await waitFor(() => {
        // Deve mostrar registros 11-20
        expect(screen.getByText('Animal 11')).toBeInTheDocument();
        expect(screen.getByText('Animal 20')).toBeInTheDocument();
        expect(screen.queryByText('Animal 1')).not.toBeInTheDocument();
      });
    });

    it('deve navegar para última página', async () => {
      render(<ReproductionTable onAddReproduction={vi.fn()} onEditReproduction={vi.fn()} />);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
      });
      
      const page4Button = screen.getByRole('button', { name: /4/i });
      fireEvent.click(page4Button);
      
      await waitFor(() => {
        // Deve mostrar registros 31-35 (última página)
        expect(screen.getByText('Animal 31')).toBeInTheDocument();
        expect(screen.getByText('Animal 35')).toBeInTheDocument();
        expect(screen.queryByText('Animal 30')).not.toBeInTheDocument();
      });
    });
  });

  describe('Filtros de fase', () => {
    it('deve filtrar por fase de reprodução', async () => {
      render(<ReproductionTable onAddReproduction={vi.fn()} onEditReproduction={vi.fn()} />);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
      });
      
      // Simular filtro por fase
      const phaseFilter = screen.getByRole('button', { name: /vazias/i });
      fireEvent.click(phaseFilter);
      
      await waitFor(() => {
        // Deve mostrar apenas registros da fase selecionada
        expect(screen.getByText('1-10 de 35 registros')).toBeInTheDocument();
      });
    });
  });

  describe('Ações de edição e exclusão', () => {
    it('deve manter paginação ao editar fase', async () => {
      render(<ReproductionTable onAddReproduction={vi.fn()} onEditReproduction={vi.fn()} />);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
      });
      
      // Navegar para página 2
      const nextButton = screen.getByRole('button', { name: /next/i });
      fireEvent.click(nextButton);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 11')).toBeInTheDocument();
      });
      
      // Clicar em editar fase
      const editButton = screen.getAllByRole('button', { name: /edit/i })[0];
      fireEvent.click(editButton);
      
      // Deve manter a página atual
      expect(screen.getByText('Animal 11')).toBeInTheDocument();
    });

    it('deve manter paginação ao excluir registro', async () => {
      render(<ReproductionTable onAddReproduction={vi.fn()} onEditReproduction={vi.fn()} />);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
      });
      
      // Navegar para página 2
      const nextButton = screen.getByRole('button', { name: /next/i });
      fireEvent.click(nextButton);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 11')).toBeInTheDocument();
      });
      
      // Clicar em excluir
      const deleteButton = screen.getAllByRole('button', { name: /delete/i })[0];
      fireEvent.click(deleteButton);
      
      // Confirmar exclusão
      const confirmButton = screen.getByRole('button', { name: /sim/i });
      fireEvent.click(confirmButton);
      
      await waitFor(() => {
        // Deve manter a página atual após exclusão
        expect(screen.getByText('Animal 11')).toBeInTheDocument();
      });
    });
  });

  describe('Comportamento responsivo', () => {
    it('deve usar paginação simplificada em mobile', async () => {
      mockUseResponsive.mockReturnValue({
        isMobile: true,
        isTablet: false,
        isDesktop: false
      });

      render(<ReproductionTable onAddReproduction={vi.fn()} onEditReproduction={vi.fn()} />);
      
      await waitFor(() => {
        // Em mobile, deve mostrar apenas 5 itens por página
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
        expect(screen.getByText('Animal 5')).toBeInTheDocument();
        expect(screen.queryByText('Animal 6')).not.toBeInTheDocument();
        
        // Não deve mostrar controles de tamanho
        expect(screen.queryByRole('combobox')).not.toBeInTheDocument();
      });
    });

    it('deve usar paginação intermediária em tablet', async () => {
      mockUseResponsive.mockReturnValue({
        isMobile: false,
        isTablet: true,
        isDesktop: false
      });

      render(<ReproductionTable onAddReproduction={vi.fn()} onEditReproduction={vi.fn()} />);
      
      await waitFor(() => {
        // Em tablet, deve mostrar 8 itens por página
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
        expect(screen.getByText('Animal 8')).toBeInTheDocument();
        expect(screen.queryByText('Animal 9')).not.toBeInTheDocument();
      });
    });
  });

  describe('Estados de carregamento', () => {
    it('deve mostrar loading sem paginação', () => {
      mockUseReproduction.mockReturnValue({
        getReproductionsByFarm: vi.fn().mockResolvedValue([]),
        deleteReproduction: vi.fn().mockResolvedValue(true),
        loading: true,
        error: null
      });

      render(<ReproductionTable onAddReproduction={vi.fn()} onEditReproduction={vi.fn()} />);
      
      expect(screen.getByRole('img', { name: /loading/i })).toBeInTheDocument();
      expect(screen.queryByText(/registros/)).not.toBeInTheDocument();
    });

    it('deve mostrar erro sem paginação', () => {
      mockUseReproduction.mockReturnValue({
        getReproductionsByFarm: vi.fn().mockResolvedValue([]),
        deleteReproduction: vi.fn().mockResolvedValue(true),
        loading: false,
        error: 'Erro ao carregar dados'
      });

      render(<ReproductionTable onAddReproduction={vi.fn()} onEditReproduction={vi.fn()} />);
      
      expect(screen.getByText('Erro ao carregar dados')).toBeInTheDocument();
      expect(screen.queryByText(/registros/)).not.toBeInTheDocument();
    });
  });

  describe('Dados vazios', () => {
    it('deve lidar com lista vazia', async () => {
      mockUseReproduction.mockReturnValue({
        getReproductionsByFarm: vi.fn().mockResolvedValue([]),
        deleteReproduction: vi.fn().mockResolvedValue(true),
        loading: false,
        error: null
      });

      render(<ReproductionTable onAddReproduction={vi.fn()} onEditReproduction={vi.fn()} />);
      
      await waitFor(() => {
        expect(screen.getByText('0-0 de 0 registros')).toBeInTheDocument();
      });
    });
  });

  describe('Adição de novos registros', () => {
    it('deve resetar para primeira página ao adicionar novo registro', async () => {
      render(<ReproductionTable onAddReproduction={vi.fn()} onEditReproduction={vi.fn()} />);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
      });
      
      // Navegar para página 2
      const nextButton = screen.getByRole('button', { name: /next/i });
      fireEvent.click(nextButton);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 11')).toBeInTheDocument();
      });
      
      // Simular adição de novo registro
      const newReproduction = {
        id: 36,
        animal_id: 36,
        animal_name: 'Animal 36',
        ear_tag: 36,
        current_phase: 1,
        insemination_date: '2024-01-01',
        pregnancy_date: null,
        expected_birth_date: null,
        veterinary_confirmation: false
      };
      
      mockUseReproduction.mockReturnValue({
        getReproductionsByFarm: vi.fn().mockResolvedValue([...mockReproductions, newReproduction]),
        deleteReproduction: vi.fn().mockResolvedValue(true),
        loading: false,
        error: null
      });
      
      const { rerender } = render(<ReproductionTable onAddReproduction={vi.fn()} onEditReproduction={vi.fn()} />);
      
      await waitFor(() => {
        // Deve voltar para primeira página
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
        expect(screen.getByText('1-10 de 36 registros')).toBeInTheDocument();
      });
    });
  });
});
