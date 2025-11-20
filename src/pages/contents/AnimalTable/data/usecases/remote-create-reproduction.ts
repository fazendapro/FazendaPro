import { api } from '../../../../../components/services/axios/api';
import { CreateReproductionRequest, Reproduction } from '../../domain/model/reproduction';

export const remoteCreateReproduction = async (data: CreateReproductionRequest): Promise<Reproduction> => {
  const response = await api().post('/reproductions', data);
  return response.data.data;
};
