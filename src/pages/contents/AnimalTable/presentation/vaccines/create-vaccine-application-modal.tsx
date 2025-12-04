import React, { useState, useEffect } from 'react'
import { Modal, Form, Input, DatePicker, Select, Button, message, Space } from 'antd'
import { useTranslation } from 'react-i18next'
import { useAnimals } from '../../hooks/useAnimals'
import { useVaccine } from '../../hooks/useVaccine'
import { useVaccineApplication } from '../../hooks/useVaccineApplication'
import { useFarm } from '../../../../../hooks/useFarm'
import { useResponsive } from '../../../../../hooks'
import { CreateVaccineApplicationRequest, VaccineApplication, UpdateVaccineApplicationRequest } from '../../domain/model/vaccine-application'
import dayjs from 'dayjs'

const { Option } = Select

interface CreateVaccineApplicationModalProps {
  visible: boolean
  onCancel: () => void
  onSuccess: () => void
  preselectedAnimalId?: number
  editingApplication?: VaccineApplication
}

export const CreateVaccineApplicationModal: React.FC<CreateVaccineApplicationModalProps> = ({
  visible,
  onCancel,
  onSuccess,
  preselectedAnimalId,
  editingApplication
}) => {
  const { t } = useTranslation()
  const { farm } = useFarm()
  const { isMobile, isTablet } = useResponsive()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  
  const { animals = [], loading: animalsLoading } = useAnimals(farm?.id || 0)
  const { vaccines = [], loading: vaccinesLoading } = useVaccine(farm?.id || 0)
  const { createVaccineApplication, updateVaccineApplication } = useVaccineApplication(farm?.id || 0)
  
  const isEditing = !!editingApplication

  useEffect(() => {
    if (visible) {
      form.resetFields()
      if (isEditing && editingApplication) {
        form.setFieldsValue({
          animalId: editingApplication.animal_id,
          vaccineId: editingApplication.vaccine_id,
          applicationDate: dayjs(editingApplication.application_date),
          batchNumber: editingApplication.batch_number,
          veterinarian: editingApplication.veterinarian,
          observations: editingApplication.observations
        })
      } else {
        form.setFieldsValue({
          applicationDate: dayjs(),
          animalId: preselectedAnimalId
        })
      }
    }
  }, [visible, form, preselectedAnimalId, isEditing, editingApplication])

  const handleSubmit = async (values: {
    animalId: number
    vaccineId: number
    applicationDate: dayjs.Dayjs
    batchNumber?: string
    veterinarian?: string
    observations?: string
  }) => {
    setLoading(true)
    
    try {
      if (isEditing && editingApplication) {
        const data: UpdateVaccineApplicationRequest = {
          animal_id: values.animalId,
          vaccine_id: values.vaccineId,
          application_date: values.applicationDate.format('YYYY-MM-DD'),
          batch_number: values.batchNumber,
          veterinarian: values.veterinarian,
          observations: values.observations
        }

        await updateVaccineApplication(editingApplication.id, data)
        message.success(t('animalTable.vaccines.applicationUpdatedSuccessfully'))
      } else {
        const data: CreateVaccineApplicationRequest = {
          animal_id: values.animalId,
          vaccine_id: values.vaccineId,
          application_date: values.applicationDate.format('YYYY-MM-DD'),
          batch_number: values.batchNumber,
          veterinarian: values.veterinarian,
          observations: values.observations
        }

        await createVaccineApplication(data)
        message.success(t('animalTable.vaccines.applicationCreatedSuccessfully'))
      }
      
      onSuccess()
      form.resetFields()
    } catch (err: unknown) {
      const error = err as Error
      const errorMessage = isEditing 
        ? t('animalTable.vaccines.updateApplicationError')
        : t('animalTable.vaccines.createApplicationError')
      message.error(error.message || errorMessage)
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
        ? t('animalTable.vaccines.editApplicationTitle') 
        : t('animalTable.vaccines.createApplicationTitle')
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
          applicationDate: dayjs()
        }}
      >
        <Form.Item
          name="animalId"
          label={t('animalTable.vaccines.selectAnimal')}
          rules={[
            { required: true, message: t('animalTable.vaccines.animalRequired') }
          ]}
        >
          <Select
            placeholder={t('animalTable.vaccines.selectAnimalPlaceholder')}
            loading={animalsLoading}
            showSearch
            optionFilterProp="children"
            filterOption={(input: string, option?: { children?: React.ReactNode }) => {
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
          name="vaccineId"
          label={t('animalTable.vaccines.selectVaccine')}
          rules={[
            { required: true, message: t('animalTable.vaccines.vaccineRequired') }
          ]}
        >
          <Select
            placeholder={t('animalTable.vaccines.selectVaccinePlaceholder')}
            loading={vaccinesLoading}
            showSearch
            optionFilterProp="children"
            filterOption={(input: string, option?: { children?: React.ReactNode }) => {
              const children = String(option?.children || '')
              return children.toLowerCase().includes(input.toLowerCase())
            }}
          >
            {vaccines.map(vaccine => (
              <Option key={vaccine.id} value={vaccine.id}>
                {vaccine.name}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          name="applicationDate"
          label={t('animalTable.vaccines.applicationDate')}
          rules={[
            { required: true, message: t('animalTable.vaccines.applicationDateRequired') }
          ]}
        >
          <DatePicker
            style={{ width: '100%' }}
            format="DD/MM/YYYY"
            placeholder={t('animalTable.vaccines.applicationDatePlaceholder')}
          />
        </Form.Item>

        <Form.Item
          name="batchNumber"
          label={t('animalTable.vaccines.batchNumber')}
        >
          <Input
            placeholder={t('animalTable.vaccines.batchNumberPlaceholder')}
          />
        </Form.Item>

        <Form.Item
          name="veterinarian"
          label={t('animalTable.vaccines.veterinarian')}
        >
          <Input
            placeholder={t('animalTable.vaccines.veterinarianPlaceholder')}
          />
        </Form.Item>

        <Form.Item
          name="observations"
          label={t('animalTable.vaccines.observations')}
        >
          <Input.TextArea
            rows={3}
            placeholder={t('animalTable.vaccines.observationsPlaceholder')}
          />
        </Form.Item>

        <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
          <Space>
            <Button onClick={handleCancel}>
              {t('animalTable.cancel')}
            </Button>
            <Button type="primary" htmlType="submit" loading={loading}>
              {isEditing
                ? t('animalTable.vaccines.update')
                : t('animalTable.vaccines.create')
              }
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  )
}

