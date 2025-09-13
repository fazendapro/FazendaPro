import { CreateMilkProductionUseCase } from '../../domain/usecases/create-milk-production-use-case'
import { MilkProduction, CreateMilkProductionRequest } from '../../domain/model/milk-production'
import { api } from '../../../../../components/services/axios/api'

export class RemoteCreateMilkProduction implements CreateMilkProductionUseCase {
  constructor(
    private readonly domain?: string,
    private readonly csrfToken?: string
  ) {}

  async createMilkProduction(data: CreateMilkProductionRequest): Promise<MilkProduction> {
    console.log('RemoteCreateMilkProduction - Creating with data:', data)
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    }
    
    if (this.csrfToken) {
      headers['X-CSRF-Token'] = this.csrfToken
    }
    
    // Garantir que liters seja um número
    const requestData = {
      ...data,
      liters: Number(data.liters)
    }
    
    console.log('Request headers:', headers)
    console.log('Request data:', requestData)
    console.log('API domain:', this.domain)
    
    const response = await api(this.domain).post('/api/v1/milk-collections', requestData, {
      headers
    })
    
    console.log('API response:', response.data)
    
    if (response.data.success) {
      return response.data.data
    } else {
      throw new Error(response.data.message || 'Erro ao criar registro de produção de leite')
    }
  }
}
