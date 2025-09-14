import { RemoteUpdateMilkProduction } from '../../data/usecases/remote-update-milk-production'

export const updateMilkProductionFactory = () => {
  return new RemoteUpdateMilkProduction(
    undefined, // domain - pode ser passado como parâmetro se necessário
    undefined // csrfToken - temporariamente undefined até implementar o contexto
  )
}
