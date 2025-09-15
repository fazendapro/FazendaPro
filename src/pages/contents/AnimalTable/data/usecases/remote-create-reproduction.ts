import { api } from '../../../../../components/services/axios/api';
import { CreateReproductionRequest, Reproduction } from '../../domain/model/reproduction';

export const remoteCreateReproduction = async (data: CreateReproductionRequest): Promise<Reproduction> => {
  console.log('remoteCreateReproduction - Sending data:', data);
  console.log('remoteCreateReproduction - API URL: /api/v1/reproductions');
  
  try {
    const response = await api().post('/api/v1/reproductions', data);
    console.log('remoteCreateReproduction - Response:', response);
    console.log('remoteCreateReproduction - Response data:', response.data);
    return response.data.data;
  } catch (error) {
    console.error('remoteCreateReproduction - Error:', error);
    throw error;
  }
};
