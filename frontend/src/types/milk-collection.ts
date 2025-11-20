export interface MilkCollection {
  id: number;
  animal_id: number;
  farm_id: number;
  collection_date: string;
  quantity: number;
  quality?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
  animal?: {
    id: number;
    animal_name: string;
    ear_tag_number_local: number;
  };
}
