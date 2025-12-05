import { api } from '../../config/api';
import { Reproduction } from '../../types/reproduction';
import i18n from '../../locale/i18n';
import { ReproductionPhase } from '../../pages/contents/AnimalTable/domain/model/reproduction';

interface ReproductionApiData {
  id?: number;
  ID?: number;
  animal_id?: number;
  AnimalID?: number;
  farm_id?: number;
  FarmID?: number;
  date?: string;
  Date?: string;
  pregnancy_date?: string;
  insemination_date?: string;
  phase?: string;
  Phase?: string;
  current_phase?: number;
  notes?: string;
  Notes?: string;
  observations?: string;
  Observations?: string;
  created_at?: string;
  createdAt?: string;
  CreatedAt?: string;
  updated_at?: string;
  updatedAt?: string;
  UpdatedAt?: string;
  animal?: unknown;
  Animal?: unknown;
  animal_name?: string;
  AnimalName?: string;
  ear_tag?: number;
  EarTag?: number;
}

interface ReproductionResponse {
  success: boolean;
  data: ReproductionApiData | ReproductionApiData[];
  message?: string;
}

const getPhaseTranslation = (phaseNumber: number): string => {
  const phaseMap: { [key: number]: string } = {
    [ReproductionPhase.LACTACAO]: 'animalTable.reproduction.phases.lactacao',
    [ReproductionPhase.SECANDO]: 'animalTable.reproduction.phases.secando',
    [ReproductionPhase.VAZIAS]: 'animalTable.reproduction.phases.vazias',
    [ReproductionPhase.PRENHAS]: 'animalTable.reproduction.phases.prenhas'
  };

  const translationKey = phaseMap[phaseNumber];
  if (translationKey) {
    return i18n.t(translationKey);
  }

  return i18n.t('common.notInformed', { defaultValue: 'Não informado' });
};

export const reproductionService = {
  async getReproductionsByAnimal(animalId: number): Promise<Reproduction[]> {
    try {
      const response = await api.get<ReproductionResponse>(`/reproductions/animal`, {
        params: { animalId }
      });
      
      if (response.data.success && response.data.data) {
        const data = response.data.data;
        
        const reproductions = Array.isArray(data) ? data : [data];
        
        return reproductions.map((rep: ReproductionApiData) => {
          const date = rep.pregnancy_date || rep.insemination_date || rep.date || rep.Date;
          
          const phaseNumber = rep.current_phase !== undefined ? rep.current_phase : (rep.phase || 0);
          const phase = typeof phaseNumber === 'number' 
            ? getPhaseTranslation(phaseNumber)
            : (rep.phase || rep.Phase || i18n.t('common.notInformed', { defaultValue: 'Não informado' }));
          
          const animal = rep.animal || rep.Animal;
          const animalData = rep.animal_name 
            ? { 
                id: rep.animal_id || rep.AnimalID || 0,
                animal_name: rep.animal_name || rep.AnimalName || '',
                ear_tag_number_local: rep.ear_tag || rep.EarTag || 0
              }
            : (animal && typeof animal === 'object' && 'id' in animal
                ? {
                    id: (animal as { id?: number }).id || 0,
                    animal_name: (animal as { animal_name?: string }).animal_name || '',
                    ear_tag_number_local: (animal as { ear_tag_number_local?: number }).ear_tag_number_local || 0
                  }
                : undefined);

          return {
            id: rep.id || rep.ID || 0,
            animal_id: rep.animal_id || rep.AnimalID || 0,
            farm_id: rep.farm_id || rep.FarmID || 0,
            date: date || '',
            phase: phase,
            notes: rep.notes || rep.Notes || rep.observations || rep.Observations,
            created_at: rep.created_at || rep.createdAt || rep.CreatedAt || '',
            updated_at: rep.updated_at || rep.updatedAt || rep.UpdatedAt || '',
            animal: animalData
          };
        });
      }
      
      return [];
    } catch (error: unknown) {
      const axiosError = error as { response?: { status?: number } };
      if (axiosError.response?.status === 404) {
        return [];
      }
      if (axiosError.response?.status !== 404) {
        console.error('Error fetching reproductions for animal:', error);
      }
      return [];
    }
  }
};
