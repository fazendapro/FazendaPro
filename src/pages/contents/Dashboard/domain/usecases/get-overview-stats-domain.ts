import { BaseHttpResponse } from '../../../../../components/services/axios/base-https-response';
import { OverviewStats, GetOverviewStatsParams } from '../../types/dashboard.types';

export interface GetOverviewStatsResponse extends BaseHttpResponse<OverviewStats> {
  success: boolean;
  message: string;
  status: number;
}

export interface GetOverviewStatsDomain {
  getOverviewStats(params: GetOverviewStatsParams): Promise<GetOverviewStatsResponse>;
}

