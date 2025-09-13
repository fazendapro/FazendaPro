import { RemoteGetMilkProductions } from '../../data/usecases/remote-get-milk-productions'

export const getMilkProductionsFactory = () => {
  return new RemoteGetMilkProductions(
    undefined, // domain - pode ser passado como parâmetro se necessário
  )
}
