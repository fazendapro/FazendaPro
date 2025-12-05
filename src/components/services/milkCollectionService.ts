import { api } from '../../config/api';
import { MilkCollection } from '../../types/milk-collection';

interface MilkCollectionApiData {
  id?: number;
  ID?: number;
  animal_id?: number;
  AnimalID?: number;
  farm_id?: number;
  FarmID?: number;
  collection_date?: string;
  date?: string;
  Date?: string;
  quantity?: number;
  liters?: number;
  Liters?: number;
  quality?: string;
  Quality?: string;
  notes?: string;
  Notes?: string;
  observations?: string;
  Observations?: string;
  created_at?: string;
  CreatedAt?: string;
  updated_at?: string;
  UpdatedAt?: string;
  animal?: unknown;
  Animal?: unknown;
}

interface MilkCollectionsResponse {
  success: boolean;
  data: MilkCollectionApiData | MilkCollectionApiData[];
  message?: string;
}

export const milkCollectionService = {
  async getMilkCollectionsByAnimal(animalId: number): Promise<MilkCollection[]> {
    try {
      const response = await api.get<MilkCollectionsResponse>(`/milk-collections/animal/${animalId}`);
      
      if (response.data.success && response.data.data) {
        const data = Array.isArray(response.data.data) ? response.data.data : [response.data.data];
        return data.map((mc: MilkCollectionApiData) => {
          const animal = mc.animal || mc.Animal;
          const animalData = animal && typeof animal === 'object' && 'id' in animal
            ? {
                id: (animal as { id?: number }).id || 0,
                animal_name: (animal as { animal_name?: string }).animal_name || '',
                ear_tag_number_local: (animal as { ear_tag_number_local?: number }).ear_tag_number_local || 0
              }
            : undefined;

          return {
            id: mc.id || mc.ID || 0,
            animal_id: mc.animal_id || mc.AnimalID || 0,
            farm_id: mc.farm_id || mc.FarmID || 0,
            collection_date: mc.collection_date || mc.date || mc.Date || '',
            quantity: mc.quantity || mc.liters || mc.Liters || 0,
            quality: mc.quality || mc.Quality,
            notes: mc.notes || mc.Notes || mc.observations || mc.Observations,
            created_at: mc.created_at || mc.CreatedAt || '',
            updated_at: mc.updated_at || mc.UpdatedAt || '',
            animal: animalData
          };
        });
      }
      
      return [];
    } catch (error) {
      console.error('Erro ao buscar ordenhas do animal:', error);
      return [];
    }
  }
};

