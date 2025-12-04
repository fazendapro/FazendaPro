import { MilkProduction, MilkProductionFilters } from '../model/milk-production'

export interface GetMilkProductionsParams {
  farmId: number
  filters?: MilkProductionFilters
  page?: number
  limit?: number
}

export interface GetMilkProductionsResponse {
  milk_collections: MilkProduction[]
  total: number
  page: number
  limit: number
}

export interface GetMilkProductionsUseCase {
  getMilkProductions: (params: GetMilkProductionsParams) => Promise<GetMilkProductionsResponse>
}
