import { api } from '../../../../../components/services/axios/api';
import { Weight } from '../../domain/model/weight';

export const remoteGetWeightByAnimal = async (animalId: number): Promise<Weight | null> => {
  try {
    const response = await api().get(`/weights/animal/${animalId}`);
    return response.data.data;
  } catch (error: any) {
    if (error.response?.status === 404) {
      return null;
    }
    throw error;
  }
};

