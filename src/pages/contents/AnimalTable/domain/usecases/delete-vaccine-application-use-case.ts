export interface DeleteVaccineApplicationUseCase {
  deleteVaccineApplication: (id: number) => Promise<void>
}

