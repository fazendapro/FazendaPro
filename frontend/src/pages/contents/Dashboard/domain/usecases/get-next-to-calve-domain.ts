import { BaseHttpResponse } from '../../../../../components/services/axios/base-https-response';
import { NextToCalveAnimal, GetNextToCalveParams } from '../../types/dashboard.types';

export interface GetNextToCalveResponse extends BaseHttpResponse<NextToCalveAnimal[]> {
  success: boolean;
  message: string;
  status: number;
}

export interface GetNextToCalveDomain {
  getNextToCalve(params: GetNextToCalveParams): Promise<GetNextToCalveResponse>;
}
