import { api } from '../../../../../components';
import { GetOverviewStatsDomain, GetOverviewStatsResponse } from '../../domain/usecases/get-overview-stats-domain';
import { GetOverviewStatsParams } from '../../types/dashboard.types';
import { AxiosError } from 'axios';
import { t } from 'i18next';

export class RemoteGetOverviewStats implements GetOverviewStatsDomain {

  async getOverviewStats(params: GetOverviewStatsParams): Promise<GetOverviewStatsResponse> {
    try {
      const { data, status } = await api().get(
        '/sales/overview',
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      const { message, data: overviewData, success } = data;
      
      return {
        data: overviewData || { males_count: 0, females_count: 0, total_sold: 0, total_revenue: 0 },
        status,
        message: message || t('dashboard.overviewStatsRetrievedSuccessfully') || 'Estatísticas gerais recuperadas com sucesso',
        success: success !== undefined ? success : true
      };
    } catch (error) {
      if (error instanceof AxiosError) {
        throw new Error(error.response?.data?.message || 'Erro ao buscar estatísticas gerais');
      }
      throw new Error('Erro desconhecido ao buscar estatísticas gerais');
    }
  }
}

