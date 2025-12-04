import { GetVaccinesUseCase } from '../../domain/usecases/get-vaccines-use-case'
import { Vaccine } from '../../domain/model/vaccine'
import { api } from '../../../../../components/services/axios/api'

export class RemoteGetVaccines implements GetVaccinesUseCase {
  constructor(
    private readonly domain?: string
  ) {}

  async getVaccines(farmId: number): Promise<Vaccine[]> {
    const response = await api(this.domain).get(`/vaccines/farm/${farmId}`)

    if (response.data.success) {
      return Array.isArray(response.data.data) ? response.data.data : []
    } else {
      throw new Error(response.data.message || 'Erro ao carregar vacinas')
    }
  }
}

