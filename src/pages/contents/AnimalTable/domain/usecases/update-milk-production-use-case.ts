import { MilkProduction, CreateMilkProductionRequest } from '../model/milk-production'

export interface UpdateMilkProductionRequest extends CreateMilkProductionRequest {
  id: number
}

export interface UpdateMilkProductionUseCase {
  updateMilkProduction: (data: UpdateMilkProductionRequest) => Promise<MilkProduction>
}
