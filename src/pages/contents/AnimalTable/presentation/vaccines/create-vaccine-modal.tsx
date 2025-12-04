import React, { useState, useEffect } from 'react'
import { Modal, Form, Input, Button, message, Space } from 'antd'
import { useTranslation } from 'react-i18next'
import { useVaccine } from '../../hooks/useVaccine'
import { useFarm } from '../../../../../hooks/useFarm'
import { useResponsive } from '../../../../../hooks'
import { CreateVaccineRequest } from '../../domain/model/vaccine'

interface CreateVaccineModalProps {
  visible: boolean
  onCancel: () => void
  onSuccess: () => void
}

export const CreateVaccineModal: React.FC<CreateVaccineModalProps> = ({
  visible,
  onCancel,
  onSuccess
}) => {
  const { t } = useTranslation()
  const { farm } = useFarm()
  const { isMobile, isTablet } = useResponsive()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  
  const { createVaccine } = useVaccine(farm?.id || 0)

  useEffect(() => {
    if (visible) {
      form.resetFields()
    }
  }, [visible, form])

  const handleSubmit = async (values: {
    name: string
    description?: string
    manufacturer?: string
  }) => {
    if (!farm?.id) {
      message.error('Fazenda não encontrada')
      return
    }

    setLoading(true)
    
    try {
      const data: CreateVaccineRequest = {
        farm_id: farm.id,
        name: values.name,
        description: values.description,
        manufacturer: values.manufacturer
      }

      await createVaccine(data)
      message.success(t('animalTable.vaccines.vaccineCreatedSuccessfully'))
      
      onSuccess()
      form.resetFields()
    } catch (err: unknown) {
      const error = err as Error
      message.error(error.message || t('animalTable.vaccines.createVaccineError'))
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = () => {
    form.resetFields()
    onCancel()
  }

  return (
    <Modal
      title={t('animalTable.vaccines.createVaccineTitle')}
      open={visible}
      onCancel={handleCancel}
      footer={null}
      width={isMobile ? '95%' : isTablet ? '80%' : 500}
      style={{
        top: isMobile ? '10px' : '50px'
      }}
      styles={{
        body: {
          maxHeight: isMobile ? '70vh' : '80vh',
          overflowY: 'auto',
          padding: isMobile ? '16px' : '24px'
        }
      }}
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
      >
        <Form.Item
          name="name"
          label={t('animalTable.vaccines.vaccineName')}
          rules={[
            { required: true, message: t('animalTable.vaccines.vaccineNameRequired') }
          ]}
        >
          <Input
            placeholder={t('animalTable.vaccines.vaccineNamePlaceholder')}
          />
        </Form.Item>

        <Form.Item
          name="description"
          label={t('animalTable.vaccines.description')}
        >
          <Input.TextArea
            rows={3}
            placeholder={t('animalTable.vaccines.descriptionPlaceholder')}
          />
        </Form.Item>

        <Form.Item
          name="manufacturer"
          label={t('animalTable.vaccines.manufacturer')}
        >
          <Input
            placeholder={t('animalTable.vaccines.manufacturerPlaceholder')}
          />
        </Form.Item>

        <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
          <Space>
            <Button onClick={handleCancel}>
              {t('animalTable.cancel')}
            </Button>
            <Button type="primary" htmlType="submit" loading={loading}>
              {t('animalTable.vaccines.create')}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  )
}

