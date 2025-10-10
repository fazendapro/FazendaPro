import { api } from '../../../../../components';
import { GetTopMilkProducersDomain, GetTopMilkProducersResponse } from '../../domain/usecases/get-top-milk-producers-domain';
import { GetTopMilkProducersParams } from '../../types/dashboard.types';
import { AxiosError } from 'axios';
import { t } from 'i18next';

export class RemoteGetTopMilkProducers implements GetTopMilkProducersDomain {

  async getTopMilkProducers(params: GetTopMilkProducersParams): Promise<GetTopMilkProducersResponse> {
    try {
      const queryParams: Record<string, string | number> = {
        farmId: params.farm_id
      };

      if (params.limit) {
        queryParams.limit = params.limit;
      }

      if (params.period_days) {
        queryParams.periodDays = params.period_days;
      }

      const { data, status } = await api().get(
        '/api/v1/milk-collections/top-producers',
        {
          params: queryParams,
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      const { message, data: topProducersData } = data;
      return {
        data: topProducersData || [],
        status,
        message: message || t('dashboard.topMilkProducersRetrievedSuccessfully'), 
        success: true
      };
    } catch (error) {
      if (error instanceof AxiosError) {
        throw new Error(error.response?.data?.message || 'Erro ao buscar maiores produtoras de leite');
      }
      throw new Error('Erro desconhecido ao buscar maiores produtoras de leite');
    }
  }
}
