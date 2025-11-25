import { api } from '../../../../../components';
import { GetMonthlySalesAndPurchasesDomain, GetMonthlySalesAndPurchasesResponse } from '../../domain/usecases/get-monthly-sales-and-purchases-domain';
import { GetMonthlySalesAndPurchasesParams } from '../../types/dashboard.types';
import { AxiosError } from 'axios';
import { t } from 'i18next';

export class RemoteGetMonthlySalesAndPurchases implements GetMonthlySalesAndPurchasesDomain {

  async getMonthlySalesAndPurchases(params: GetMonthlySalesAndPurchasesParams): Promise<GetMonthlySalesAndPurchasesResponse> {
    try {
      const queryParams: Record<string, string | number> = {
        farmId: params.farm_id
      };

      if (params.months) {
        queryParams.months = params.months;
      }

      const { data, status } = await api().get(
        '/sales/monthly-data',
        {
          params: queryParams,
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      const { message, data: monthlyData, success } = data;
      
      return {
        data: monthlyData || { sales: [], purchases: [] },
        status,
        message: message || t('dashboard.monthlySalesDataRetrievedSuccessfully') || 'Dados mensais recuperados com sucesso', 
        success: success !== undefined ? success : true
      };
    } catch (error) {
      if (error instanceof AxiosError) {
        throw new Error(error.response?.data?.message || 'Erro ao buscar dados mensais de vendas e compras');
      }
      throw new Error('Erro desconhecido ao buscar dados mensais de vendas e compras');
    }
  }
}


