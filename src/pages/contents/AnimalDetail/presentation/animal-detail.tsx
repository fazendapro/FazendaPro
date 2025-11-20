import React from 'react';
import { Alert, Spin } from 'antd';
import { useAnimalDetailContext } from '../hooks';
import { AnimalDetailDisplay } from './animal-detail-display';
import { AnimalDetailForm } from './animal-detail-form';
import { AnimalDetailFormData } from '../types';

export const AnimalDetail: React.FC = () => {
  const { animal, loading, error, isEditing, setIsEditing, updateAnimal, uploadPhoto, refreshAnimal } = useAnimalDetailContext();

  const handleEdit = () => {
    setIsEditing(true);
  };

  const handleSave = async (data: AnimalDetailFormData, photoFile: File | null) => {
    try {
      await updateAnimal(data);
      if (photoFile) {
        await uploadPhoto(photoFile);
      }
      await refreshAnimal();
    } catch (error) {
      console.error('Erro ao salvar animal:', error);
    }
  };

  const handleCancel = () => {
    setIsEditing(false);
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (error) {
    return (
      <Alert
        message="Erro"
        description={error}
        type="error"
        showIcon
        style={{ marginBottom: '16px' }}
      />
    );
  }

  if (!animal) {
    return (
      <Alert
        message="Animal não encontrado"
        description="O animal solicitado não foi encontrado."
        type="warning"
        showIcon
      />
    );
  }

  return (
    <div style={{ padding: '24px' }}>
      {isEditing ? (
        <AnimalDetailForm 
          onSave={handleSave} 
          onCancel={handleCancel}
        />
      ) : (
        <AnimalDetailDisplay onEdit={handleEdit} />
      )}
    </div>
  );
};

export default AnimalDetail;
