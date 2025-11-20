export interface Reproduction {
  id: number;
  animal_id: number;
  farm_id: number;
  date: string;
  phase: string;
  notes?: string;
  created_at: string;
  updated_at: string;
  animal?: {
    id: number;
    animal_name: string;
    ear_tag_number_local: number;
  };
}
