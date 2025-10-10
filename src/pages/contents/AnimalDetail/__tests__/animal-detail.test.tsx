import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { AnimalDetailProvider } from '../hooks/AnimalDetailProvider';
import { AnimalDetail } from '../presentation';

// const mockAnimal = {
//   id: 1,
//   farm_id: 1,
//   animal_name: 'Vaca Teste',
//   ear_tag_number_local: 123,
//   ear_tag_number_register: 456,
//   type: 'Bovino',
//   sex: 0,
//   breed: 'Holandesa',
//   birth_date: '2020-01-01',
//   photo: '',
//   animal_type: 0,
//   status: 0,
//   confinement: false,
//   fertilization: false,
//   castrated: false,
//   purpose: 1,
//   current_batch: 1,
//   father_id: undefined,
//   mother_id: undefined,
//   createdAt: '2023-01-01T00:00:00Z',
//   updatedAt: '2023-01-01T00:00:00Z'
// };

const MockRouter = ({ children }: { children: React.ReactNode }) => (
  <BrowserRouter>{children}</BrowserRouter>
);

describe('AnimalDetail', () => {
  it('should render animal details correctly', async () => {
    render(
      <MockRouter>
        <AnimalDetailProvider animalId={1}>
          <AnimalDetail />
        </AnimalDetailProvider>
      </MockRouter>
    );

    await waitFor(() => {
      expect(screen.getByDisplayValue('Vaca Teste')).toBeInTheDocument();
      expect(screen.getByDisplayValue('123')).toBeInTheDocument();
      expect(screen.getByDisplayValue('456')).toBeInTheDocument();
      expect(screen.getByDisplayValue('Bovino')).toBeInTheDocument();
      expect(screen.getByDisplayValue('Holandesa')).toBeInTheDocument();
    });
  });

  it('should allow editing animal details', async () => {
    render(
      <MockRouter>
        <AnimalDetailProvider animalId={1}>
          <AnimalDetail />
        </AnimalDetailProvider>
      </MockRouter>
    );

    await waitFor(() => {
      const editButton = screen.getByText('Editar');
      fireEvent.click(editButton);
    });

    const nameInput = screen.getByDisplayValue('Vaca Teste');
    fireEvent.change(nameInput, { target: { value: 'Vaca Teste Editada' } });

    expect(nameInput).toHaveValue('Vaca Teste Editada');
  });

  it('should save changes when save button is clicked', async () => {
    render(
      <MockRouter>
        <AnimalDetailProvider animalId={1}>
          <AnimalDetail />
        </AnimalDetailProvider>
      </MockRouter>
    );

    await waitFor(() => {
      const editButton = screen.getByText('Editar');
      fireEvent.click(editButton);
    });

    const nameInput = screen.getByDisplayValue('Vaca Teste');
    fireEvent.change(nameInput, { target: { value: 'Vaca Teste Editada' } });

    const saveButton = screen.getByText('Salvar');
    fireEvent.click(saveButton);

    await waitFor(() => {
      expect(screen.getByText('Animal atualizado com sucesso')).toBeInTheDocument();
    });
  });

  it('should cancel editing when cancel button is clicked', async () => {
    render(
      <MockRouter>
        <AnimalDetailProvider animalId={1}>
          <AnimalDetail />
        </AnimalDetailProvider>
      </MockRouter>
    );

    await waitFor(() => {
      const editButton = screen.getByText('Editar');
      fireEvent.click(editButton);
    });

    const nameInput = screen.getByDisplayValue('Vaca Teste');
    fireEvent.change(nameInput, { target: { value: 'Vaca Teste Editada' } });

    const cancelButton = screen.getByText('Cancelar');
    fireEvent.click(cancelButton);

    await waitFor(() => {
      expect(screen.getByDisplayValue('Vaca Teste')).toBeInTheDocument();
    });
  });

  it('should display parent selection dropdowns', async () => {
    render(
      <MockRouter>
        <AnimalDetailProvider animalId={1}>
          <AnimalDetail />
        </AnimalDetailProvider>
      </MockRouter>
    );

    await waitFor(() => {
      const editButton = screen.getByText('Editar');
      fireEvent.click(editButton);
    });

    expect(screen.getByText('Pai')).toBeInTheDocument();
    expect(screen.getByText('Mãe')).toBeInTheDocument();
  });

  it('should handle photo upload', async () => {
    render(
      <MockRouter>
        <AnimalDetailProvider animalId={1}>
          <AnimalDetail />
        </AnimalDetailProvider>
      </MockRouter>
    );

    await waitFor(() => {
      const editButton = screen.getByText('Editar');
      fireEvent.click(editButton);
    });

    const fileInput = screen.getByLabelText('Foto do Animal');
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
    
    fireEvent.change(fileInput, { target: { files: [file] } });

    await waitFor(() => {
      expect((fileInput as HTMLInputElement).files?.[0]).toBe(file);
    });
  });
});
