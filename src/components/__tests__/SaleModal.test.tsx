import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import { SaleModal } from '../SaleModal/SaleModal';
import { Sale } from '../../types/sale';
import * as useSaleModule from '../../hooks/useSale';

jest.mock('../../hooks/useSale', () => ({
  useSaleForm: () => ({
    createSale: jest.fn(),
    updateSale: jest.fn(),
    loading: false,
  }),
}));

jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

jest.mock('dayjs', () => {
  const originalDayjs = jest.requireActual('dayjs');
  return {
    ...originalDayjs,
    default: (date?: string | number | Date) => originalDayjs(date),
  };
});

describe('SaleModal', () => {
  const mockOnCancel = jest.fn();
  const mockOnSuccess = jest.fn();

  const mockSale: Sale = {
    id: 1,
    animal_id: 1,
    farm_id: 1,
    buyer_name: 'João Silva',
    price: 1500.50,
    sale_date: '2024-01-15',
    notes: 'Test sale',
    created_at: '2024-01-15T10:00:00Z',
    updated_at: '2024-01-15T10:00:00Z',
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders create modal correctly', () => {
    render(
      <SaleModal
        visible={true}
        onCancel={mockOnCancel}
        onSuccess={mockOnSuccess}
        mode="create"
      />
    );

    expect(screen.getByText('saleModal.title.create')).toBeInTheDocument();
    expect(screen.getByText('saleModal.fields.buyerName')).toBeInTheDocument();
    expect(screen.getByText('saleModal.fields.price')).toBeInTheDocument();
    expect(screen.getByText('saleModal.fields.saleDate')).toBeInTheDocument();
    expect(screen.getByText('saleModal.fields.notes')).toBeInTheDocument();
  });

  it('renders edit modal correctly', () => {
    render(
      <SaleModal
        visible={true}
        onCancel={mockOnCancel}
        onSuccess={mockOnSuccess}
        sale={mockSale}
        mode="edit"
      />
    );

    expect(screen.getByText('saleModal.title.edit')).toBeInTheDocument();
    expect(screen.getByDisplayValue('João Silva')).toBeInTheDocument();
    expect(screen.getByDisplayValue(1500.50)).toBeInTheDocument();
  });

  it('calls onCancel when cancel button is clicked', () => {
    render(
      <SaleModal
        visible={true}
        onCancel={mockOnCancel}
        onSuccess={mockOnSuccess}
        mode="create"
      />
    );

    fireEvent.click(screen.getByText('common.cancel'));
    expect(mockOnCancel).toHaveBeenCalled();
  });

  it('shows animal name in create mode', () => {
    render(
      <SaleModal
        visible={true}
        onCancel={mockOnCancel}
        onSuccess={mockOnSuccess}
        mode="create"
      />
    );

    expect(screen.getByDisplayValue('Test Animal')).toBeInTheDocument();
  });

  it('validates required fields', async () => {
    render(
      <SaleModal
        visible={true}
        onCancel={mockOnCancel}
        onSuccess={mockOnSuccess}
        mode="create"
      />
    );

    fireEvent.click(screen.getByText('common.create'));

    await waitFor(() => {
      expect(screen.getByText('saleModal.validation.buyerNameRequired')).toBeInTheDocument();
    });
  });

  it('submits form with valid data', async () => {
    const mockCreateSale = jest.fn().mockResolvedValue(undefined);
    vi.mocked(useSaleModule.useSaleForm).mockReturnValue({
      createSale: mockCreateSale,
      updateSale: jest.fn(),
      loading: false,
    });

    render(
      <SaleModal
        visible={true}
        onCancel={mockOnCancel}
        onSuccess={mockOnSuccess}
        mode="create"
      />
    );

    fireEvent.change(screen.getByPlaceholderText('saleModal.placeholders.buyerName'), {
      target: { value: 'João Silva' },
    });
    fireEvent.change(screen.getByPlaceholderText('saleModal.placeholders.price'), {
      target: { value: '1500.50' },
    });

    fireEvent.click(screen.getByText('common.create'));

    await waitFor(() => {
      expect(mockCreateSale).toHaveBeenCalledWith({
        animal_id: 1,
        buyer_name: 'João Silva',
        price: 1500.50,
        sale_date: expect.any(String),
        notes: '',
      });
    });
  });
});
