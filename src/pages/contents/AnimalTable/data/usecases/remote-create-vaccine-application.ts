import { CreateVaccineApplicationUseCase } from '../../domain/usecases/create-vaccine-application-use-case'
import { VaccineApplication, CreateVaccineApplicationRequest } from '../../domain/model/vaccine-application'
import { api } from '../../../../../components/services/axios/api'

export class RemoteCreateVaccineApplication implements CreateVaccineApplicationUseCase {
  constructor(
    private readonly domain?: string,
    private readonly csrfToken?: string
  ) {}

  async createVaccineApplication(data: CreateVaccineApplicationRequest): Promise<VaccineApplication> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    }
    
    if (this.csrfToken) {
      headers['X-CSRF-Token'] = this.csrfToken
    }
    
    const response = await api(this.domain).post('/vaccine-applications', data, {
      headers
    })
    
    if (response.data.success) {
      return response.data.data
    } else {
      throw new Error(response.data.message || 'Erro ao criar aplicação de vacina')
    }
  }
}

