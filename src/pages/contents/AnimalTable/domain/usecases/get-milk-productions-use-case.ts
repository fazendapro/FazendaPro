import { MilkProduction, MilkProductionFilters } from '../model/milk-production'

export interface GetMilkProductionsUseCase {
  getMilkProductions: (farmId: number, filters?: MilkProductionFilters) => Promise<MilkProduction[]>
}
