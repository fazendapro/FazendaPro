export interface Weight {
  id: number;
  animal_id: number;
  animal_name?: string;
  ear_tag?: number;
  date: string;
  animal_weight: number;
  created_at: string;
  updated_at: string;
}

export interface CreateOrUpdateWeightRequest {
  animal_id: number;
  date: string;
  animal_weight: number;
}

export interface UpdateWeightRequest {
  id: number;
  animal_id: number;
  date: string;
  animal_weight: number;
}

