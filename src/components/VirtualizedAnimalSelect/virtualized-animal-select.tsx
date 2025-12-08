import React, { useState, useCallback, useRef, useEffect, useMemo } from 'react';
import { Select, Spin } from 'antd';
import VirtualList from 'rc-virtual-list';
import { useInfiniteAnimals } from '../../pages/contents/AnimalTable/hooks/useInfiniteAnimals';
import { Animal } from '../../pages/contents/AnimalTable/types/type';

interface VirtualizedAnimalSelectProps {
  farmId?: number;
  onAnimalLoad?: (animals: Animal[]) => void;
  itemHeight?: number;
  dropdownHeight?: number;
  pageSize?: number;
  placeholder?: string;
  loading?: boolean;
  disabled?: boolean;
  allowClear?: boolean;
  value?: number;
  onChange?: (value: number) => void;
  filterAnimals?: (animal: Animal) => boolean;
  getPopupContainer?: (triggerNode: HTMLElement) => HTMLElement;
}

export const VirtualizedAnimalSelect: React.FC<VirtualizedAnimalSelectProps> = ({
  farmId,
  onAnimalLoad,
  itemHeight = 32,
  dropdownHeight = 256,
  pageSize = 20,
  placeholder = 'Selecione um animal',
  loading: externalLoading,
  disabled,
  allowClear,
  value,
  onChange,
  filterAnimals,
  getPopupContainer
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [open, setOpen] = useState(false);
  const debounceTimer = useRef<NodeJS.Timeout | null>(null);

  const { animals: allAnimals, loading, hasMore, loadMore } = useInfiniteAnimals({
    farmId,
    searchTerm,
    pageSize
  });

  const animals = useMemo(() => {
    console.log('[VirtualizedAnimalSelect] 🎯 Processando animais:', {
      animaisDoHook: allAnimals.length,
      temFiltro: !!filterAnimals,
      idsDoHook: allAnimals.slice(0, 10).map(a => a.id)
    });

    const filtered = filterAnimals ? allAnimals.filter(filterAnimals) : allAnimals;
    
    console.log('[VirtualizedAnimalSelect] 🔍 Após filtro:', {
      animaisFiltrados: filtered.length,
      animaisRemovidos: allAnimals.length - filtered.length,
      idsFiltrados: filtered.slice(0, 10).map(a => a.id)
    });

    const normalizeId = (id: string | number): string => String(id);
    const uniqueAnimalsMap = new Map<string, Animal>();
    filtered.forEach(animal => {
      uniqueAnimalsMap.set(normalizeId(animal.id), animal);
    });
    
    const finalAnimals = Array.from(uniqueAnimalsMap.values());
    
    console.log('[VirtualizedAnimalSelect] ✅ Animais finais:', {
      animaisFinais: finalAnimals.length,
      duplicatasRemovidas: filtered.length - finalAnimals.length,
      idsFinais: finalAnimals.slice(0, 10).map(a => a.id),
      totalIds: finalAnimals.map(a => a.id)
    });
    
    return finalAnimals;
  }, [allAnimals, filterAnimals]);

  useEffect(() => {
    if (onAnimalLoad && animals.length > 0) {
      onAnimalLoad(animals);
    }
  }, [animals, onAnimalLoad]);

  const handleSearch = useCallback((value: string) => {
    if (debounceTimer.current) {
      clearTimeout(debounceTimer.current);
    }

    debounceTimer.current = setTimeout(() => {
      setSearchTerm(value);
    }, 300);
  }, []);

  const handleScroll = useCallback((e: React.UIEvent<HTMLElement>) => {
    const target = e.currentTarget;
    const scrollTop = target.scrollTop;
    const scrollHeight = target.scrollHeight;
    const clientHeight = target.clientHeight;

    const threshold = itemHeight * 3;
    const distanceFromBottom = scrollHeight - scrollTop - clientHeight;

    if (distanceFromBottom < threshold && hasMore && !loading) {
      console.log('[VirtualizedAnimalSelect] 📜 Scroll próximo ao fim, carregando mais:', {
        scrollTop,
        scrollHeight,
        clientHeight,
        distanceFromBottom,
        threshold,
        hasMore,
        loading
      });
      loadMore();
    }
  }, [hasMore, loading, loadMore, itemHeight]);

  const handleOpenChange = useCallback((visible: boolean) => {
    console.log('[VirtualizedAnimalSelect] 🔄 Dropdown', visible ? 'aberto' : 'fechado');
    setOpen(visible);
    if (!visible) {
      setSearchTerm('');
      if (debounceTimer.current) {
        clearTimeout(debounceTimer.current);
      }
    }
  }, []);

  const dropdownRender = useCallback(() => {
    console.log('[VirtualizedAnimalSelect] 🎨 Renderizando dropdown:', {
      animaisParaRenderizar: animals.length,
      loading,
      hasMore,
      idsParaRenderizar: animals.map(a => a.id)
    });

    if (animals.length === 0 && !loading) {
      return <div style={{ padding: '8px 12px', textAlign: 'center', color: '#999' }}>Nenhum animal encontrado</div>;
    }

    return (
      <div>
        <VirtualList
          data={animals}
          height={dropdownHeight}
          itemHeight={itemHeight}
          itemKey={(animal: Animal) => String(animal.id)}
          onScroll={handleScroll}
        >
          {(animal: Animal) => {
            const animalId = typeof animal.id === 'string' ? parseInt(animal.id) : animal.id;
            return (
              <div
                key={animal.id}
                className="ant-select-item ant-select-item-option"
                style={{
                  padding: '5px 12px',
                  cursor: 'pointer',
                  height: itemHeight,
                  display: 'flex',
                  alignItems: 'center'
                }}
                onClick={() => {
                  if (onChange) {
                    onChange(animalId);
                  }
                  setOpen(false);
                }}
              >
                {animal.animal_name} - {animal.ear_tag_number_local}
              </div>
            );
          }}
        </VirtualList>
        {loading && (
          <div style={{ padding: '8px', textAlign: 'center', borderTop: '1px solid #f0f0f0' }}>
            <Spin size="small" />
          </div>
        )}
      </div>
    );
  }, [animals, dropdownHeight, itemHeight, handleScroll, loading, onChange]);

  const selectOptions = useMemo(() => {
    const optionsMap = new Map<number, { label: string; value: number }>();
    animals.forEach(animal => {
      const animalId = typeof animal.id === 'string' ? parseInt(animal.id) : animal.id;
      if (!optionsMap.has(animalId)) {
        optionsMap.set(animalId, {
          label: `${animal.animal_name} - ${animal.ear_tag_number_local}`,
          value: animalId
        });
      }
    });
    return Array.from(optionsMap.values());
  }, [animals]);

  return (
    <Select
      open={open}
      onOpenChange={handleOpenChange}
      placeholder={placeholder}
      showSearch
      filterOption={false}
      onSearch={handleSearch}
      notFoundContent={loading ? <Spin size="small" /> : 'Nenhum animal encontrado'}
      dropdownRender={dropdownRender}
      loading={externalLoading || loading}
      disabled={disabled}
      allowClear={allowClear}
      value={value}
      onChange={onChange}
      getPopupContainer={getPopupContainer}
      options={selectOptions}
    />
  );
};
