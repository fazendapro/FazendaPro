import { UpdateFarmDomain, GetFarmDomain } from '../../domain/usecases/update-farm-domain';
import { RemoteUpdateFarm, RemoteGetFarm } from '../../data/usecases/remote-update-farm';

export class UpdateFarmFactory {
  static create(): UpdateFarmDomain {
    return new RemoteUpdateFarm();
  }
}

export class GetFarmFactory {
  static create(): GetFarmDomain {
    return new RemoteGetFarm();
  }
}
