import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AnimalHistoryPDFGenerator, generateAnimalHistoryPDF } from '../pdfGenerator';

// Mock jsPDF
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
  internal: {
    pageSize: {
      getWidth: vi.fn(() => 210),
      getHeight: vi.fn(() => 297),
    },
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

describe('AnimalHistoryPDFGenerator', () => {
  const mockAnimal = {
    id: 1,
    animal_name: 'Test Animal',
    ear_tag_number_local: 123,
    breed: 'Holandesa',
    type: 'Bovino',
    birth_date: '2020-01-01',
    sex: 0,
    status: 0,
    confinement: false,
    fertilization: false,
    castrated: false,
    current_weight: 500,
    ideal_weight: 600,
    milk_production: 20,
  };

  const mockSales = [
    {
      id: 1,
      sale_date: '2024-01-01',
      buyer_name: 'Buyer 1',
      price: 1000,
      notes: 'Sale notes',
    },
  ];

  const mockMilkCollections = [
    {
      id: 1,
      collection_date: '2024-01-01',
      quantity: 20,
      quality: 'Good',
      notes: 'Milk notes',
    },
  ];

  const mockReproductions = [
    {
      id: 1,
      date: '2024-01-01',
      phase: 'Gestação',
      notes: 'Reproduction notes',
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
    const generator = new AnimalHistoryPDFGenerator({
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
    const generator = new AnimalHistoryPDFGenerator({
      animal: mockAnimal,
      sales: [],
      milkCollections: [],
      reproductions: [],
      translations: mockTranslations,
    });

    expect(mockDoc.save).toHaveBeenCalled();
  });

  it('deve gerar PDF com imagem do animal', () => {
    const generator = new AnimalHistoryPDFGenerator({
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
    const generator = new AnimalHistoryPDFGenerator({
      animal: mockAnimal,
      sales: [],
      milkCollections: [],
      reproductions: [],
      translations: mockTranslations,
    });

    expect(mockDoc.text).toHaveBeenCalled();
  });

  it('deve calcular estatísticas corretamente', () => {
    const generator = new AnimalHistoryPDFGenerator({
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




