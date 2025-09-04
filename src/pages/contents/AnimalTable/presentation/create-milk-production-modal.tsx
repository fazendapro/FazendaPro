import { useState, useEffect } from 'react'
import { Modal, Form, Input, DatePicker, Select, Button, message, Space } from 'antd'
import { useTranslation } from 'react-i18next'
import { useAnimals } from '../hooks/useAnimals'
import { useMilkProduction } from '../hooks/useMilkProduction'
import { useFarm } from '../../../../hooks/useFarm'
import { CreateMilkProductionRequest } from '../types/milk-production'
import dayjs from 'dayjs'

const { Option } = Select

interface CreateMilkProductionModalProps {
  visible: boolean
  onCancel: () => void
  onSuccess: () => void
}

export const CreateMilkProductionModal: React.FC<CreateMilkProductionModalProps> = ({
  visible,
  onCancel,
  onSuccess
}) => {
  const { t } = useTranslation()
  const { farm } = useFarm()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  
  const { animals, loading: animalsLoading } = useAnimals(farm.id)
  const { createMilkProduction } = useMilkProduction(farm.id)

  useEffect(() => {
    if (visible) {
      form.resetFields()
      form.setFieldsValue({
        date: dayjs()
      })
    }
  }, [visible, form])

  const handleSubmit = async (values: any) => {
    setLoading(true)
    
    try {
      const data: CreateMilkProductionRequest = {
        animalId: values.animalId,
        liters: parseFloat(values.liters),
        date: values.date.format('YYYY-MM-DD')
      }

      await createMilkProduction(data)
      message.success(t('animalTable.milkProductionContainer.createdSuccessfully'))
      onSuccess()
      form.resetFields()
    } catch {
      message.error(t('animalTable.milkProductionContainer.createError'))
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
      title={t('animalTable.milkProductionContainer.createTitle')}
      open={visible}
      onCancel={handleCancel}
      footer={null}
      width={500}
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
            filterOption={(input, option) =>
              (option?.children as unknown as string)
                ?.toLowerCase()
                .includes(input.toLowerCase())
            }
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
              {t('animalTable.milkProductionContainer.create')}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  )
}
