import { UpdateMilkProductionUseCase, UpdateMilkProductionRequest } from '../../domain/usecases/update-milk-production-use-case'
import { MilkProduction } from '../../domain/model/milk-production'
import { api } from '../../../../../components/services/axios/api'

export class RemoteUpdateMilkProduction implements UpdateMilkProductionUseCase {
  constructor(
    private readonly domain?: string,
    private readonly csrfToken?: string
  ) {}

  async updateMilkProduction(data: UpdateMilkProductionRequest): Promise<MilkProduction> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    }

    if (this.csrfToken) {
      headers['X-CSRF-Token'] = this.csrfToken
    }

    const requestData = {
      animal_id: data.animal_id,
      liters: Number(data.liters),
      date: data.date
    }

    const response = await api(this.domain).put(`/api/v1/milk-collections/${data.id}`, requestData, {
      headers
    })

    if (response.data.success) {
      return response.data.data
    } else {
      throw new Error(response.data.message || 'Erro ao atualizar registro de produção de leite')
    }
  }
}
