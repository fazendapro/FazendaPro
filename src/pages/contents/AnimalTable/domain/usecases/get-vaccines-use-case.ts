import { Vaccine } from '../model/vaccine'

export interface GetVaccinesUseCase {
  getVaccines: (farmId: number) => Promise<Vaccine[]>
}

