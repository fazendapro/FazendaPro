import { render, screen, fireEvent } from '@testing-library/react';
import { AnimalHistoryExport } from '../AnimalHistoryExport/AnimalHistoryExport';
import { Animal } from '../../pages/contents/AnimalTable/types/type';
import { Sale } from '../../types/sale';
import { MilkCollection } from '../../types/milk-collection';
import { Reproduction } from '../../types/reproduction';

jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

jest.mock('jspdf', () => {
  const mockDoc = {
    setFontSize: jest.fn(),
    setFont: jest.fn(),
    text: jest.fn(),
    autoTable: jest.fn(),
    getNumberOfPages: jest.fn(() => 1),
    setPage: jest.fn(),
    save: jest.fn(),
    internal: {
      pageSize: {
        getWidth: jest.fn(() => 210),
        getHeight: jest.fn(() => 297),
      },
    },
  };
  return jest.fn(() => mockDoc);
});

jest.mock('antd', () => ({
  ...jest.requireActual('antd'),
  message: {
    success: jest.fn(),
    error: jest.fn(),
  },
}));

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
    jest.clearAllMocks();
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
    const { jsPDF } = require('jspdf');
    const mockDoc = new jsPDF();

    render(
      <AnimalHistoryExport
        animal={mockAnimal}
        sales={mockSales}
        milkCollections={mockMilkCollections}
        reproductions={mockReproductions}
      />
    );

    fireEvent.click(screen.getByText('animalDetail.exportHistory'));

    expect(mockDoc.setFontSize).toHaveBeenCalled();
    expect(mockDoc.text).toHaveBeenCalled();
    expect(mockDoc.save).toHaveBeenCalled();
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

    expect(screen.getByText('animalDetail.exportHistory')).toBeInTheDocument();
  });

  it('generates correct filename', () => {
    const { jsPDF } = require('jspdf');
    const mockDoc = new jsPDF();

    render(
      <AnimalHistoryExport
        animal={mockAnimal}
        sales={mockSales}
        milkCollections={mockMilkCollections}
        reproductions={mockReproductions}
      />
    );

    fireEvent.click(screen.getByText('animalDetail.exportHistory'));

    expect(mockDoc.save).toHaveBeenCalledWith(
      expect.stringMatching(/Test_Animal_historico_\d{4}-\d{2}-\d{2}\.pdf/)
    );
  });
});
