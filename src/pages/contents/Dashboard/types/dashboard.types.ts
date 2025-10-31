export interface NextToCalveAnimal {
  id: number;
  animal_name: string;
  ear_tag_number_local: number;
  photo?: string;
  pregnancy_date: string;
  expected_birth_date: string;
  days_until_birth: number;
  status: 'Alto' | 'Médio' | 'Baixo';
}

export interface GetNextToCalveParams {
  farm_id: number;
}

export interface GetNextToCalveResponse {
  data: NextToCalveAnimal[];
  success: boolean;
  message: string;
  status: number;
}

export interface TopMilkProducer {
  id: number;
  animal_name: string;
  ear_tag_number_local: number;
  photo?: string;
  total_production: number;
  average_daily_production: number;
  fat_content: number;
  last_collection_date: string;
  days_in_lactation: number;
}

export interface GetTopMilkProducersParams {
  farm_id: number;
  limit?: number;
  period_days?: number;
}

export interface GetTopMilkProducersResponse {
  data: TopMilkProducer[];
  success: boolean;
  message: string;
  status: number;
}

export interface MonthlyDataPoint {
  month: string;
  year: number;
  sales?: number;
  total?: number;
  count?: number;
}

export interface MonthlySalesAndPurchasesData {
  sales: MonthlyDataPoint[];
  purchases: MonthlyDataPoint[];
}

export interface GetMonthlySalesAndPurchasesParams {
  farm_id: number;
  months?: number;
}

export interface GetMonthlySalesAndPurchasesResponse {
  data: MonthlySalesAndPurchasesData;
  success: boolean;
  message: string;
  status: number;
}
