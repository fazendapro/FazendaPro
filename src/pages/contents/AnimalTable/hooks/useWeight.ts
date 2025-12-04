import { useState, useCallback } from 'react';
import { useFarm } from '../../../../hooks/useFarm';
import { remoteCreateOrUpdateWeight } from '../data/usecases/remote-create-or-update-weight';
import { remoteGetWeightByAnimal } from '../data/usecases/remote-get-weight-by-animal';
import { remoteGetWeightsByFarm } from '../data/usecases/remote-get-weights-by-farm';
import { remoteUpdateWeight } from '../data/usecases/remote-update-weight';
import { CreateOrUpdateWeightRequest, UpdateWeightRequest, Weight } from '../domain/model/weight';

export const useWeight = () => {
  const { farm } = useFarm();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createOrUpdateWeight = useCallback(async (data: CreateOrUpdateWeightRequest): Promise<Weight | null> => {
    if (!farm?.id) {
      setError('Fazenda não encontrada');
      return null;
    }

    setLoading(true);
    setError(null);

    try {
      const result = await remoteCreateOrUpdateWeight(data);
      return result;
    } catch (err: unknown) {
      const errorMessage = (err as { response?: { data?: { message?: string } } })?.response?.data?.message || 'Erro ao criar ou atualizar registro de peso';
      setError(errorMessage);
      return null;
    } finally {
      setLoading(false);
    }
  }, [farm?.id]);

  const getWeightByAnimal = useCallback(async (animalId: number): Promise<Weight | null> => {
    setLoading(true);
    setError(null);

    try {
      const result = await remoteGetWeightByAnimal(animalId);
      return result;
    } catch (err: unknown) {
      const errorMessage = (err as { response?: { data?: { message?: string } } })?.response?.data?.message || 'Erro ao buscar peso do animal';
      setError(errorMessage);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const getWeightsByFarm = useCallback(async (): Promise<Weight[]> => {
    if (!farm?.id) {
      setError('Fazenda não encontrada');
      return [];
    }

    setLoading(true);
    setError(null);

    try {
      const result = await remoteGetWeightsByFarm(farm.id);
      return result;
    } catch (err: unknown) {
      const errorMessage = (err as { response?: { data?: { message?: string } } })?.response?.data?.message || 'Erro ao buscar pesos da fazenda';
      setError(errorMessage);
      return [];
    } finally {
      setLoading(false);
    }
  }, [farm?.id]);

  const updateWeight = useCallback(async (data: UpdateWeightRequest): Promise<boolean> => {
    setLoading(true);
    setError(null);

    try {
      await remoteUpdateWeight(data);
      return true;
    } catch (err: unknown) {
      const errorMessage = (err as { response?: { data?: { message?: string } } })?.response?.data?.message || 'Erro ao atualizar registro de peso';
      setError(errorMessage);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    createOrUpdateWeight,
    getWeightByAnimal,
    getWeightsByFarm,
    updateWeight,
    loading,
    error,
    clearError: () => setError(null)
  };
};

