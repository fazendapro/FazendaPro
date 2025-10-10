import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import { AnimalTable } from '../animal-table';
import { useAnimals } from '../../../hooks/useAnimals';
import { useFarm } from '../../../../../../hooks/useFarm';
import { useAnimalColumnBuilder } from '../column-builder.tsx';
import { useResponsive } from '../../../../../../hooks';

vi.mock('../../../hooks/useAnimals');
vi.mock('../../../../../../hooks/useFarm');
vi.mock('../column-builder.tsx');
vi.mock('../../../../../../hooks');

const mockUseAnimals = useAnimals as any;
const mockUseFarm = useFarm as any;
const mockUseAnimalColumnBuilder = useAnimalColumnBuilder as any;
const mockUseResponsive = useResponsive as any;

describe('AnimalTable Pagination Integration', () => {
  const mockAnimals = Array.from({ length: 25 }, (_, i) => ({
    id: String(i + 1),
    animal_name: `Animal ${i + 1}`,
    ear_tag_number_local: i + 1,
    ear_tag_number_register: i + 1,
    breed: 'Holandês',
    type: 'Vaca',
    sex: 'Fêmea',
    birth_date: '2020-01-01',
    farm_id: 1,
    animal_type: 1,
    status: 1,
    confinement: true,
    weight: 500,
    color: 'Preto e Branco',
    mother_id: null,
    father_id: null,
    created_at: '2024-01-01',
    updated_at: '2024-01-01',
    createdAt: '2024-01-01',
    updatedAt: '2024-01-01',
    fertilization: false,
    castrated: false,
    purpose: 1,
    current_batch: null,
    batch_id: null,
    batch_name: null
  }));

  const mockColumns = [
    {
      title: 'Nome',
      dataIndex: 'animal_name',
      key: 'animal_name'
    },
    {
      title: 'Brinco',
      dataIndex: 'ear_tag_number_local',
      key: 'ear_tag_number_local'
    }
  ];

  beforeEach(() => {
    vi.clearAllMocks();
    
    mockUseFarm.mockReturnValue({
      farm: { 
        id: 1, 
        name: 'Test Farm',
        location: 'Test Location',
        created_at: '2024-01-01',
        updated_at: '2024-01-01'
      }
    });
    
    mockUseAnimals.mockReturnValue({
      animals: mockAnimals,
      loading: false,
      error: null,
      refetch: vi.fn(),
      getAnimalsByFarm: vi.fn(),
      clearError: vi.fn()
    });
    
    mockUseAnimalColumnBuilder.mockReturnValue({
      getAllColumns: vi.fn().mockReturnValue([]),
      getDefaultColumns: vi.fn().mockReturnValue([]),
      getDefaultColumnKeys: vi.fn().mockReturnValue([]),
      getColumnsByKeys: vi.fn().mockReturnValue([]),
      getColumnOptions: vi.fn().mockReturnValue([]),
      buildTableColumns: vi.fn().mockReturnValue(mockColumns)
    });
    
    mockUseResponsive.mockReturnValue({
      isMobile: false,
      isTablet: false,
      isDesktop: true,
      isLargeDesktop: false,
      screenWidth: 1200
    });
  });

  describe('Paginação básica', () => {
    it('deve renderizar primeira página com 10 itens por padrão', () => {
      render(<AnimalTable />);
      
      expect(screen.getByText('Animal 1')).toBeInTheDocument();
      expect(screen.getByText('Animal 10')).toBeInTheDocument();
      expect(screen.queryByText('Animal 11')).not.toBeInTheDocument();
      
      expect(screen.getByText('1-10 de 25 registros')).toBeInTheDocument();
    });

    it('deve navegar para próxima página', async () => {
      render(<AnimalTable />);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      fireEvent.click(nextButton);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 11')).toBeInTheDocument();
        expect(screen.getByText('Animal 20')).toBeInTheDocument();
        expect(screen.queryByText('Animal 1')).not.toBeInTheDocument();
      });
    });

    it('deve navegar para página específica', async () => {
      render(<AnimalTable />);
      
      const page3Button = screen.getByRole('button', { name: /3/i });
      fireEvent.click(page3Button);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 21')).toBeInTheDocument();
        expect(screen.getByText('Animal 25')).toBeInTheDocument();
        expect(screen.queryByText('Animal 20')).not.toBeInTheDocument();
      });
    });
  });

  describe('Controle de tamanho da página', () => {
    it('deve alterar tamanho da página', async () => {
      render(<AnimalTable />);
      
      const sizeChanger = screen.getByRole('combobox');
      fireEvent.click(sizeChanger);
      
      const option20 = screen.getByText('20 / page');
      fireEvent.click(option20);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
        expect(screen.getByText('Animal 20')).toBeInTheDocument();
        expect(screen.getByText('1-20 de 25 registros')).toBeInTheDocument();
      });
    });

    it('deve resetar para primeira página ao alterar tamanho', async () => {
      render(<AnimalTable />);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      fireEvent.click(nextButton);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 11')).toBeInTheDocument();
      });
      
      const sizeChanger = screen.getByRole('combobox');
      fireEvent.click(sizeChanger);
      
      const option5 = screen.getByText('5 / page');
      fireEvent.click(option5);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
        expect(screen.getByText('1-5 de 25 registros')).toBeInTheDocument();
      });
    });
  });

  describe('Comportamento responsivo', () => {
    it('deve usar paginação simplificada em mobile', () => {
      mockUseResponsive.mockReturnValue({
      isMobile: true,
      isTablet: false,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 400
      });

      render(<AnimalTable />);
      
      expect(screen.getByText('Animal 1')).toBeInTheDocument();
      expect(screen.getByText('Animal 5')).toBeInTheDocument();
      expect(screen.queryByText('Animal 6')).not.toBeInTheDocument();
      
      expect(screen.queryByRole('combobox')).not.toBeInTheDocument();
    });

    it('deve usar paginação intermediária em tablet', () => {
      mockUseResponsive.mockReturnValue({
      isMobile: false,
      isTablet: true,
      isDesktop: false,
      isLargeDesktop: false,
      screenWidth: 800
      });

      render(<AnimalTable />);
      
      expect(screen.getByText('Animal 1')).toBeInTheDocument();
      expect(screen.getByText('Animal 8')).toBeInTheDocument();
      expect(screen.queryByText('Animal 9')).not.toBeInTheDocument();
    });
  });

  describe('Filtros e busca', () => {
    it('deve manter paginação com filtros aplicados', async () => {
      render(<AnimalTable searchTerm="Animal 1" />);
      
      expect(screen.getByText('Animal 1')).toBeInTheDocument();
      expect(screen.getByText('Animal 10')).toBeInTheDocument();
      expect(screen.getByText('Animal 11')).toBeInTheDocument();
      expect(screen.getByText('Animal 12')).toBeInTheDocument();
      expect(screen.getByText('Animal 13')).toBeInTheDocument();
      expect(screen.getByText('Animal 14')).toBeInTheDocument();
      expect(screen.getByText('Animal 15')).toBeInTheDocument();
      expect(screen.getByText('Animal 16')).toBeInTheDocument();
      expect(screen.getByText('Animal 17')).toBeInTheDocument();
      expect(screen.getByText('Animal 18')).toBeInTheDocument();
      expect(screen.getByText('Animal 19')).toBeInTheDocument();
      
      expect(screen.queryByText('Animal 2')).not.toBeInTheDocument();
      expect(screen.queryByText('Animal 3')).not.toBeInTheDocument();
    });

    it('deve resetar para primeira página ao aplicar filtro', async () => {
      render(<AnimalTable />);
      
      const nextButton = screen.getByRole('button', { name: /next/i });
      fireEvent.click(nextButton);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 11')).toBeInTheDocument();
      });
      
      render(<AnimalTable searchTerm="Animal 1" />);
      
      await waitFor(() => {
        expect(screen.getByText('Animal 1')).toBeInTheDocument();
      });
    });
  });

  describe('Estados de carregamento e erro', () => {
    it('deve mostrar loading sem paginação', () => {
      mockUseAnimals.mockReturnValue({
        animals: [],
        loading: true,
        error: null,
        refetch: vi.fn(),
        getAnimalsByFarm: vi.fn(),
        clearError: vi.fn()
      });

      render(<AnimalTable />);
      
      expect(screen.getByRole('img', { name: /loading/i })).toBeInTheDocument();
      expect(screen.queryByText(/registros/)).not.toBeInTheDocument();
    });

    it('deve mostrar erro sem paginação', () => {
      mockUseAnimals.mockReturnValue({
        animals: [],
        loading: false,
        error: 'Erro ao carregar dados',
        refetch: vi.fn(),
        getAnimalsByFarm: vi.fn(),
        clearError: vi.fn()
      });

      render(<AnimalTable />);
      
      expect(screen.getByText('Erro ao carregar dados')).toBeInTheDocument();
      expect(screen.queryByText(/registros/)).not.toBeInTheDocument();
    });
  });

  describe('Dados vazios', () => {
    it('deve lidar com lista vazia', () => {
      mockUseAnimals.mockReturnValue({
        animals: [],
        loading: false,
        error: null,
        refetch: vi.fn(),
        getAnimalsByFarm: vi.fn(),
        clearError: vi.fn()
      });

      render(<AnimalTable />);
      
      expect(screen.getByText('0-0 de 0 registros')).toBeInTheDocument();
    });
  });
});
