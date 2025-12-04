import { api } from '../../../../../components/services/axios/api';
import { Reproduction } from '../../domain/model/reproduction';

export interface GetReproductionsByFarmParams {
  farm_id: number;
  page?: number;
  limit?: number;
}

export interface GetReproductionsByFarmResponse {
  reproductions: Reproduction[];
  total: number;
  page: number;
  limit: number;
}

export const remoteGetReproductionsByFarm = async (params: GetReproductionsByFarmParams): Promise<GetReproductionsByFarmResponse> => {
  const requestParams: Record<string, string | number> = {
    farmId: params.farm_id
  };

  if (params.page !== undefined) {
    requestParams.page = params.page;
  }

  if (params.limit !== undefined) {
    requestParams.limit = params.limit;
  }

  const { data } = await api().get('/reproductions/farm', {
    params: requestParams,
    headers: {
      'Content-Type': 'application/json'
    }
  });

  const { message, data: responseData } = data;

  let paginatedData;
  if (Array.isArray(responseData)) {
    // Backward compatibility: if response is array, wrap it
    paginatedData = {
      reproductions: responseData,
      total: responseData.length,
      page: 1,
      limit: responseData.length
    };
  } else {
    // Paginated response
    paginatedData = responseData || {
      reproductions: [],
      total: 0,
      page: 1,
      limit: 10
    };
  }

  return paginatedData;
};
