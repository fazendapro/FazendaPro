import { RemoteCreateMilkProduction } from '../../data/usecases/remote-create-milk-production'

export const createMilkProductionFactory = () => {
  return new RemoteCreateMilkProduction(
    undefined, // domain - pode ser passado como parâmetro se necessário
    undefined // csrfToken - temporariamente undefined até implementar o contexto
  )
}
