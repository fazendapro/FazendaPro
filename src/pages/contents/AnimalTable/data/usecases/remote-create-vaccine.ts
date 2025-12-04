import { CreateVaccineUseCase } from '../../domain/usecases/create-vaccine-use-case'
import { Vaccine, CreateVaccineRequest } from '../../domain/model/vaccine'
import { api } from '../../../../../components/services/axios/api'

export class RemoteCreateVaccine implements CreateVaccineUseCase {
  constructor(
    private readonly domain?: string,
    private readonly csrfToken?: string
  ) {}

  async createVaccine(data: CreateVaccineRequest): Promise<Vaccine> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    }
    
    if (this.csrfToken) {
      headers['X-CSRF-Token'] = this.csrfToken
    }
    
    const response = await api(this.domain).post('/vaccines', data, {
      headers
    })
    
    if (response.data.success) {
      return response.data.data
    } else {
      throw new Error(response.data.message || 'Erro ao criar vacina')
    }
  }
}

