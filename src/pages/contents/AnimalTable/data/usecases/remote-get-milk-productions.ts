import { GetMilkProductionsUseCase, GetMilkProductionsParams, GetMilkProductionsResponse } from '../../domain/usecases/get-milk-productions-use-case'
import { api } from '../../../../../components/services/axios/api'

export class RemoteGetMilkProductions implements GetMilkProductionsUseCase {
  constructor(
    private readonly domain?: string
  ) {}

  async getMilkProductions(params: GetMilkProductionsParams): Promise<GetMilkProductionsResponse> {
    const requestParams: Record<string, string | number> = {}

    if (params.filters?.startDate) {
      requestParams.start_date = params.filters.startDate
    }
    if (params.filters?.endDate) {
      requestParams.end_date = params.filters.endDate
    }

    if (params.page !== undefined) {
      requestParams.page = params.page
    }

    if (params.limit !== undefined) {
      requestParams.limit = params.limit
    }

    const response = await api(this.domain).get(`/milk-collections/farm/${params.farmId}`, { params: requestParams })

    if (response.data.success) {
      const responseData = response.data.data

      if (responseData && typeof responseData === 'object' && 'milk_collections' in responseData) {
        return {
          milk_collections: responseData.milk_collections || [],
          total: responseData.total || 0,
          page: responseData.page || 1,
          limit: responseData.limit || 10
        }
      }

      if (Array.isArray(responseData)) {
        return {
          milk_collections: responseData,
          total: responseData.length,
          page: 1,
          limit: responseData.length
        }
      }

      return {
        milk_collections: [],
        total: 0,
        page: 1,
        limit: 10
      }
    } else {
      throw new Error(response.data.message || 'Erro ao carregar dados de produção de leite')
    }
  }
}
