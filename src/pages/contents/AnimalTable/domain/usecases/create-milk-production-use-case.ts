import { MilkProduction, CreateMilkProductionRequest } from '../model/milk-production'

export interface CreateMilkProductionUseCase {
  createMilkProduction: (data: CreateMilkProductionRequest) => Promise<MilkProduction>
}
