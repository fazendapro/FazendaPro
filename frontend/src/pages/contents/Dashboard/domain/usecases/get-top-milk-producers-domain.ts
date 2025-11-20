import { BaseHttpResponse } from '../../../../../components/services/axios/base-https-response';
import { TopMilkProducer, GetTopMilkProducersParams } from '../../types/dashboard.types';

export interface GetTopMilkProducersResponse extends BaseHttpResponse<TopMilkProducer[]> {
  success: boolean;
  message: string;
  status: number;
}

export interface GetTopMilkProducersDomain {
  getTopMilkProducers(params: GetTopMilkProducersParams): Promise<GetTopMilkProducersResponse>;
}
