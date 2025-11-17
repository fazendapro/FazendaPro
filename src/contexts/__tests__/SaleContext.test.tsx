import { render, screen, act, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import { SaleProvider, useSaleContext } from '../SaleContext';
import * as saleServiceModule from '../../components/services/saleService';

jest.mock('../../components/services/saleService', () => ({
  saleService: {
    createSale: jest.fn(),
    getSalesByFarm: jest.fn(),
    getSalesHistory: jest.fn(),
    getSalesByAnimal: jest.fn(),
    getSalesByDateRange: jest.fn(),
    getSaleById: jest.fn(),
    updateSale: jest.fn(),
    deleteSale: jest.fn(),
  },
}));

jest.mock('react-toastify', () => ({
  toast: {
    success: jest.fn(),
    error: jest.fn(),
  },
}));

jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

const TestComponent = () => {
  const { sales, loading, error, createSale, updateSale, deleteSale } = useSaleContext();

  return (
    <div>
      <div data-testid="sales-count">{sales.length}</div>
      <div data-testid="loading">{loading.toString()}</div>
      <div data-testid="error">{error}</div>
      <button
        data-testid="create-sale"
        onClick={() => createSale({
          animal_id: 1,
          buyer_name: 'Test Buyer',
          price: 1000,
          sale_date: '2024-01-01',
          notes: 'Test notes',
        })}
      >
        Create Sale
      </button>
      <button
        data-testid="update-sale"
        onClick={() => updateSale(1, {
          buyer_name: 'Updated Buyer',
          price: 1500,
          sale_date: '2024-01-02',
          notes: 'Updated notes',
        })}
      >
        Update Sale
      </button>
      <button
        data-testid="delete-sale"
        onClick={() => deleteSale(1)}
      >
        Delete Sale
      </button>
    </div>
  );
};

describe('SaleContext', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('provides initial state', () => {
    render(
      <SaleProvider>
        <TestComponent />
      </SaleProvider>
    );

    expect(screen.getByTestId('sales-count')).toHaveTextContent('0');
    expect(screen.getByTestId('loading')).toHaveTextContent('false');
    expect(screen.getByTestId('error')).toHaveTextContent('');
  });

  it('handles createSale successfully', async () => {
    const mockSale = {
      id: 1,
      animal_id: 1,
      farm_id: 1,
      buyer_name: 'Test Buyer',
      price: 1000,
      sale_date: '2024-01-01',
      notes: 'Test notes',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    };
    vi.mocked(saleServiceModule.saleService.createSale).mockResolvedValue(mockSale);

    render(
      <SaleProvider>
        <TestComponent />
      </SaleProvider>
    );

    await act(async () => {
      fireEvent.click(screen.getByTestId('create-sale'));
    });

    expect(saleServiceModule.saleService.createSale).toHaveBeenCalledWith({
      animal_id: 1,
      buyer_name: 'Test Buyer',
      price: 1000,
      sale_date: '2024-01-01',
      notes: 'Test notes',
    });
  });

  it('handles createSale error', async () => {
    vi.mocked(saleServiceModule.saleService.createSale).mockRejectedValue(new Error('Network error'));

    render(
      <SaleProvider>
        <TestComponent />
      </SaleProvider>
    );

    await act(async () => {
      fireEvent.click(screen.getByTestId('create-sale'));
    });

    expect(saleServiceModule.saleService.createSale).toHaveBeenCalled();
  });

  it('handles updateSale successfully', async () => {
    const mockUpdatedSale = {
      id: 1,
      animal_id: 1,
      farm_id: 1,
      buyer_name: 'Updated Buyer',
      price: 1500,
      sale_date: '2024-01-02',
      notes: 'Updated notes',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-02T00:00:00Z',
    };
    vi.mocked(saleServiceModule.saleService.updateSale).mockResolvedValue(mockUpdatedSale);

    render(
      <SaleProvider>
        <TestComponent />
      </SaleProvider>
    );

    await act(async () => {
      fireEvent.click(screen.getByTestId('update-sale'));
    });

    expect(saleServiceModule.saleService.updateSale).toHaveBeenCalledWith(1, {
      buyer_name: 'Updated Buyer',
      price: 1500,
      sale_date: '2024-01-02',
      notes: 'Updated notes',
    });
  });

  it('handles deleteSale successfully', async () => {
    vi.mocked(saleServiceModule.saleService.deleteSale).mockResolvedValue(undefined);

    render(
      <SaleProvider>
        <TestComponent />
      </SaleProvider>
    );

    await act(async () => {
      fireEvent.click(screen.getByTestId('delete-sale'));
    });

    expect(saleService.deleteSale).toHaveBeenCalledWith(1);
  });

  it('throws error when used outside provider', () => {
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
    
    expect(() => {
      render(<TestComponent />);
    }).toThrow('useSaleContext must be used within a SaleProvider');
    
    consoleSpy.mockRestore();
  });
});
