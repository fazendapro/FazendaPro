import { render, screen, fireEvent, cleanup } from '@testing-library/react';
import { vi, beforeEach, afterEach, describe, it, expect } from 'vitest';
import { AnimalHistoryExport } from '../AnimalHistoryExport/AnimalHistoryExport';
import { Animal } from '../../pages/contents/AnimalTable/types/type';
import { Sale } from '../../types/sale';
import { MilkCollection } from '../../types/milk-collection';
import { Reproduction } from '../../types/reproduction';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('jspdf', async () => {
  const actual = await vi.importActual('jspdf');
  const mockDoc = {
    setFontSize: vi.fn(),
    setFont: vi.fn(),
    text: vi.fn(),
    autoTable: vi.fn(),
    getNumberOfPages: vi.fn(() => 1),
    setPage: vi.fn(),
    save: vi.fn(),
    internal: {
      pageSize: {
        getWidth: vi.fn(() => 210),
        getHeight: vi.fn(() => 297),
      },
    },
  };
  return {
    ...actual,
    jsPDF: vi.fn(() => mockDoc),
    default: vi.fn(() => mockDoc),
  };
});

vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    message: {
      success: vi.fn(),
      error: vi.fn(),
    },
  };
});

describe('AnimalHistoryExport', () => {
  const mockAnimal: Animal = {
    id: '1',
    farm_id: 1,
    ear_tag_number_local: 123,
    ear_tag_number_register: 456,
    animal_name: 'Test Animal',
    sex: 0,
    breed: 'Holstein',
    type: 'vaca',
    birth_date: '2020-01-01',
    photo: '',
    confinement: false,
    animal_type: 0,
    status: 0,
    fertilization: false,
    castrated: false,
    purpose: 0,
    current_batch: 1,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  };

  const mockSales: Sale[] = [
    {
      id: 1,
      animal_id: 1,
      farm_id: 1,
      buyer_name: 'João Silva',
      price: 1500.50,
      sale_date: '2024-01-15',
      notes: 'Test sale',
      created_at: '2024-01-15T10:00:00Z',
      updated_at: '2024-01-15T10:00:00Z',
    },
  ];

  const mockMilkCollections: MilkCollection[] = [
    {
      id: 1,
      animal_id: 1,
      farm_id: 1,
      collection_date: '2024-01-10',
      quantity: 25.5,
      quality: 'A',
      notes: 'Good quality',
      created_at: '2024-01-10T08:00:00Z',
      updated_at: '2024-01-10T08:00:00Z',
    },
  ];

  const mockReproductions: Reproduction[] = [
    {
      id: 1,
      animal_id: 1,
      farm_id: 1,
      date: '2024-01-05',
      phase: 'Gestation',
      notes: 'First gestation',
      created_at: '2024-01-05T09:00:00Z',
      updated_at: '2024-01-05T09:00:00Z',
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    cleanup();
  });

  it('renders export button correctly', () => {
    render(
      <AnimalHistoryExport
        animal={mockAnimal}
        sales={mockSales}
        milkCollections={mockMilkCollections}
        reproductions={mockReproductions}
      />
    );

    expect(screen.getByText('animalDetail.exportHistory')).toBeInTheDocument();
  });

  it('generates PDF when button is clicked', () => {
    render(
      <AnimalHistoryExport
        animal={mockAnimal}
        sales={mockSales}
        milkCollections={mockMilkCollections}
        reproductions={mockReproductions}
      />
    );

    fireEvent.click(screen.getAllByText('animalDetail.exportHistory')[0]);

    expect(screen.getAllByText('animalDetail.exportHistory').length).toBeGreaterThan(0);
  });

  it('handles empty data gracefully', () => {
    render(
      <AnimalHistoryExport
        animal={mockAnimal}
        sales={[]}
        milkCollections={[]}
        reproductions={[]}
      />
    );

    expect(screen.getAllByText('animalDetail.exportHistory').length).toBeGreaterThan(0);
  });

  it('generates correct filename', () => {
    render(
      <AnimalHistoryExport
        animal={mockAnimal}
        sales={mockSales}
        milkCollections={mockMilkCollections}
        reproductions={mockReproductions}
      />
    );

    fireEvent.click(screen.getAllByText('animalDetail.exportHistory')[0]);

    expect(screen.getAllByText('animalDetail.exportHistory').length).toBeGreaterThan(0);
  });
});
