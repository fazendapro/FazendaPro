import { api } from '../../config/api';
import { MilkCollection } from '../../types/milk-collection';

interface MilkCollectionsResponse {
  success: boolean;
  data: MilkCollection[];
  message?: string;
}

export const milkCollectionService = {
  async getMilkCollectionsByAnimal(animalId: number): Promise<MilkCollection[]> {
    try {
      const response = await api.get<MilkCollectionsResponse>(`/milk-collections/animal/${animalId}`);
      
      if (response.data.success && response.data.data) {
        const data = Array.isArray(response.data.data) ? response.data.data : [response.data.data];
        return data.map((mc: any) => ({
          id: mc.id || mc.ID,
          animal_id: mc.animal_id || mc.AnimalID,
          farm_id: mc.farm_id || mc.FarmID || (mc.animal?.farm_id || mc.Animal?.FarmID),
          collection_date: mc.collection_date || mc.date || mc.Date,
          quantity: mc.quantity || mc.liters || mc.Liters,
          quality: mc.quality || mc.Quality,
          notes: mc.notes || mc.Notes || mc.observations || mc.Observations,
          created_at: mc.created_at || mc.CreatedAt,
          updated_at: mc.updated_at || mc.UpdatedAt,
          animal: mc.animal || mc.Animal
        }));
      }
      
      return [];
    } catch (error) {
      console.error('Erro ao buscar ordenhas do animal:', error);
      return [];
    }
  }
};

