import { GetVaccineApplicationsUseCase } from '../../domain/usecases/get-vaccine-applications-use-case'
import { VaccineApplication, VaccineApplicationFilters } from '../../domain/model/vaccine-application'
import { api } from '../../../../../components/services/axios/api'

export interface GetVaccineApplicationsResponse {
  vaccine_applications: VaccineApplication[];
  total: number;
  page: number;
  limit: number;
}

export class RemoteGetVaccineApplications implements GetVaccineApplicationsUseCase {
  constructor(
    private readonly domain?: string
  ) {}

  async getVaccineApplications(farmId: number, filters?: VaccineApplicationFilters): Promise<VaccineApplication[] | GetVaccineApplicationsResponse> {
    const params: Record<string, string | number> = {}

    if (filters?.startDate) {
      params.start_date = filters.startDate
    }
    if (filters?.endDate) {
      params.end_date = filters.endDate
    }
    if (filters?.page !== undefined) {
      params.page = filters.page
    }
    if (filters?.limit !== undefined) {
      params.limit = filters.limit
    }

    const response = await api(this.domain).get(`/vaccine-applications/farm/${farmId}`, { params })

    if (response.data.success) {
      const responseData = response.data.data
      
      if (responseData && typeof responseData === 'object' && 'vaccine_applications' in responseData) {
        return responseData as GetVaccineApplicationsResponse
      }
      
      return Array.isArray(responseData) ? responseData : []
    } else {
      throw new Error(response.data.message || 'Erro ao carregar aplicações de vacinas')
    }
  }
}

