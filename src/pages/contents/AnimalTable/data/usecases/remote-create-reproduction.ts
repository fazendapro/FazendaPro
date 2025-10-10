import { api } from '../../../../../components/services/axios/api';
import { CreateReproductionRequest, Reproduction } from '../../domain/model/reproduction';

export const remoteCreateReproduction = async (data: CreateReproductionRequest): Promise<Reproduction> => {
  try {
    const response = await api().post('/api/v1/reproductions', data);
    return response.data.data;
  } catch (error) {
    throw error;
  }
};
