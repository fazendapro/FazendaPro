import { BaseHttpResponse } from '../../../../../components/services/axios/base-https-response';
import { CreateAnimalParams } from '../../types/type';

export interface CreateAnimalResponse extends BaseHttpResponse<{
  id: string;
  animal_name: string;
  ear_tag_number_local: string;
  ear_tag_number_global: string;
  farm_id: number;
  type: string;
  sex: string;
  breed: string;
  birth_date: string;
  created_at: string;
  updated_at: string;
}> {
  success: boolean;
  message: string;
  status: number;
}

export interface CreateAnimalDomain {
  create(params: CreateAnimalParams): Promise<CreateAnimalResponse>;
}