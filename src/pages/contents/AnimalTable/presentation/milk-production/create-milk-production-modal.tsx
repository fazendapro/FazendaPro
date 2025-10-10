import React, { useState, useEffect } from 'react'
import { Modal, Form, Input, DatePicker, Select, Button, message, Space } from 'antd'
import { useTranslation } from 'react-i18next'
import { useAnimals } from '../../hooks/useAnimals'
import { useMilkProduction } from '../../hooks/useMilkProduction'
import { useFarm } from '../../../../../hooks/useFarm'
import { useResponsive } from '../../../../../hooks'
import { CreateMilkProductionRequest, MilkProduction } from '../../domain/model/milk-production'
import { UpdateMilkProductionRequest } from '../../domain/usecases/update-milk-production-use-case'
import dayjs from 'dayjs'

const { Option } = Select

interface MilkProductionModalProps {
  visible: boolean
  onCancel: () => void
  onSuccess: () => void
  preselectedAnimalId?: number
  editingProduction?: MilkProduction
}

export const CreateMilkProductionModal: React.FC<MilkProductionModalProps> = ({
  visible,
  onCancel,
  onSuccess,
  preselectedAnimalId,
  editingProduction
}) => {
  const { t } = useTranslation()
  const { farm } = useFarm()
  const { isMobile, isTablet } = useResponsive()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  
  const { animals = [], loading: animalsLoading } = useAnimals(farm?.id)
  const { createMilkProduction, updateMilkProduction } = useMilkProduction(farm?.id)
  
  const isEditing = !!editingProduction

  useEffect(() => {
    if (visible) {
      form.resetFields()
      if (isEditing && editingProduction) {
        form.setFieldsValue({
          animalId: editingProduction.animal_id,
          liters: editingProduction.liters,
          date: dayjs(editingProduction.date)
        })
      } else {
        form.setFieldsValue({
          date: dayjs(),
          animalId: preselectedAnimalId
        })
      }
    }
  }, [visible, form, preselectedAnimalId, isEditing, editingProduction])

  const handleSubmit = async (values: {
    animalId: number
    liters: number
    date: dayjs.Dayjs
  }) => {
    setLoading(true)
    
    try {
      if (isEditing && editingProduction) {
        const data: UpdateMilkProductionRequest = {
          id: editingProduction.id,
          animal_id: values.animalId,
          liters: values.liters,
          date: values.date.format('YYYY-MM-DD')
        }

        const result = await updateMilkProduction(data)
        message.success(t('animalTable.milkProductionContainer.updatedSuccessfully'))
      } else {
        const data: CreateMilkProductionRequest = {
          animal_id: values.animalId,
          liters: values.liters,
          date: values.date.format('YYYY-MM-DD')
        }

        const result = await createMilkProduction(data)
        message.success(t('animalTable.milkProductionContainer.createdSuccessfully'))
      }
      
      onSuccess()
      form.resetFields()
    } catch (error) {
      const errorMessage = isEditing 
        ? t('animalTable.milkProductionContainer.updateError')
        : t('animalTable.milkProductionContainer.createError')
      message.error(errorMessage)
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
      title={isEditing 
        ? t('animalTable.milkProductionContainer.editTitle') 
        : t('animalTable.milkProductionContainer.createTitle')
      }
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
        initialValues={{
          date: dayjs()
        }}
      >
        <Form.Item
          name="animalId"
          label={t('animalTable.milkProductionContainer.selectAnimal')}
          rules={[
            { required: true, message: t('animalTable.milkProductionContainer.animalRequired') }
          ]}
        >
          <Select
            placeholder={t('animalTable.milkProductionContainer.selectAnimalPlaceholder')}
            loading={animalsLoading}
            showSearch
            optionFilterProp="children"
            filterOption={(input, option) => {
              const children = String(option?.children || '')
              return children.toLowerCase().includes(input.toLowerCase())
            }}
          >
            {animals.map(animal => (
              <Option key={animal.id} value={animal.id}>
                {animal.animal_name} - {animal.ear_tag_number_local}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          name="liters"
          label={t('animalTable.milkProductionContainer.liters')}
          rules={[
            { required: true, message: t('animalTable.milkProductionContainer.litersRequired') },
            { 
              pattern: /^\d+(\.\d{1,2})?$/, 
              message: t('animalTable.milkProductionContainer.litersFormat') 
            }
          ]}
        >
          <Input
            type="number"
            step="0.1"
            min="0"
            placeholder={t('animalTable.milkProductionContainer.litersPlaceholder')}
            suffix="L"
          />
        </Form.Item>

        <Form.Item
          name="date"
          label={t('animalTable.milkProductionContainer.date')}
          rules={[
            { required: true, message: t('animalTable.milkProductionContainer.dateRequired') }
          ]}
        >
          <DatePicker
            style={{ width: '100%' }}
            format="DD/MM/YYYY"
            placeholder={t('animalTable.milkProductionContainer.datePlaceholder')}
          />
        </Form.Item>

        <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
          <Space>
            <Button onClick={handleCancel}>
              {t('animalTable.cancel')}
            </Button>
            <Button type="primary" htmlType="submit" loading={loading}>
              {isEditing
                ? t('animalTable.milkProductionContainer.update')
                : t('animalTable.milkProductionContainer.create')
              }
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  )
}
