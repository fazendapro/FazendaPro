import { useCallback, useEffect, useImperativeHandle, useRef, useState, forwardRef } from 'react';
import { Table, Button, Tag, Space, message, Popconfirm } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useReproduction } from '../../hooks/useReproduction';
import { useFarm } from '../../../../../hooks/useFarm';
import { Reproduction, ReproductionPhase, ReproductionPhaseLabels, ReproductionPhaseColors } from '../../domain/model/reproduction';
import { CreateReproductionModal } from './create-reproduction-modal';
import { UpdateReproductionPhaseModal } from './update-reproduction-phase-modal';

interface ReproductionTableRef {
  refetch: () => void;
}

interface ReproductionTableProps {
  onAddReproduction: () => void;
  onEditReproduction: (reproduction: Reproduction) => void;
}

export const ReproductionTable = forwardRef<ReproductionTableRef, ReproductionTableProps>(({ onAddReproduction, onEditReproduction }, ref) => {
  const { t } = useTranslation();
  const { farm } = useFarm();
  const { getReproductionsByFarm, updateReproduction, deleteReproduction, loading, error } = useReproduction();
  const [reproductions, setReproductions] = useState<Reproduction[]>([]);
  const [isCreateModalVisible, setIsCreateModalVisible] = useState(false);
  const [isUpdatePhaseModalVisible, setIsUpdatePhaseModalVisible] = useState(false);
  const [selectedReproduction, setSelectedReproduction] = useState<Reproduction | null>(null);

  const fetchReproductions = useCallback(async () => {
    if (!farm?.id) return;
    
    const data = await getReproductionsByFarm();
    setReproductions(data);
  }, [farm?.id, getReproductionsByFarm]);

  useImperativeHandle(ref, () => ({
    refetch: fetchReproductions
  }), [fetchReproductions]);

  useEffect(() => {
    fetchReproductions();
  }, [fetchReproductions]);

  useEffect(() => {
    if (error) {
      message.error(error);
    }
  }, [error]);

  const handleDeleteReproduction = async (id: number) => {
    const success = await deleteReproduction(id);
    if (success) {
      message.success(t('animalTable.reproduction.deletedSuccessfully'));
      fetchReproductions();
    } else {
      message.error(t('animalTable.reproduction.deleteError'));
    }
  };

  const handleUpdatePhase = (reproduction: Reproduction) => {
    setSelectedReproduction(reproduction);
    setIsUpdatePhaseModalVisible(true);
  };

  const handleCreateModalSuccess = () => {
    setIsCreateModalVisible(false);
    fetchReproductions();
  };

  const handleUpdatePhaseModalSuccess = () => {
    setIsUpdatePhaseModalVisible(false);
    setSelectedReproduction(null);
    fetchReproductions();
  };

  const columns = [
    {
      title: t('animalTable.reproduction.animalName'),
      dataIndex: 'animal_name',
      key: 'animal_name',
      sorter: (a: Reproduction, b: Reproduction) => a.animal_id - b.animal_id,
    },
    {
      title: t('animalTable.reproduction.earTag'),
      dataIndex: 'ear_tag',
      key: 'ear_tag',
    },
    {
      title: t('animalTable.reproduction.currentPhase'),
      dataIndex: 'current_phase',
      key: 'current_phase',
      render: (phase: number) => (
        <Tag color={ReproductionPhaseColors[phase as ReproductionPhase]}>
          {ReproductionPhaseLabels[phase as ReproductionPhase]}
        </Tag>
      ),
      filters: Object.entries(ReproductionPhaseLabels).map(([value, text]) => ({
        text,
        value: parseInt(value),
      })),
      onFilter: (value: any, record: Reproduction) => record.current_phase === value,
    },
    {
      title: t('animalTable.reproduction.inseminationDate'),
      dataIndex: 'insemination_date',
      key: 'insemination_date',
      render: (date: string) => date ? new Date(date).toLocaleDateString('pt-BR') : '-',
    },
    {
      title: t('animalTable.reproduction.pregnancyDate'),
      dataIndex: 'pregnancy_date',
      key: 'pregnancy_date',
      render: (date: string) => date ? new Date(date).toLocaleDateString('pt-BR') : '-',
    },
    {
      title: t('animalTable.reproduction.expectedBirthDate'),
      dataIndex: 'expected_birth_date',
      key: 'expected_birth_date',
      render: (date: string) => date ? new Date(date).toLocaleDateString('pt-BR') : '-',
    },
    {
      title: t('animalTable.reproduction.veterinaryConfirmation'),
      dataIndex: 'veterinary_confirmation',
      key: 'veterinary_confirmation',
      render: (confirmed: boolean) => (
        <Tag color={confirmed ? 'green' : 'red'}>
          {confirmed ? 'Confirmado' : 'Não confirmado'}
        </Tag>
      ),
    },
    {
      title: t('animalTable.reproduction.actions'),
      key: 'actions',
      render: (_: any, record: Reproduction) => (
        <Space size="middle">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleUpdatePhase(record)}
            title={t('animalTable.reproduction.updatePhase')}
          />
          <Popconfirm
            title={t('animalTable.reproduction.deleteConfirm')}
            onConfirm={() => handleDeleteReproduction(record.id)}
            okText={t('animalTable.reproduction.yes')}
            cancelText={t('animalTable.reproduction.no')}
          >
            <Button
              type="link"
              danger
              icon={<DeleteOutlined />}
              title={t('animalTable.reproduction.delete')}
            />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h3>{t('animalTable.reproduction.title')}</h3>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setIsCreateModalVisible(true)}
        >
          {t('animalTable.reproduction.addReproduction')}
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={reproductions}
        rowKey="id"
        loading={loading}
        pagination={{
          pageSize: 10,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total, range) =>
            `${range[0]}-${range[1]} de ${total} registros`,
        }}
        scroll={{ x: 800 }}
      />

      <CreateReproductionModal
        visible={isCreateModalVisible}
        onCancel={() => setIsCreateModalVisible(false)}
        onSuccess={handleCreateModalSuccess}
      />

      <UpdateReproductionPhaseModal
        visible={isUpdatePhaseModalVisible}
        onCancel={() => {
          setIsUpdatePhaseModalVisible(false);
          setSelectedReproduction(null);
        }}
        onSuccess={handleUpdatePhaseModalSuccess}
        reproduction={selectedReproduction}
      />
    </div>
  );
});
