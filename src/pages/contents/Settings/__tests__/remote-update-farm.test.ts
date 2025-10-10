import { describe, it, expect, vi, beforeEach } from 'vitest'
import { RemoteUpdateFarm, RemoteGetFarm } from '../data/usecases/remote-update-farm'
import { api } from '../../../../components/services/axios/api'

vi.mock('../../../../components/services/axios/api', () => ({
  api: {
    put: vi.fn(),
    get: vi.fn(),
  },
}))

const mockApi = api as unknown as {
  put: ReturnType<typeof vi.fn>
  get: ReturnType<typeof vi.fn>
}

describe('RemoteUpdateFarm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('update', () => {
    it('deve atualizar fazenda com sucesso', async () => {
      const mockResponse = {
        data: {
          success: true,
          message: 'Fazenda atualizada com sucesso',
          data: {
            ID: 1,
            Logo: 'data:image/jpeg;base64,test',
            CompanyID: 1,
            Company: {
              ID: 1,
              CompanyName: 'Empresa Teste',
              Location: 'Rua Teste, 123',
              FarmCNPJ: '12345678000199'
            },
            CreatedAt: '2021-01-01T00:00:00Z',
            UpdatedAt: '2021-01-01T00:00:00Z'
          }
        }
      }

      mockApi.put.mockResolvedValue(mockResponse)

      const remoteUpdateFarm = new RemoteUpdateFarm()
      const result = await remoteUpdateFarm.update(1, { logo: 'data:image/jpeg;base64,test' })

      expect(mockApi.put).toHaveBeenCalledWith('/api/v1/farm', {
        logo: 'data:image/jpeg;base64,test'
      })
      expect(result.success).toBe(true)
      expect(result.message).toBe('Fazenda atualizada com sucesso')
      expect(result.data).toBeDefined()
    })

    it('deve tratar erro quando a API retorna erro', async () => {
      const mockError = new Error('Erro na API')
      mockApi.put.mockRejectedValue(mockError)

      const remoteUpdateFarm = new RemoteUpdateFarm()

      await expect(remoteUpdateFarm.update(1, { logo: 'test' })).rejects.toThrow('Erro na API')
    })

    it('deve tratar erro quando a resposta não é bem-sucedida', async () => {
      const mockResponse = {
        data: {
          success: false,
          message: 'Erro ao atualizar fazenda',
          data: null
        }
      }

      mockApi.put.mockResolvedValue(mockResponse)

      const remoteUpdateFarm = new RemoteUpdateFarm()

      await expect(remoteUpdateFarm.update(1, { logo: 'test' })).rejects.toThrow('Erro ao atualizar fazenda')
    })
  })
})

describe('RemoteGetFarm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('get', () => {
    it('deve buscar fazenda com sucesso', async () => {
      const mockResponse = {
        data: {
          success: true,
          message: 'Fazenda recuperada com sucesso',
          data: {
            ID: 1,
            Logo: 'data:image/jpeg;base64,test',
            CompanyID: 1,
            Company: {
              ID: 1,
              CompanyName: 'Empresa Teste',
              Location: 'Rua Teste, 123',
              FarmCNPJ: '12345678000199'
            },
            CreatedAt: '2021-01-01T00:00:00Z',
            UpdatedAt: '2021-01-01T00:00:00Z'
          }
        }
      }

      mockApi.get.mockResolvedValue(mockResponse)

      const remoteGetFarm = new RemoteGetFarm()
      const result = await remoteGetFarm.get(1)

      expect(mockApi.get).toHaveBeenCalledWith('/api/v1/farm')
      expect(result.success).toBe(true)
      expect(result.message).toBe('Fazenda recuperada com sucesso')
      expect(result.data).toBeDefined()
    })

    it('deve tratar erro quando a API retorna erro', async () => {
      const mockError = new Error('Erro na API')
      mockApi.get.mockRejectedValue(mockError)

      const remoteGetFarm = new RemoteGetFarm()

      await expect(remoteGetFarm.get(1)).rejects.toThrow('Erro na API')
    })

    it('deve tratar erro quando a resposta não é bem-sucedida', async () => {
      const mockResponse = {
        data: {
          success: false,
          message: 'Fazenda não encontrada',
          data: null
        }
      }

      mockApi.get.mockResolvedValue(mockResponse)

      const remoteGetFarm = new RemoteGetFarm()

      await expect(remoteGetFarm.get(1)).rejects.toThrow('Fazenda não encontrada')
    })
  })
})
