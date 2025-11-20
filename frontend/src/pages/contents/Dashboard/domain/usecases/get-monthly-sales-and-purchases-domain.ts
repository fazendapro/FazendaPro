import { BaseHttpResponse } from '../../../../../components/services/axios/base-https-response';
import { MonthlySalesAndPurchasesData, GetMonthlySalesAndPurchasesParams } from '../../types/dashboard.types';

export interface GetMonthlySalesAndPurchasesResponse extends BaseHttpResponse<MonthlySalesAndPurchasesData> {
  success: boolean;
  message: string;
  status: number;
}

export interface GetMonthlySalesAndPurchasesDomain {
  getMonthlySalesAndPurchases(params: GetMonthlySalesAndPurchasesParams): Promise<GetMonthlySalesAndPurchasesResponse>;
}


