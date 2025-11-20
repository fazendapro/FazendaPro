import { RemoteCreateMilkProduction } from '../../data/usecases/remote-create-milk-production'

export const createMilkProductionFactory = () => {
  return new RemoteCreateMilkProduction(
    undefined,
    undefined
  )
}
