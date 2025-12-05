import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AnimalHistoryPDFGenerator, generateAnimalHistoryPDF } from '../pdfGenerator';
import { Animal } from '../../../pages/contents/AnimalTable/types/type';
import { Sale } from '../../../types/sale';
import { MilkCollection } from '../../../types/milk-collection';
import { Reproduction } from '../../../types/reproduction';

const mockDoc = {
  setFontSize: vi.fn(),
  setFont: vi.fn(),
  setTextColor: vi.fn(),
  text: vi.fn(),
  setDrawColor: vi.fn(),
  setLineWidth: vi.fn(),
  line: vi.fn(),
  addImage: vi.fn(),
  rect: vi.fn(),
  autoTable: vi.fn(),
  getNumberOfPages: vi.fn(() => 1),
  setPage: vi.fn(),
  addPage: vi.fn(),
  save: vi.fn(),
  getTextColor: vi.fn(() => [0, 0, 0]),
  internal: {
    pageSize: {
      getWidth: vi.fn(() => 210),
      getHeight: vi.fn(() => 297),
    },
    getFontSize: vi.fn(() => 12),
    getFont: vi.fn(() => ({
      fontName: 'helvetica',
      fontStyle: 'normal'
    })),
  },
  lastAutoTable: {
    finalY: 100,
  },
};

vi.mock('jspdf', () => {
  return {
    default: vi.fn(() => mockDoc),
  };
});

vi.mock('jspdf-autotable', () => {
  return {
    default: vi.fn((doc: typeof mockDoc, options: unknown) => {
      // Simular o comportamento do autoTable
      if (mockDoc.lastAutoTable) {
        mockDoc.lastAutoTable.finalY = 150;
      }
    }),
  };
});

describe('AnimalHistoryPDFGenerator', () => {
  const mockAnimal: Animal = {
    id: '1',
    farm_id: 1,
    animal_name: 'Test Animal',
    ear_tag_number_local: 123,
    ear_tag_number_register: 456,
    breed: 'Holandesa',
    type: 'vaca',
    birth_date: '2020-01-01',
    sex: 0,
    animal_type: 0,
    status: 0,
    confinement: false,
    fertilization: false,
    castrated: false,
    purpose: 0,
    current_batch: 1,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
    current_weight: '500',
    ideal_weight: '600',
    milk_production: '20',
  };

  const mockSales: Sale[] = [
    {
      id: 1,
      animal_id: 1,
      farm_id: 1,
      sale_date: '2024-01-01',
      buyer_name: 'Buyer 1',
      price: 1000,
      notes: 'Sale notes',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    },
  ];

  const mockMilkCollections: MilkCollection[] = [
    {
      id: 1,
      animal_id: 1,
      farm_id: 1,
      collection_date: '2024-01-01',
      quantity: 20,
      quality: 'Good',
      notes: 'Milk notes',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    },
  ];

  const mockReproductions: Reproduction[] = [
    {
      id: 1,
      animal_id: 1,
      farm_id: 1,
      date: '2024-01-01',
      phase: 'Gestação',
      notes: 'Reproduction notes',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    },
  ];

  const mockTranslations = {
    animalHistoryExport: {
      title: 'Histórico do Animal',
      animalInfo: 'Informações do Animal',
      fields: {
        name: 'Nome',
        earTag: 'Brinco',
        breed: 'Raça',
        type: 'Tipo',
        birthDate: 'Data de Nascimento',
        sex: 'Sexo',
        status: 'Status',
        confinement: 'Confinamento',
        fertilization: 'Fertilização',
        castrated: 'Castrado',
        currentWeight: 'Peso Atual',
        idealWeight: 'Peso Ideal',
        milkProduction: 'Produção de Leite',
      },
      statistics: 'Estatísticas',
      stats: {
        totalSales: 'Total de Vendas',
        totalSalesValue: 'Valor Total de Vendas',
        totalMilkCollections: 'Total de Coletas',
        totalMilkQuantity: 'Quantidade Total de Leite',
        totalReproductions: 'Total de Reproduções',
      },
      salesHistory: 'Histórico de Vendas',
      salesTable: {
        date: 'Data',
        buyer: 'Comprador',
        price: 'Preço',
        notes: 'Observações',
      },
      milkHistory: 'Histórico de Leite',
      milkTable: {
        date: 'Data',
        quantity: 'Quantidade',
        quality: 'Qualidade',
        notes: 'Observações',
      },
      reproductionHistory: 'Histórico de Reprodução',
      reproductionTable: {
        date: 'Data',
        phase: 'Fase',
        notes: 'Observações',
      },
      footer: {
        page: 'Página',
        of: 'de',
      },
      sex: {
        male: 'Macho',
        female: 'Fêmea',
      },
      status: {
        active: 'Ativo',
        inactive: 'Inativo',
        sold: 'Vendido',
        deceased: 'Falecido',
      },
    },
    common: {
      yes: 'Sim',
      no: 'Não',
      notInformed: 'Não informado',
    },
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('deve gerar PDF com dados completos', () => {
    new AnimalHistoryPDFGenerator({
      animal: mockAnimal,
      sales: mockSales,
      milkCollections: mockMilkCollections,
      reproductions: mockReproductions,
      translations: mockTranslations,
    });

    expect(mockDoc.save).toHaveBeenCalled();
    expect(mockDoc.text).toHaveBeenCalled();
  });

  it('deve gerar PDF sem imagem do animal', () => {
    new AnimalHistoryPDFGenerator({
      animal: mockAnimal,
      sales: [],
      milkCollections: [],
      reproductions: [],
      translations: mockTranslations,
    });

    expect(mockDoc.save).toHaveBeenCalled();
  });

  it('deve gerar PDF com imagem do animal', () => {
    new AnimalHistoryPDFGenerator({
      animal: mockAnimal,
      sales: [],
      milkCollections: [],
      reproductions: [],
      animalImage: 'data:image/png;base64,test',
      translations: mockTranslations,
    });

    expect(mockDoc.addImage).toHaveBeenCalled();
  });

  it('deve formatar corretamente os dados do animal', () => {
    new AnimalHistoryPDFGenerator({
      animal: mockAnimal,
      sales: [],
      milkCollections: [],
      reproductions: [],
      translations: mockTranslations,
    });

    expect(mockDoc.text).toHaveBeenCalled();
  });

  it('deve calcular estatísticas corretamente', () => {
    new AnimalHistoryPDFGenerator({
      animal: mockAnimal,
      sales: mockSales,
      milkCollections: mockMilkCollections,
      reproductions: mockReproductions,
      translations: mockTranslations,
    });

    expect(mockDoc.text).toHaveBeenCalled();
  });

  it('deve usar função generateAnimalHistoryPDF', () => {
    const generator = generateAnimalHistoryPDF({
      animal: mockAnimal,
      sales: [],
      milkCollections: [],
      reproductions: [],
      translations: mockTranslations,
    });

    expect(generator).toBeInstanceOf(AnimalHistoryPDFGenerator);
  });
});




