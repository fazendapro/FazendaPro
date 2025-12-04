import { api } from '../../../../../components/services/axios/api';
import { CreateOrUpdateWeightRequest, Weight } from '../../domain/model/weight';

export const remoteCreateOrUpdateWeight = async (data: CreateOrUpdateWeightRequest): Promise<Weight> => {
  const response = await api().post('/weights', data);
  return response.data.data;
};

