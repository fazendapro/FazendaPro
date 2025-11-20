import { BaseHttpResponse } from '../../../../../components/services/axios/base-https-response';
import { Animal, GetAnimalsByFarmParams } from '../../types/type';

export interface GetAnimalsByFarmResponse extends BaseHttpResponse<Animal[]> {
  success: boolean;
  message: string;
  status: number;
}

export interface GetAnimalsByFarmDomain {
  getAnimalsByFarm(params: GetAnimalsByFarmParams): Promise<GetAnimalsByFarmResponse>;
} 