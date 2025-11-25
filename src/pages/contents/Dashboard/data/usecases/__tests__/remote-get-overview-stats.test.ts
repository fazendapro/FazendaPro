import { describe, it, expect, vi, beforeEach } from 'vitest';
import { RemoteGetOverviewStats } from '../remote-get-overview-stats';
import { api } from '../../../../../../components';
import { AxiosError, InternalAxiosRequestConfig } from 'axios';

vi.mock('../../../../../../components', () => ({
  api: vi.fn()
}));

vi.mock('i18next', () => ({
  t: vi.fn((key: string) => key)
}));

const mockApi = vi.mocked(api);

describe('RemoteGetOverviewStats', () => {
  let remoteGetOverviewStats: RemoteGetOverviewStats;

  beforeEach(() => {
    vi.clearAllMocks();
    remoteGetOverviewStats = new RemoteGetOverviewStats();
  });

  it('deve buscar estatísticas gerais com sucesso', async () => {
    const mockResponse = {
      data: {
        message: 'Success',
        data: {
          males_count: 10,
          females_count: 15,
          total_sold: 5,
          total_revenue: 15000.50
        },
        success: true
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    } as unknown as ReturnType<typeof api>;

    mockApi.mockReturnValue(mockApiInstance as unknown as ReturnType<typeof api>);

    const result = await remoteGetOverviewStats.getOverviewStats({ farm_id: 1 });

    expect(mockApiInstance.get).toHaveBeenCalledWith(
      '/sales/overview',
      {
        params: { farmId: 1 },
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
    } as unknown as ReturnType<typeof api>;

    mockApi.mockReturnValue(mockApiInstance as unknown as ReturnType<typeof api>);

    const result = await remoteGetOverviewStats.getOverviewStats({ farm_id: 1 });

    expect(result).toEqual({
      data: { males_count: 0, females_count: 0, total_sold: 0, total_revenue: 0 },
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
      config: {} as unknown as InternalAxiosRequestConfig
    };

    const mockApiInstance = {
      get: vi.fn().mockRejectedValue(mockError)
    } as unknown as ReturnType<typeof api>;

    mockApi.mockReturnValue(mockApiInstance as unknown as ReturnType<typeof api>);

    await expect(remoteGetOverviewStats.getOverviewStats({ farm_id: 1 }))
      .rejects.toThrow('Network Error');
  });

  it('deve lidar com erro sem resposta', async () => {
    const mockError = new AxiosError('Request failed');
    mockError.response = undefined;

    const mockApiInstance = {
      get: vi.fn().mockRejectedValue(mockError)
    } as unknown as ReturnType<typeof api>;

    mockApi.mockReturnValue(mockApiInstance as unknown as ReturnType<typeof api>);

    await expect(remoteGetOverviewStats.getOverviewStats({ farm_id: 1 }))
      .rejects.toThrow('Erro ao buscar estatísticas gerais');
  });

  it('deve lidar com erro desconhecido', async () => {
    const mockApiInstance = {
      get: vi.fn().mockRejectedValue(new Error('Unknown error'))
    };

    mockApi.mockReturnValue(mockApiInstance as unknown as ReturnType<typeof api>);

    await expect(remoteGetOverviewStats.getOverviewStats({ farm_id: 1 }))
      .rejects.toThrow('Erro desconhecido ao buscar estatísticas gerais');
  });

  it('deve usar mensagem padrão quando não houver mensagem na resposta', async () => {
    const mockResponse = {
      data: {
        data: {
          males_count: 10,
          females_count: 15,
          total_sold: 5,
          total_revenue: 15000.50
        },
        success: true
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    } as unknown as ReturnType<typeof api>;

    mockApi.mockReturnValue(mockApiInstance as unknown as ReturnType<typeof api>);

    const result = await remoteGetOverviewStats.getOverviewStats({ farm_id: 1 });

    expect(result.message).toBe('dashboard.overviewStatsRetrievedSuccessfully');
  });
});

