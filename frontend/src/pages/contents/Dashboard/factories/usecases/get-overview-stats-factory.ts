import { RemoteGetOverviewStats } from '../../data/usecases/remote-get-overview-stats';
import { GetOverviewStatsDomain } from '../../domain/usecases/get-overview-stats-domain';

export const GetOverviewStatsFactory = (): GetOverviewStatsDomain => {
  return new RemoteGetOverviewStats();
};

