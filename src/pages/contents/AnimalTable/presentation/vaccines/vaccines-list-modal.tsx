import React from 'react'
import { Modal, Table, Button, Tag, Spin, Alert } from 'antd'
import { useTranslation } from 'react-i18next'
import { useVaccine } from '../../hooks/useVaccine'
import { useFarm } from '../../../../../hooks/useFarm'
import { useResponsive } from '../../../../../hooks'
import { Vaccine } from '../../domain/model/vaccine'

interface VaccinesListModalProps {
  visible: boolean
  onCancel: () => void
}

export const VaccinesListModal: React.FC<VaccinesListModalProps> = ({
  visible,
  onCancel
}) => {
  const { t } = useTranslation()
  const { farm } = useFarm()
  const { isMobile } = useResponsive()
  const { vaccines, loading, error } = useVaccine(farm?.id || 0)

  const columns = [
    {
      title: t('animalTable.vaccines.name') || 'Nome',
      dataIndex: 'name',
      key: 'name',
      sorter: (a: Vaccine, b: Vaccine) => a.name.localeCompare(b.name),
      render: (text: string) => <strong>{text}</strong>
    },
    {
      title: t('animalTable.vaccines.description') || 'Descrição',
      dataIndex: 'description',
      key: 'description',
      render: (text: string) => text || '-'
    },
    {
      title: t('animalTable.vaccines.manufacturer') || 'Fabricante',
      dataIndex: 'manufacturer',
      key: 'manufacturer',
      render: (text: string) => text || '-'
    },
    {
      title: t('animalTable.vaccines.createdAt') || 'Data de Criação',
      dataIndex: 'created_at',
      key: 'created_at',
      sorter: (a: Vaccine, b: Vaccine) => 
        new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
      render: (date: string) => new Date(date).toLocaleDateString('pt-BR')
    }
  ]

  return (
    <Modal
      title={t('animalTable.vaccines.vaccinesList') || 'Lista de Vacinas'}
      open={visible}
      onCancel={onCancel}
      footer={[
        <Button key="close" onClick={onCancel}>
          {t('animalTable.vaccines.close') || 'Fechar'}
        </Button>
      ]}
      width={isMobile ? '95%' : 800}
      style={{
        top: isMobile ? '10px' : '50px'
      }}
    >
      {loading ? (
        <div style={{ textAlign: 'center', padding: '50px' }}>
          <Spin size="large" />
        </div>
      ) : error ? (
        <Alert
          message={t('animalTable.vaccines.error') || 'Erro'}
          description={error}
          type="error"
          showIcon
          style={{ marginBottom: '16px' }}
        />
      ) : (
        <>
          <div style={{ marginBottom: '16px' }}>
            <Tag color="blue">
              {t('animalTable.vaccines.totalVaccines') || 'Total de Vacinas'}: {vaccines.length}
            </Tag>
          </div>
          <Table
            columns={columns}
            dataSource={vaccines}
            rowKey="id"
            pagination={{
              pageSize: 10,
              showSizeChanger: !isMobile,
              showTotal: (total: number) => 
                `${t('animalTable.vaccines.total') || 'Total'}: ${total}`
            }}
            scroll={{
              x: isMobile ? 600 : undefined
            }}
            size={isMobile ? 'small' : 'middle'}
          />
        </>
      )}
    </Modal>
  )
}

