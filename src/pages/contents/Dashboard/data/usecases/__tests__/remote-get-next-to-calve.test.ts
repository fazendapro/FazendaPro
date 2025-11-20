import { describe, it, expect, vi, beforeEach } from 'vitest';
import { RemoteGetNextToCalve } from '../remote-get-next-to-calve';
import { api } from '../../../../../../components';
import { AxiosError } from 'axios';

vi.mock('../../../../../../components', () => ({
  api: vi.fn()
}));

vi.mock('i18next', () => ({
  t: vi.fn((key: string) => key)
}));

const mockApi = vi.mocked(api);

describe('RemoteGetNextToCalve', () => {
  let remoteGetNextToCalve: RemoteGetNextToCalve;

  beforeEach(() => {
    vi.clearAllMocks();
    remoteGetNextToCalve = new RemoteGetNextToCalve();
  });

  it('deve buscar próximas vacas a parir com sucesso', async () => {
    const mockResponse = {
      data: {
        message: 'Success',
        data: [
          {
            id: 1,
            animal_name: 'Tata Salt',
            ear_tag_number_local: 123,
            photo: 'src/assets/images/mocked/cows/tata.png',
            pregnancy_date: '2024-01-01',
            expected_birth_date: '2024-10-01',
            days_until_birth: 15,
            status: 'Alto'
          }
        ]
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    } as any;

    mockApi.mockReturnValue(mockApiInstance);

    const result = await remoteGetNextToCalve.getNextToCalve({ farm_id: 1 });

    expect(mockApiInstance.get).toHaveBeenCalledWith(
      '/reproductions/next-to-calve',
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
        data: null
      },
      status: 200
    };

    const mockApiInstance = {
      get: vi.fn().mockResolvedValue(mockResponse)
    } as any;

    mockApi.mockReturnValue(mockApiInstance);

    const result = await remoteGetNextToCalve.getNextToCalve({ farm_id: 1 });

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
      config: {} as any
    };

    const mockApiInstance = {
      get: vi.fn().mockRejectedValue(mockError)
    } as any;

    mockApi.mockReturnValue(mockApiInstance);

    await expect(remoteGetNextToCalve.getNextToCalve({ farm_id: 1 }))
      .rejects.toThrow('Network Error');
  });

  it('deve lidar com erro sem resposta', async () => {
    const mockError = new AxiosError('Request failed');
    mockError.response = undefined;

    const mockApiInstance = {
      get: vi.fn().mockRejectedValue(mockError)
    } as any;

    mockApi.mockReturnValue(mockApiInstance);

    await expect(remoteGetNextToCalve.getNextToCalve({ farm_id: 1 }))
      .rejects.toThrow('Erro ao buscar próximas vacas a parir');
  });

  it('deve lidar com erro desconhecido', async () => {
    const mockApiInstance = {
      get: vi.fn().mockRejectedValue(new Error('Unknown error'))
    } as any;

    mockApi.mockReturnValue(mockApiInstance);

    await expect(remoteGetNextToCalve.getNextToCalve({ farm_id: 1 }))
      .rejects.toThrow('Erro desconhecido ao buscar próximas vacas a parir');
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
    } as any;

    mockApi.mockReturnValue(mockApiInstance);

    const result = await remoteGetNextToCalve.getNextToCalve({ farm_id: 1 });

    expect(result.message).toBe('dashboard.nextToCalveRetrievedSuccessfully');
  });
});
