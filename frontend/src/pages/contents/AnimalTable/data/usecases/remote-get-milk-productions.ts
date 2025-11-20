import { GetMilkProductionsUseCase } from '../../domain/usecases/get-milk-productions-use-case'
import { MilkProduction, MilkProductionFilters } from '../../domain/model/milk-production'
import { api } from '../../../../../components/services/axios/api'

export class RemoteGetMilkProductions implements GetMilkProductionsUseCase {
  constructor(
    private readonly domain?: string
  ) {}

  async getMilkProductions(farmId: number, filters?: MilkProductionFilters): Promise<MilkProduction[]> {
    const params: Record<string, string> = {}

    if (filters?.startDate) {
      params.start_date = filters.startDate
    }
    if (filters?.endDate) {
      params.end_date = filters.endDate
    }

    const response = await api(this.domain).get(`/milk-collections/farm/${farmId}`, { params })

    if (response.data.success) {
      return response.data.data
    } else {
      throw new Error(response.data.message || 'Erro ao carregar dados de produção de leite')
    }
  }
}
