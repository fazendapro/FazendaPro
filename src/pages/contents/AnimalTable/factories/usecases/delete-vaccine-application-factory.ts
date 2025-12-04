import { RemoteDeleteVaccineApplication } from '../../data/usecases/remote-delete-vaccine-application'

export const deleteVaccineApplicationFactory = () => {
  return new RemoteDeleteVaccineApplication(
    undefined,
    undefined
  )
}

