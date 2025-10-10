import { CreateMilkProductionUseCase } from '../../domain/usecases/create-milk-production-use-case'
import { MilkProduction, CreateMilkProductionRequest } from '../../domain/model/milk-production'
import { api } from '../../../../../components/services/axios/api'

export class RemoteCreateMilkProduction implements CreateMilkProductionUseCase {
  constructor(
    private readonly domain?: string,
    private readonly csrfToken?: string
  ) {}

  async createMilkProduction(data: CreateMilkProductionRequest): Promise<MilkProduction> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    }
    
    if (this.csrfToken) {
      headers['X-CSRF-Token'] = this.csrfToken
    }
    
    const requestData = {
      ...data,
      liters: Number(data.liters)
    }
    
    const response = await api(this.domain).post('/api/v1/milk-collections', requestData, {
      headers
    })
    
    if (response.data.success) {
      return response.data.data
    } else {
      throw new Error(response.data.message || 'Erro ao criar registro de produção de leite')
    }
  }
}
