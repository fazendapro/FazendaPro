import { RemoteCreateVaccine } from '../../data/usecases/remote-create-vaccine'

export const createVaccineFactory = () => {
  return new RemoteCreateVaccine(
    undefined,
    undefined
  )
}

