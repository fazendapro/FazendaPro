import React from 'react';
import { Table, Spin, Alert } from 'antd';
import { useAnimals } from '../hooks/useAnimals';
import { useFarm } from '../../../../hooks/useFarm';

const AnimalTable: React.FC = () => {
  const { farm } = useFarm();
  const { animals, loading, error } = useAnimals(farm.id);

  const columns = [
    { 
      title: 'Nome do Animal', 
      dataIndex: 'AnimalName', 
      key: 'AnimalName' 
    },
    { 
      title: 'Tipo', 
      dataIndex: 'Type', 
      key: 'Type' 
    },
    { 
      title: 'Número da Orelha (Local)', 
      dataIndex: 'EarTagNumberLocal', 
      key: 'EarTagNumberLocal' 
    },
    { 
      title: 'Número da Orelha (Registro)', 
      dataIndex: 'EarTagNumberRegister', 
      key: 'EarTagNumberRegister' 
    },
    { 
      title: 'Raça', 
      dataIndex: 'Breed', 
      key: 'Breed' 
    },
    { 
      title: 'Data de Nascimento', 
      dataIndex: 'BirthDate', 
      key: 'BirthDate',
      render: (date: string) => {
        if (!date) return '-';
        return new Date(date).toLocaleDateString('pt-BR');
      }
    },
    { 
      title: 'Sexo', 
      dataIndex: 'Sex', 
      key: 'Sex',
      render: (sex: number) => {
        return sex === 0 ? 'Macho' : 'Fêmea';
      }
    },
    { 
      title: 'Castrado', 
      dataIndex: 'Castrated', 
      key: 'Castrated',
      render: (castrated: boolean) => {
        return castrated ? 'Sim' : 'Não';
      }
    },
    { 
      title: 'Confinamento', 
      dataIndex: 'Confinement', 
      key: 'Confinement',
      render: (confinement: boolean) => {
        return confinement ? 'Sim' : 'Não';
      }
    },
    { 
      title: 'Fertilização', 
      dataIndex: 'Fertilization', 
      key: 'Fertilization',
      render: (fertilization: boolean) => {
        return fertilization ? 'Sim' : 'Não';
      }
    },
    { 
      title: 'Status', 
      dataIndex: 'Status', 
      key: 'Status',
      render: (status: number) => {
        const statusMap = {
          0: 'Ativo',
          1: 'Inativo',
          2: 'Vendido',
          3: 'Abatido'
        };
        return statusMap[status as keyof typeof statusMap] || 'Desconhecido';
      }
    },
    { 
      title: 'Propósito', 
      dataIndex: 'Purpose', 
      key: 'Purpose',
      render: (purpose: number) => {
        const purposeMap = {
          0: 'Leite',
          1: 'Carne',
          2: 'Reprodução',
          3: 'Misto'
        };
        return purposeMap[purpose as keyof typeof purposeMap] || 'Desconhecido';
      }
    },
    { 
      title: 'Lote Atual', 
      dataIndex: 'CurrentBatch', 
      key: 'CurrentBatch' 
    },
    { 
      title: 'Criado em', 
      dataIndex: 'CreatedAt', 
      key: 'CreatedAt',
      render: (date: string) => {
        if (!date) return '-';
        return new Date(date).toLocaleDateString('pt-BR');
      }
    }
  ];

  console.log('Animals data:', animals);
  console.log('Animals type:', typeof animals);
  console.log('Is array:', Array.isArray(animals));

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

  // Garantir que animals seja sempre um array
  const animalsData = Array.isArray(animals) ? animals : [];

  return (
    <Table 
      columns={columns} 
      dataSource={animalsData} 
      pagination={{ showSizeChanger: true }}
      rowKey="ID"
      scroll={{ x: 1500 }}
    />
  );
};

export { AnimalTable };