import { useState, useEffect, useCallback } from 'react';
import { animalDetailService } from '../data';
import { AnimalDetail, AnimalDetailFormData, AnimalParent } from '../types';

interface UseAnimalDetailReturn {
  animal: AnimalDetail | undefined;
  loading: boolean;
  error: string | undefined;
  isEditing: boolean;
  setIsEditing: (editing: boolean) => void;
  updateAnimal: (data: AnimalDetailFormData) => Promise<void>;
  uploadPhoto: (file: File) => Promise<string | null>;
  fathers: AnimalParent[];
  mothers: AnimalParent[];
  loadingParents: boolean;
  refreshAnimal: () => Promise<void>;
}

export const useAnimalDetail = (animalId: number): UseAnimalDetailReturn => {
  const [animal, setAnimal] = useState<AnimalDetail | undefined>();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | undefined>();
  const [isEditing, setIsEditing] = useState(false);
  const [fathers, setFathers] = useState<AnimalParent[]>([]);
  const [mothers, setMothers] = useState<AnimalParent[]>([]);
  const [loadingParents, setLoadingParents] = useState(false);

  const loadAnimal = useCallback(async () => {
    try {
      setLoading(true);
      setError(undefined);
      
      const response = await animalDetailService.getAnimalById(animalId);
      
      if (response.success && response.data) {
        setAnimal(response.data);
      } else {
        setError(response.message || 'Erro ao carregar animal');
      }
    } catch {
      setError('Erro ao carregar animal');
    } finally {
      setLoading(false);
    }
  }, [animalId]);

  const loadParents = useCallback(async (farmId: number) => {
    try {
      setLoadingParents(true);
      
      const [fathersResponse, mothersResponse] = await Promise.all([
        animalDetailService.getAnimalsBySex(farmId, 1),
        animalDetailService.getAnimalsBySex(farmId, 0)
      ]);

      if (fathersResponse.success && fathersResponse.data) {
        setFathers(fathersResponse.data);
      }

      if (mothersResponse.success && mothersResponse.data) {
        setMothers(mothersResponse.data);
      }
    } catch {
      console.error('Erro ao carregar pais');
    } finally {
      setLoadingParents(false);
    }
  }, []);

  const refreshAnimal = useCallback(async () => {
    await loadAnimal();
  }, [loadAnimal]);

  const updateAnimal = useCallback(async (data: AnimalDetailFormData) => {
    if (!animal) return;

    try {
      const updateData = {
        id: animal.id,
        farm_id: animal.farm_id,
        ...data
      };

      const response = await animalDetailService.updateAnimal(updateData);
      
      if (response.success && response.data) {
        setAnimal(response.data);
        setIsEditing(false);
        await refreshAnimal();
      } else {
        console.error(response.message || 'Erro ao atualizar animal');
      }
    } catch (error) {
      console.error('Erro ao atualizar animal:', error);
    }
  }, [animal, refreshAnimal]);

  const uploadPhoto = useCallback(async (file: File): Promise<string | null> => {
    if (!animal) return null;

    try {
      const response = await animalDetailService.uploadAnimalPhoto(animal.id, file);
      
      if (response.success && response.data) {
        setAnimal(response.data);
        return response.data.photo || null;
      } else {
        console.error(response.message || 'Erro ao fazer upload da foto');
        return null;
      }
    } catch (error) {
      console.error('Erro ao fazer upload da foto:', error);
      return null;
    }
  }, [animal]);

  useEffect(() => {
    loadAnimal();
  }, [loadAnimal]);

  useEffect(() => {
    if (animal?.farm_id) {
      loadParents(animal.farm_id);
    }
  }, [animal?.farm_id, loadParents]);

  return {
    animal,
    loading,
    error,
    isEditing,
    setIsEditing,
    updateAnimal,
    uploadPhoto,
    fathers,
    mothers,
    loadingParents,
    refreshAnimal
  };
};
