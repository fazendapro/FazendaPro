import { DeleteVaccineApplicationUseCase } from '../../domain/usecases/delete-vaccine-application-use-case'
import { api } from '../../../../../components/services/axios/api'

export class RemoteDeleteVaccineApplication implements DeleteVaccineApplicationUseCase {
  constructor(
    private readonly domain?: string,
    private readonly csrfToken?: string
  ) {}

  async deleteVaccineApplication(id: number): Promise<void> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    }
    
    if (this.csrfToken) {
      headers['X-CSRF-Token'] = this.csrfToken
    }
    
    const response = await api(this.domain).delete(`/vaccine-applications/${id}`, {
      headers
    })
    
    if (!response.data.success) {
      throw new Error(response.data.message || 'Erro ao deletar aplicação de vacina')
    }
  }
}

