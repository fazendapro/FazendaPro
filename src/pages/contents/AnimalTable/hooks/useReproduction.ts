import { useState, useCallback } from 'react';
import { useFarm } from '../../../../hooks/useFarm';
import { remoteCreateReproduction } from '../data/usecases/remote-create-reproduction';
import { remoteGetReproductionsByFarm } from '../data/usecases/remote-get-reproductions-by-farm';
import { remoteUpdateReproductionPhase } from '../data/usecases/remote-update-reproduction-phase';
import { remoteUpdateReproduction } from '../data/usecases/remote-update-reproduction';
import { remoteDeleteReproduction } from '../data/usecases/remote-delete-reproduction';
import { CreateReproductionRequest, UpdateReproductionPhaseRequest, Reproduction } from '../domain/model/reproduction';

export const useReproduction = () => {
  const { farm } = useFarm();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createReproduction = useCallback(async (data: CreateReproductionRequest): Promise<Reproduction | null> => {
    console.log('useReproduction - createReproduction called with:', data);
    console.log('useReproduction - farm:', farm);
    
    if (!farm?.id) {
      console.error('useReproduction - Farm not found');
      setError('Fazenda não encontrada');
      return null;
    }

    setLoading(true);
    setError(null);

    try {
      console.log('useReproduction - Calling remoteCreateReproduction...');
      const result = await remoteCreateReproduction(data);
      console.log('useReproduction - Result received:', result);
      return result;
    } catch (err: any) {
      console.error('useReproduction - Error:', err);
      const errorMessage = err.response?.data?.message || 'Erro ao criar registro de reprodução';
      setError(errorMessage);
      return null;
    } finally {
      setLoading(false);
    }
  }, [farm?.id]);

  const getReproductionsByFarm = useCallback(async (): Promise<Reproduction[]> => {
    if (!farm?.id) {
      setError('Fazenda não encontrada');
      return [];
    }

    setLoading(true);
    setError(null);

    try {
      const result = await remoteGetReproductionsByFarm(farm.id);
      return result;
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Erro ao buscar registros de reprodução';
      setError(errorMessage);
      return [];
    } finally {
      setLoading(false);
    }
  }, [farm?.id]);

  const updateReproductionPhase = useCallback(async (data: UpdateReproductionPhaseRequest): Promise<boolean> => {
    setLoading(true);
    setError(null);

    try {
      await remoteUpdateReproductionPhase(data);
      return true;
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Erro ao atualizar fase de reprodução';
      setError(errorMessage);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const updateReproduction = useCallback(async (data: CreateReproductionRequest): Promise<boolean> => {
    setLoading(true);
    setError(null);

    try {
      await remoteUpdateReproduction(data);
      return true;
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Erro ao atualizar registro de reprodução';
      setError(errorMessage);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const deleteReproduction = useCallback(async (id: number): Promise<boolean> => {
    setLoading(true);
    setError(null);

    try {
      await remoteDeleteReproduction(id);
      return true;
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Erro ao deletar registro de reprodução';
      setError(errorMessage);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    createReproduction,
    getReproductionsByFarm,
    updateReproductionPhase,
    updateReproduction,
    deleteReproduction,
    loading,
    error,
    clearError: () => setError(null)
  };
};
