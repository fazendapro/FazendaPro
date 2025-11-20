import { RemoteGetTopMilkProducers } from '../../data/usecases/remote-get-top-milk-producers';

export const GetTopMilkProducersFactory = () => {
  return new RemoteGetTopMilkProducers();
};
