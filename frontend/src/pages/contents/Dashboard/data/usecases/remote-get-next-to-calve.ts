import { api } from '../../../../../components';
import { GetNextToCalveDomain, GetNextToCalveResponse } from '../../domain/usecases/get-next-to-calve-domain';
import { GetNextToCalveParams } from '../../types/dashboard.types';
import { AxiosError } from 'axios';
import { t } from 'i18next';

export class RemoteGetNextToCalve implements GetNextToCalveDomain {

  async getNextToCalve(params: GetNextToCalveParams): Promise<GetNextToCalveResponse> {
    try {
      const { data, status } = await api().get(
        '/reproductions/next-to-calve',
        {
          params: {
            farmId: params.farm_id
          },
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      const { message, data: nextToCalveData } = data;
      return {
        data: nextToCalveData || [],
        status,
        message: message || t('dashboard.nextToCalveRetrievedSuccessfully'), 
        success: true
      };
    } catch (error) {
      if (error instanceof AxiosError) {
        throw new Error(error.response?.data?.message || 'Erro ao buscar próximas vacas a parir');
      }
      throw new Error('Erro desconhecido ao buscar próximas vacas a parir');
    }
  }
}
