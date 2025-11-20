export interface Farm {
  id: number;
  name: string;
  location?: string;
  created_at: string;
  updated_at: string;
}

export interface GetFarmParams {
  farm_id: number;
} 