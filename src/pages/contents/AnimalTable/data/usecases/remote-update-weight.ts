import { api } from '../../../../../components/services/axios/api';
import { UpdateWeightRequest, Weight } from '../../domain/model/weight';

export const remoteUpdateWeight = async (data: UpdateWeightRequest): Promise<Weight> => {
  const response = await api().put('/weights', data);
  return response.data.data;
};

