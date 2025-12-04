import { BaseHttpResponse } from '../../../../../components/services/axios/base-https-response';
import { GetAnimalsByFarmParams, GetAnimalsByFarmPaginatedResponse } from '../../types/type';

export interface GetAnimalsByFarmResponse extends BaseHttpResponse<GetAnimalsByFarmPaginatedResponse> {
  success: boolean;
  message: string;
  status: number;
}

export interface GetAnimalsByFarmDomain {
  getAnimalsByFarm(params: GetAnimalsByFarmParams): Promise<GetAnimalsByFarmResponse>;
} 