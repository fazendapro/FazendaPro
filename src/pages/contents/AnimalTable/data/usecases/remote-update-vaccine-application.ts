import { UpdateVaccineApplicationUseCase } from '../../domain/usecases/update-vaccine-application-use-case'
import { VaccineApplication, UpdateVaccineApplicationRequest } from '../../domain/model/vaccine-application'
import { api } from '../../../../../components/services/axios/api'

export class RemoteUpdateVaccineApplication implements UpdateVaccineApplicationUseCase {
  constructor(
    private readonly domain?: string,
    private readonly csrfToken?: string
  ) {}

  async updateVaccineApplication(id: number, data: UpdateVaccineApplicationRequest): Promise<VaccineApplication> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    }
    
    if (this.csrfToken) {
      headers['X-CSRF-Token'] = this.csrfToken
    }
    
    const response = await api(this.domain).put(`/vaccine-applications/${id}`, data, {
      headers
    })
    
    if (response.data.success) {
      return response.data.data
    } else {
      throw new Error(response.data.message || 'Erro ao atualizar aplicação de vacina')
    }
  }
}

