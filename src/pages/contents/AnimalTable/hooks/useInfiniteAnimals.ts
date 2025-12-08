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
  const uniqueAnimalsCountRef = useRef(0);

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

      console.log('[useInfiniteAnimals] 📥 Requisição recebida:', {
        page: pageToLoad,
        resetList,
        animaisRecebidos: animalsData.length,
        totalEsperado: totalCount,
        idsRecebidos: animalsData.map(a => a.id)
      });

      if (animalsData.length === 0 && !resetList && pageToLoad > 1) {
        console.warn('[useInfiniteAnimals] ⚠️ Página vazia recebida, definindo hasMore=false');
        setHasMore(false);
        setLoading(false);
        loadingRef.current = false;
        return;
      }

      setTotal(totalCount);

      const normalizeId = (id: string | number): string => {
        return String(id);
      };

      if (resetList || pageToLoad === 1) {
        const uniqueAnimals = Array.from(
          new Map(animalsData.map(animal => [normalizeId(animal.id), animal])).values()
        );
        uniqueAnimalsCountRef.current = uniqueAnimals.length;
        
        console.log('[useInfiniteAnimals] 🔄 Página 1/Reset:', {
          animaisRecebidos: animalsData.length,
          animaisUnicos: uniqueAnimals.length,
          duplicatasRemovidas: animalsData.length - uniqueAnimals.length,
          idsUnicos: uniqueAnimals.map(a => a.id),
          hasMore: uniqueAnimals.length < totalCount
        });
        
        setAnimals(uniqueAnimals);
        setHasMore(uniqueAnimals.length < totalCount);
      } else {
        let finalAnimalsCount = 0;
        
        setAnimals(prev => {
          const animalsMap = new Map<string, Animal>();
          const prevIds = prev.map(a => normalizeId(a.id));

          prev.forEach(animal => {
            animalsMap.set(normalizeId(animal.id), animal);
          });

          const novosIds: string[] = [];
          animalsData.forEach(animal => {
            const normalizedId = normalizeId(animal.id);
            if (!animalsMap.has(normalizedId)) {
              novosIds.push(normalizedId);
            }
            animalsMap.set(normalizedId, animal);
          });

          const finalAnimals = Array.from(animalsMap.values());
          finalAnimalsCount = finalAnimals.length;
          uniqueAnimalsCountRef.current = finalAnimals.length;
          
          console.log('[useInfiniteAnimals] ➕ Acumulando página:', {
            page: pageToLoad,
            animaisAnteriores: prev.length,
            animaisRecebidos: animalsData.length,
            animaisNovos: novosIds.length,
            animaisUnicosFinais: finalAnimals.length,
            idsNovos: novosIds,
            idsAnteriores: prevIds.slice(0, 5),
            hasMore: finalAnimals.length < totalCount
          });
          
          return finalAnimals;
        });
        
        // Atualizar hasMore fora do setState para evitar problemas de sincronização
        setHasMore(finalAnimalsCount < totalCount);
      }
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
      console.log('[useInfiniteAnimals] 🔄 SearchTerm mudou, resetando:', {
        anterior: previousSearchTerm.current,
        novo: searchTerm
      });
      previousSearchTerm.current = searchTerm;
      setAnimals([]);
      setPage(1);
      setHasMore(true);
      setTotal(0);
      uniqueAnimalsCountRef.current = 0;
      loadAnimals(1, true);
    }
  }, [searchTerm, loadAnimals]);

  const loadMore = useCallback(() => {
    console.log('[useInfiniteAnimals] 🔽 loadMore chamado:', {
      loading,
      hasMore,
      pageAtual: page,
      proximaPage: page + 1,
      animaisAtuais: uniqueAnimalsCountRef.current,
      total
    });
    
    if (!loading && hasMore) {
      const nextPage = page + 1;
      
      // Verificar se a próxima página faz sentido
      const expectedMaxPage = Math.ceil(total / pageSize);
      if (nextPage > expectedMaxPage) {
        console.warn('[useInfiniteAnimals] ⚠️ Tentando carregar página além do esperado:', {
          nextPage,
          expectedMaxPage,
          total,
          pageSize
        });
        setHasMore(false);
        return;
      }
      
      setPage(nextPage);
      loadAnimals(nextPage);
    } else {
      console.log('[useInfiniteAnimals] ⛔ loadMore bloqueado:', {
        loading,
        hasMore
      });
    }
  }, [loading, hasMore, page, loadAnimals, total, pageSize]);

  const reset = useCallback(() => {
    setAnimals([]);
    setPage(1);
    setHasMore(true);
    setTotal(0);
    setError(null);
    loadAnimals(1);
  }, [loadAnimals]);

  useEffect(() => {
    console.log('[useInfiniteAnimals] 📊 Estado atualizado:', {
      animais: animals.length,
      loading,
      hasMore,
      total,
      page,
      faltam: total - animals.length,
      ids: animals.slice(0, 10).map(a => a.id)
    });
    
    if (animals.length < total && hasMore && !loading) {
      console.warn('[useInfiniteAnimals] ⚠️ ATENÇÃO: Temos menos animais que o total e hasMore=true, mas loadMore não foi chamado!');
    }
  }, [animals, loading, hasMore, total, page]);

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

