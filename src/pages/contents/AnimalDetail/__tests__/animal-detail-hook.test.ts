import { renderHook, act } from '@testing-library/react';
import { useAnimalDetail } from '../hooks/useAnimalDetail';

const mockAnimal = {
  id: 1,
  farm_id: 1,
  animal_name: 'Vaca Teste',
  ear_tag_number_local: 123,
  ear_tag_number_register: 456,
  type: 'Bovino',
  sex: 0,
  breed: 'Holandesa',
  birth_date: '2020-01-01',
  photo: '',
  animal_type: 0,
  status: 0,
  confinement: false,
  fertilization: false,
  castrated: false,
  purpose: 1,
  current_batch: 1,
  father_id: undefined,
  mother_id: undefined,
  createdAt: '2023-01-01T00:00:00Z',
  updatedAt: '2023-01-01T00:00:00Z'
};

describe('useAnimalDetail', () => {
  it('should initialize with loading state', () => {
    const { result } = renderHook(() => useAnimalDetail(1));

    expect(result.current.loading).toBe(true);
    expect(result.current.animal).toBeUndefined();
    expect(result.current.error).toBeUndefined();
  });

  it('should load animal data successfully', async () => {
    const { result } = renderHook(() => useAnimalDetail(1));

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 100));
    });

    expect(result.current.loading).toBe(false);
    expect(result.current.animal).toEqual(mockAnimal);
    expect(result.current.error).toBeUndefined();
  });

  it('should handle update animal', async () => {
    const { result } = renderHook(() => useAnimalDetail(1));

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 100));
    });

    const updatedData = {
      ...mockAnimal,
      animal_name: 'Vaca Teste Editada'
    };

    await act(async () => {
      await result.current.updateAnimal(updatedData);
    });

    expect(result.current.animal?.animal_name).toBe('Vaca Teste Editada');
  });

  it('should handle error states', async () => {
    const { result } = renderHook(() => useAnimalDetail(999));

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 100));
    });

    expect(result.current.loading).toBe(false);
    expect(result.current.animal).toBeUndefined();
    expect(result.current.error).toBeDefined();
  });

  it('should toggle edit mode', () => {
    const { result } = renderHook(() => useAnimalDetail(1));

    expect(result.current.isEditing).toBe(false);

    act(() => {
      result.current.setIsEditing(true);
    });

    expect(result.current.isEditing).toBe(true);

    act(() => {
      result.current.setIsEditing(false);
    });

    expect(result.current.isEditing).toBe(false);
  });
});
