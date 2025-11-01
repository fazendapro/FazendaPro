import { describe, it, expect, vi, beforeEach } from 'vitest';
import { RemoteGetMonthlySalesAndPurchases } from '../remote-get-monthly-sales-and-purchases';
import { api } from '../../../../../../components';
import { AxiosError } from 'axios';

vi.mock('../../../../../../components', () => ({
  api: vi.fn()
}));

vi.mock('i18next', () => ({
  t: vi.fn((key: string) => key)
}));

const mockApi = vi.mocked(api);

describe('RemoteGetMonthlySalesAndPurchases', () => {
  let remoteGetMonthlySalesAndPurchases: RemoteGetMonthlySalesAndPurchases;

  beforeEach(() => {
    vi.clearAllMocks();
    remoteGetMonthlySalesAndPurchases = new RemoteGetMonthlySalesAndPurchases();
  });

  it('deve buscar dados mensais de vendas e compras com sucesso', async () => {
    const mockResponse = {
      data: {
        message: 'Success',
        data: {
          sales: [
            { month: 'Jan', year: 2024, sales: 15000, count: 5 },
            { month: 'Fev', year: 2024, sales: 20000, count: 7 },
          ],
          purchases: [
            { month: 'Jan', year: 2024, total: 0 },
            { month: 'Fev', year: 2024, total: 0 },
          ],
        },
        success: true
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    };

    mockApi.mockReturnValue(mockApiInstance as any);

    const result = await remoteGetMonthlySalesAndPurchases.getMonthlySalesAndPurchases({ 
      farm_id: 1, 
      months: 12 
    });

    expect(mockApiInstance.get).toHaveBeenCalledWith(
      '/sales/monthly-data',
      {
        params: { months: 12 },
        headers: { 'Content-Type': 'application/json' }
      }
    );

    expect(result).toEqual({
      data: mockResponse.data.data,
      status: 200,
      message: 'Success',
      success: true
    });
  });

  it('deve buscar dados mensais com parâmetros mínimos', async () => {
    const mockResponse = {
      data: {
        message: 'Success',
        data: {
          sales: [],
          purchases: [],
        },
        success: true
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    };

    mockApi.mockReturnValue(mockApiInstance as any);

    const result = await remoteGetMonthlySalesAndPurchases.getMonthlySalesAndPurchases({ farm_id: 1 });

    expect(mockApiInstance.get).toHaveBeenCalledWith(
      '/sales/monthly-data',
      {
        params: {},
        headers: { 'Content-Type': 'application/json' }
      }
    );

    expect(result).toEqual({
      data: { sales: [], purchases: [] },
      status: 200,
      message: 'Success',
      success: true
    });
  });

  it('deve lidar com resposta vazia', async () => {
    const mockResponse = {
      data: {
        message: 'No data',
        data: null,
        success: true
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    };

    mockApi.mockReturnValue(mockApiInstance as any);

    const result = await remoteGetMonthlySalesAndPurchases.getMonthlySalesAndPurchases({ farm_id: 1 });

    expect(result).toEqual({
      data: { sales: [], purchases: [] },
      status: 200,
      message: 'No data',
      success: true
    });
  });

  it('deve lidar com erro de rede', async () => {
    const mockError = new AxiosError('Network Error');
    mockError.response = {
      data: { message: 'Network Error' },
      status: 500,
      statusText: 'Internal Server Error',
      headers: {},
      config: {} as any
    };

    const mockApiInstance = {
      get: vi.fn().mockRejectedValue(mockError)
    };

    mockApi.mockReturnValue(mockApiInstance as any);

    await expect(remoteGetMonthlySalesAndPurchases.getMonthlySalesAndPurchases({ farm_id: 1 }))
      .rejects.toThrow('Network Error');
  });

  it('deve lidar com erro sem resposta', async () => {
    const mockError = new AxiosError('Request failed');
    mockError.response = undefined;

    const mockApiInstance = {
      get: vi.fn().mockRejectedValue(mockError)
    };

    mockApi.mockReturnValue(mockApiInstance as any);

    await expect(remoteGetMonthlySalesAndPurchases.getMonthlySalesAndPurchases({ farm_id: 1 }))
      .rejects.toThrow('Erro ao buscar dados mensais de vendas e compras');
  });

  it('deve lidar com erro desconhecido', async () => {
    const mockApiInstance = {
      get: vi.fn().mockRejectedValue(new Error('Unknown error'))
    };

    mockApi.mockReturnValue(mockApiInstance as any);

    await expect(remoteGetMonthlySalesAndPurchases.getMonthlySalesAndPurchases({ farm_id: 1 }))
      .rejects.toThrow('Erro desconhecido ao buscar dados mensais de vendas e compras');
  });

  it('deve usar mensagem padrão quando não houver mensagem na resposta', async () => {
    const mockResponse = {
      data: {
        data: {
          sales: [],
          purchases: [],
        },
        success: true
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    };

    mockApi.mockReturnValue(mockApiInstance as any);

    const result = await remoteGetMonthlySalesAndPurchases.getMonthlySalesAndPurchases({ farm_id: 1 });

    expect(result.message).toBe('dashboard.monthlySalesDataRetrievedSuccessfully');
  });

  it('deve construir parâmetros de query corretamente quando months for fornecido', async () => {
    const mockResponse = {
      data: {
        message: 'Success',
        data: {
          sales: [],
          purchases: [],
        },
        success: true
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    };

    mockApi.mockReturnValue(mockApiInstance as any);

    await remoteGetMonthlySalesAndPurchases.getMonthlySalesAndPurchases({ 
      farm_id: 2, 
      months: 6 
    });

    expect(mockApiInstance.get).toHaveBeenCalledWith(
      '/sales/monthly-data',
      {
        params: { months: 6 },
        headers: { 'Content-Type': 'application/json' }
      }
    );
  });
});

