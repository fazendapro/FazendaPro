import { RemoteUpdateMilkProduction } from '../../data/usecases/remote-update-milk-production'

export const updateMilkProductionFactory = () => {
  return new RemoteUpdateMilkProduction(
    undefined,
    undefined
  )
}
