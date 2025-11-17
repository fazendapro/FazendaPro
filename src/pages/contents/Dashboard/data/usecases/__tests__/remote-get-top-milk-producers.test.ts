import { describe, it, expect, vi, beforeEach } from 'vitest';
import { RemoteGetTopMilkProducers } from '../remote-get-top-milk-producers';
import { api } from '../../../../../../components';
import { AxiosError, AxiosRequestConfig } from 'axios';

vi.mock('../../../../../../components', () => ({
  api: vi.fn()
}));

vi.mock('i18next', () => ({
  t: vi.fn((key: string) => key)
}));

const mockApi = vi.mocked(api);

describe('RemoteGetTopMilkProducers', () => {
  let remoteGetTopMilkProducers: RemoteGetTopMilkProducers;

  beforeEach(() => {
    vi.clearAllMocks();
    remoteGetTopMilkProducers = new RemoteGetTopMilkProducers();
  });

  it('deve buscar maiores produtoras de leite com sucesso', async () => {
    const mockResponse = {
      data: {
        message: 'Success',
        data: [
          {
            id: 1,
            animal_name: 'Surf Excel',
            ear_tag_number_local: 123,
            photo: 'src/assets/images/mocked/cows/surf.png',
            total_production: 900,
            average_daily_production: 30,
            fat_content: 3.5,
            last_collection_date: '2024-01-15',
            days_in_lactation: 120
          }
        ]
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    } as unknown as typeof api;

    mockApi.mockReturnValue(mockApiInstance);

    const result = await remoteGetTopMilkProducers.getTopMilkProducers({ 
      farm_id: 1, 
      limit: 10, 
      period_days: 30 
    });

    expect(mockApiInstance.get).toHaveBeenCalledWith(
      '/api/v1/milk-collections/top-producers',
      {
        params: { 
          farmId: 1, 
          limit: 10, 
          periodDays: 30 
        },
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

  it('deve buscar maiores produtoras com parâmetros mínimos', async () => {
    const mockResponse = {
      data: {
        message: 'Success',
        data: []
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    } as unknown as typeof api;

    mockApi.mockReturnValue(mockApiInstance);

    const result = await remoteGetTopMilkProducers.getTopMilkProducers({ farm_id: 1 });

    expect(mockApiInstance.get).toHaveBeenCalledWith(
      '/api/v1/milk-collections/top-producers',
      {
        params: { farmId: 1 },
        headers: { 'Content-Type': 'application/json' }
      }
    );

    expect(result).toEqual({
      data: [],
      status: 200,
      message: 'Success',
      success: true
    });
  });

  it('deve lidar com resposta vazia', async () => {
    const mockResponse = {
      data: {
        message: 'No data',
        data: null
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    } as unknown as typeof api;

    mockApi.mockReturnValue(mockApiInstance);

    const result = await remoteGetTopMilkProducers.getTopMilkProducers({ farm_id: 1 });

    expect(result).toEqual({
      data: [],
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
      config: {} as AxiosRequestConfig
    };

    const mockApiInstance = {
      get: vi.fn().mockRejectedValue(mockError)
    } as unknown as typeof api;

    mockApi.mockReturnValue(mockApiInstance);

    await expect(remoteGetTopMilkProducers.getTopMilkProducers({ farm_id: 1 }))
      .rejects.toThrow('Network Error');
  });

  it('deve lidar com erro sem resposta', async () => {
    const mockError = new AxiosError('Request failed');
    mockError.response = undefined;

    const mockApiInstance = {
      get: vi.fn().mockRejectedValue(mockError)
    } as unknown as typeof api;

    mockApi.mockReturnValue(mockApiInstance);

    await expect(remoteGetTopMilkProducers.getTopMilkProducers({ farm_id: 1 }))
      .rejects.toThrow('Erro ao buscar maiores produtoras de leite');
  });

  it('deve lidar com erro desconhecido', async () => {
    const mockApiInstance = {
      get: vi.fn().mockRejectedValue(new Error('Unknown error'))
    } as unknown as typeof api;

    mockApi.mockReturnValue(mockApiInstance);

    await expect(remoteGetTopMilkProducers.getTopMilkProducers({ farm_id: 1 }))
      .rejects.toThrow('Erro desconhecido ao buscar maiores produtoras de leite');
  });

  it('deve usar mensagem padrão quando não houver mensagem na resposta', async () => {
    const mockResponse = {
      data: {
        data: []
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    } as unknown as typeof api;

    mockApi.mockReturnValue(mockApiInstance);

    const result = await remoteGetTopMilkProducers.getTopMilkProducers({ farm_id: 1 });

    expect(result.message).toBe('dashboard.topMilkProducersRetrievedSuccessfully');
  });

  it('deve construir parâmetros de query corretamente', async () => {
    const mockResponse = {
      data: {
        message: 'Success',
        data: []
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    } as unknown as typeof api;

    mockApi.mockReturnValue(mockApiInstance);

    await remoteGetTopMilkProducers.getTopMilkProducers({ 
      farm_id: 2, 
      limit: 5, 
      period_days: 15 
    });

    expect(mockApiInstance.get).toHaveBeenCalledWith(
      '/api/v1/milk-collections/top-producers',
      {
        params: { 
          farmId: 2, 
          limit: 5, 
          periodDays: 15 
        },
        headers: { 'Content-Type': 'application/json' }
      }
    );
  });
});
