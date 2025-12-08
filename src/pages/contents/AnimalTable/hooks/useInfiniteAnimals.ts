import { useState, useCallback, useEffect, useRef } from 'react';
import { useFarm } from '../../../../hooks/useFarm';
import { GetAnimalsByFarmFactory } from '../factories/usecases/get-animals-by-farm-factory';
import { Animal } from '../types/type';

interface UseInfiniteAnimalsOptions {
  pageSize?: number;
  searchTerm?: string;
  farmId?: number;
}

export const useInfiniteAnimals = (options: UseInfiniteAnimalsOptions = {}) => {
  const { pageSize = 50, searchTerm = '', farmId } = options;
  const { farm } = useFarm();
  
  const [animals, setAnimals] = useState<Animal[]>([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState(0);
  
  const previousSearchTerm = useRef(searchTerm);
  const isInitialMount = useRef(true);
  const loadingRef = useRef(false);

  const loadAnimals = useCallback(async (pageToLoad: number, resetList = false) => {
    const targetFarmId = farmId || farm?.id;
    if (!targetFarmId) {
      setError('Fazenda não encontrada');
      return;
    }

    if (loadingRef.current) return;

    loadingRef.current = true;
    setLoading(true);
    setError(null);

    try {
      const getAnimalsByFarmUseCase = GetAnimalsByFarmFactory();
      const result = await getAnimalsByFarmUseCase.getAnimalsByFarm({
        farm_id: targetFarmId,
        page: pageToLoad,
        limit: pageSize,
        searchTerm: searchTerm.trim() || undefined
      });

      const paginatedData = result.data;
      const animalsData = paginatedData?.animals || [];
      const totalCount = paginatedData?.total || 0;

      setTotal(totalCount);

      if (resetList || pageToLoad === 1) {
        const uniqueAnimals = Array.from(
          new Map(animalsData.map(animal => [animal.id, animal])).values()
        );
        setAnimals(uniqueAnimals);
      } else {
        setAnimals(prev => {
          const animalsMap = new Map<string, Animal>();

          prev.forEach(animal => {
            animalsMap.set(animal.id, animal);
          });

          animalsData.forEach(animal => {
            animalsMap.set(animal.id, animal);
          });

          return Array.from(animalsMap.values());
        });
      }

      const loadedCount = pageToLoad * pageSize;
      setHasMore(loadedCount < totalCount);
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Erro ao buscar animais';
      setError(errorMessage);
    } finally {
      setLoading(false);
      loadingRef.current = false;
    }
  }, [farmId, farm?.id, pageSize, searchTerm]);

  useEffect(() => {
    if (isInitialMount.current) {
      isInitialMount.current = false;
      loadAnimals(1, true);
    }
  }, [loadAnimals]);

  useEffect(() => {
    if (previousSearchTerm.current !== searchTerm) {
      previousSearchTerm.current = searchTerm;
      setAnimals([]);
      setPage(1);
      setHasMore(true);
      setTotal(0);
      loadAnimals(1, true);
    }
  }, [searchTerm, loadAnimals]);

  const loadMore = useCallback(() => {
    if (!loading && hasMore) {
      const nextPage = page + 1;
      setPage(nextPage);
      loadAnimals(nextPage);
    }
  }, [loading, hasMore, page, loadAnimals]);

  const reset = useCallback(() => {
    setAnimals([]);
    setPage(1);
    setHasMore(true);
    setTotal(0);
    setError(null);
    loadAnimals(1);
  }, [loadAnimals]);

  return {
    animals,
    loading,
    error,
    hasMore,
    total,
    loadMore,
    reset
  };
};

