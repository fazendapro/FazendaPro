import { RemoteGetVaccines } from '../../data/usecases/remote-get-vaccines'

export const getVaccinesFactory = () => {
  return new RemoteGetVaccines(
    undefined,
  )
}

