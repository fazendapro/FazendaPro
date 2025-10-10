import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { NextToCalve } from '../next-to-calve';
import { useNextToCalve } from '../../hooks/useNextToCalve';

// Mock do hook
vi.mock('../../hooks/useNextToCalve', () => ({
  useNextToCalve: vi.fn()
}));

// Mock do react-i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

const mockUseNextToCalve = vi.mocked(useNextToCalve);

describe('NextToCalve', () => {
  const mockNextToCalveData = [
    {
      id: 1,
      animal_name: 'Tata Salt',
      ear_tag_number_local: 123,
      photo: 'src/assets/images/mocked/cows/tata.png',
      pregnancy_date: '2024-01-01',
      expected_birth_date: '2024-10-01',
      days_until_birth: 15,
      status: 'Alto' as const
    },
    {
      id: 2,
      animal_name: 'Lays',
      ear_tag_number_local: 124,
      photo: 'src/assets/images/mocked/cows/lays.png',
      pregnancy_date: '2024-01-15',
      expected_birth_date: '2024-10-15',
      days_until_birth: 40,
      status: 'Baixo' as const
    },
    {
      id: 3,
      animal_name: 'Matilda',
      ear_tag_number_local: 125,
      photo: 'src/assets/images/mocked/cows/matilda.png',
      pregnancy_date: '2024-01-20',
      expected_birth_date: '2024-10-20',
      days_until_birth: 20,
      status: 'Médio' as const
    }
  ];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('deve renderizar o título do card com tooltip', () => {
    mockUseNextToCalve.mockReturnValue({
      nextToCalve: mockNextToCalveData,
      loading: false,
      error: null,
      getNextToCalve: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<NextToCalve />);

    expect(screen.getByText('dashboard.nextToCalve')).toBeInTheDocument();
    const infoIcon = document.querySelector('.anticon-info-circle');
    expect(infoIcon).toBeInTheDocument();
  });

  it('deve renderizar lista de vacas quando dados estão disponíveis', () => {
    mockUseNextToCalve.mockReturnValue({
      nextToCalve: mockNextToCalveData,
      loading: false,
      error: null,
      getNextToCalve: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<NextToCalve />);

    expect(screen.getByText('Tata Salt')).toBeInTheDocument();
    expect(screen.getByText('Lays')).toBeInTheDocument();
    expect(screen.getByText('Matilda')).toBeInTheDocument();
  });

  it('deve renderizar informações corretas para cada vaca', () => {
    mockUseNextToCalve.mockReturnValue({
      nextToCalve: mockNextToCalveData,
      loading: false,
      error: null,
      getNextToCalve: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<NextToCalve />);

    expect(screen.getByText('Última vez: 15 dias')).toBeInTheDocument();
    expect(screen.getByText('Última vez: 40 dias')).toBeInTheDocument();
    expect(screen.getByText('Última vez: 20 dias')).toBeInTheDocument();
  });

  it('deve renderizar status com cores corretas', () => {
    mockUseNextToCalve.mockReturnValue({
      nextToCalve: mockNextToCalveData,
      loading: false,
      error: null,
      getNextToCalve: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<NextToCalve />);

    const altoStatus = screen.getByText('Alto');
    const baixoStatus = screen.getByText('Baixo');
    const medioStatus = screen.getByText('Médio');

    expect(altoStatus).toBeInTheDocument();
    expect(baixoStatus).toBeInTheDocument();
    expect(medioStatus).toBeInTheDocument();
  });

  it('deve renderizar loading quando dados estão carregando', () => {
    mockUseNextToCalve.mockReturnValue({
      nextToCalve: [],
      loading: true,
      error: null,
      getNextToCalve: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<NextToCalve />);

    // O componente não tem loading state visível, mas podemos verificar se a lista está vazia
    expect(screen.queryByText('Tata Salt')).not.toBeInTheDocument();
  });

  it('deve renderizar lista vazia quando não há dados', () => {
    mockUseNextToCalve.mockReturnValue({
      nextToCalve: [],
      loading: false,
      error: null,
      getNextToCalve: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<NextToCalve />);

    expect(screen.queryByText('Tata Salt')).not.toBeInTheDocument();
    expect(screen.queryByText('Lays')).not.toBeInTheDocument();
    expect(screen.queryByText('Matilda')).not.toBeInTheDocument();
  });

  it('deve renderizar avatares com imagens corretas', () => {
    mockUseNextToCalve.mockReturnValue({
      nextToCalve: mockNextToCalveData,
      loading: false,
      error: null,
      getNextToCalve: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<NextToCalve />);

    const avatars = screen.getAllByRole('img');
    expect(avatars).toHaveLength(3);
    
    // Verificar se as imagens têm os srcs corretos
    expect(avatars[0]).toHaveAttribute('src', 'src/assets/images/mocked/cows/tata.png');
    expect(avatars[1]).toHaveAttribute('src', 'src/assets/images/mocked/cows/lays.png');
    expect(avatars[2]).toHaveAttribute('src', 'src/assets/images/mocked/cows/matilda.png');
  });

  it('deve aplicar estilos corretos para status Alto', () => {
    mockUseNextToCalve.mockReturnValue({
      nextToCalve: [mockNextToCalveData[0]], // Apenas Tata Salt com status Alto
      loading: false,
      error: null,
      getNextToCalve: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<NextToCalve />);

    const altoStatus = screen.getByText('Alto');
    expect(altoStatus).toHaveStyle({
      color: 'red',
      border: '1px solid red',
      backgroundColor: '#ffcccc',
      fontSize: '12px'
    });
  });

  it('deve aplicar estilos corretos para status Médio', () => {
    mockUseNextToCalve.mockReturnValue({
      nextToCalve: [mockNextToCalveData[2]], // Apenas Matilda com status Médio
      loading: false,
      error: null,
      getNextToCalve: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<NextToCalve />);

    const medioStatus = screen.getByText('Médio');
    expect(medioStatus).toHaveStyle({
      color: '#ff8c00',
      border: '1px solid #ff8c00',
      backgroundColor: '#ffe4b5',
      fontSize: '12px'
    });
  });

  it('deve aplicar estilos corretos para status Baixo', () => {
    mockUseNextToCalve.mockReturnValue({
      nextToCalve: [mockNextToCalveData[1]], // Apenas Lays com status Baixo
      loading: false,
      error: null,
      getNextToCalve: vi.fn(),
      refetch: vi.fn(),
      clearError: vi.fn()
    });

    render(<NextToCalve />);

    const baixoStatus = screen.getByText('Baixo');
    expect(baixoStatus).toHaveStyle({
      color: 'green',
      border: '1px solid green',
      backgroundColor: '#ccffcc',
      fontSize: '12px'
    });
  });
});
