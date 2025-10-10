import { BaseHttpResponse } from '../../../../../components/services/axios/base-https-response';
import { FarmData, BackendFarmData, UpdateFarmParams } from '../../types/farm-types';

export interface UpdateFarmResponse extends BaseHttpResponse<FarmData> {
  success: boolean;
  message: string;
  status: number;
}

export interface GetFarmResponse extends BaseHttpResponse<BackendFarmData> {
  success: boolean;
  message: string;
  status: number;
}

export interface UpdateFarmDomain {
  update(farmId: number, params: UpdateFarmParams): Promise<UpdateFarmResponse>;
}

export interface GetFarmDomain {
  get(farmId: number): Promise<GetFarmResponse>;
}
