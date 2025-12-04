import { GetVaccineApplicationsUseCase } from '../../domain/usecases/get-vaccine-applications-use-case'
import { VaccineApplication, VaccineApplicationFilters } from '../../domain/model/vaccine-application'
import { api } from '../../../../../components/services/axios/api'

export class RemoteGetVaccineApplications implements GetVaccineApplicationsUseCase {
  constructor(
    private readonly domain?: string
  ) {}

  async getVaccineApplications(farmId: number, filters?: VaccineApplicationFilters): Promise<VaccineApplication[]> {
    const params: Record<string, string> = {}

    if (filters?.startDate) {
      params.start_date = filters.startDate
    }
    if (filters?.endDate) {
      params.end_date = filters.endDate
    }

    const response = await api(this.domain).get(`/vaccine-applications/farm/${farmId}`, { params })

    if (response.data.success) {
      return Array.isArray(response.data.data) ? response.data.data : []
    } else {
      throw new Error(response.data.message || 'Erro ao carregar aplicações de vacinas')
    }
  }
}

