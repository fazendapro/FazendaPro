import { api } from '../../../../../components/services/axios/api';
import { UpdateReproductionPhaseRequest } from '../../domain/model/reproduction';

export const remoteUpdateReproductionPhase = async (data: UpdateReproductionPhaseRequest): Promise<void> => {
  await api().put('/reproductions/phase', data);
};
