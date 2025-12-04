import { useEffect } from 'react';
import { Modal, Form, Select, DatePicker, InputNumber, Tooltip, Space, message } from 'antd';
import { InfoCircleOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';
import { useWeight } from '../../hooks/useWeight';
import { useAnimals } from '../../hooks/useAnimals';
import { useResponsive } from '../../../../../hooks';
import { CreateOrUpdateWeightRequest } from '../../domain/model/weight';

interface CreateWeightModalProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
  preselectedAnimalId?: number;
}

export const CreateWeightModal = ({ 
  visible, 
  onCancel, 
  onSuccess, 
  preselectedAnimalId
}: CreateWeightModalProps) => {
  const { t } = useTranslation();
  const { createOrUpdateWeight, loading } = useWeight();
  const { animals = [], loading: animalsLoading } = useAnimals();
  const { isMobile, isTablet } = useResponsive();
  const [form] = Form.useForm();

  useEffect(() => {
    if (visible) {
      if (preselectedAnimalId) {
        form.setFieldsValue({ 
          animal_id: preselectedAnimalId,
          date: dayjs()
        });
      } else {
        form.setFieldsValue({
          date: dayjs()
        });
      }
    }
  }, [visible, preselectedAnimalId, form]);

  const handleSubmit = async (values: {
    animal_id: number;
    date: dayjs.Dayjs;
    animal_weight: number;
  }) => {
    const data: CreateOrUpdateWeightRequest = {
      animal_id: values.animal_id,
      date: values.date.format('YYYY-MM-DD'),
      animal_weight: values.animal_weight,
    };

    try {
      const result = await createOrUpdateWeight(data);
      
      if (result) {
        message.success(t('animalTable.weight.createdSuccessfully') || 'Peso registrado com sucesso');
        form.resetFields();
        onSuccess();
      } else {
        message.error(t('animalTable.weight.createError') || 'Erro ao registrar peso');
      }
    } catch {
      message.error(t('animalTable.weight.createError') || 'Erro ao registrar peso');
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title={
        <Space>
          {t('animalTable.weight.addWeight') || 'Adicionar Peso'}
          <Tooltip 
            title="Se o animal já tiver um peso registrado, ele será atualizado. Caso contrário, um novo registro será criado."
            placement="top"
          >
            <InfoCircleOutlined style={{ color: '#1890ff', cursor: 'help' }} />
          </Tooltip>
        </Space>
      }
      open={visible}
      onCancel={handleCancel}
      onOk={() => form.submit()}
      confirmLoading={loading}
      width={isMobile ? '95%' : isTablet ? '80%' : 600}
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
          name="animal_id"
          label={t('animalTable.weight.selectAnimal') || 'Selecionar Animal'}
          rules={[{ required: true, message: t('animalTable.weight.animalRequired') || 'Animal é obrigatório' }]}
        >
          <Select
            placeholder={t('animalTable.weight.selectAnimalPlaceholder') || 'Selecione um animal'}
            showSearch
            loading={animalsLoading}
            optionFilterProp="children"
            filterOption={(input: string, option?: { children: unknown }) =>
              (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
            }
          >
            {animals.map((animal) => (
              <Select.Option key={animal.id} value={parseInt(animal.id)}>
                {animal.animal_name} - {animal.ear_tag_number_local}
              </Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          name="animal_weight"
          label={t('animalTable.weight.weight') || 'Peso (kg)'}
          rules={[
            { required: true, message: t('animalTable.weight.weightRequired') || 'Peso é obrigatório' },
            { 
              type: 'number', 
              min: 0.1, 
              message: t('animalTable.weight.weightMin') || 'Peso deve ser maior que zero' 
            }
          ]}
        >
          <InputNumber
            style={{ width: '100%' }}
            placeholder={t('animalTable.weight.weightPlaceholder') || 'Digite o peso em kg'}
            min={0.1}
            step={0.1}
            precision={2}
            suffix=" kg"
          />
        </Form.Item>

        <Form.Item
          name="date"
          label={t('animalTable.weight.date') || 'Data'}
          rules={[{ required: true, message: t('animalTable.weight.dateRequired') || 'Data é obrigatória' }]}
        >
          <DatePicker 
            style={{ width: '100%' }} 
            format="DD/MM/YYYY"
            placeholder={t('animalTable.weight.datePlaceholder') || 'Selecione a data'}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

