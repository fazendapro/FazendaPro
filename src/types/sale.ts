export interface Sale {
  id: number;
  animal_id: number;
  farm_id: number;
  buyer_name: string;
  price: number;
  sale_date: string;
  notes?: string;
  created_at: string;
  updated_at: string;
  animal?: {
    id: number;
    animal_name: string;
    ear_tag_number_local: number;
    breed: string;
    type: string;
    sex: number;
    status: number;
  };
}

export interface CreateSaleRequest {
  animal_id: number;
  buyer_name: string;
  price: number;
  sale_date: string;
  notes?: string;
}

export interface UpdateSaleRequest {
  buyer_name: string;
  price: number;
  sale_date: string;
  notes?: string;
}

export interface SaleFilters {
  start_date?: string;
  end_date?: string;
  animal_id?: number;
}
