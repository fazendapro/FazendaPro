import React from 'react';
import { Button, Select, Modal, Row, Col, Form } from 'antd';
import { useTranslation } from 'react-i18next';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { animalSchema } from './animal-schema';
import { toast } from 'react-toastify';
import { CreateAnimalFactory } from '../factories';
import { useCsrfTokenContext } from '../../../../contexts';
import { Form as CustomForm } from '../../../../components';
import { FieldType } from '../../../../types/field-types';
import { AnimalForm, AnimalSex } from '../types/type';

const { Option } = Select;

interface CreateAnimalModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const CreateAnimalModal: React.FC<CreateAnimalModalProps> = ({ isOpen, onClose }) => {
  const { t } = useTranslation();
  const { csrfToken } = useCsrfTokenContext();

  const methods = useForm<AnimalForm>({
    resolver: yupResolver(animalSchema),
    defaultValues: {
      animal_name: '',
      ear_tag_number_local: 0,
      ear_tag_number_register: 0,
      type: 'vaca',
      sex: AnimalSex.MALE,
      breed: '',
      birth_date: ''
    },
  });

  const onSubmit = async (data: AnimalForm) => {
    try {
      const farmId = 1; // TODO: get farmId from context

      const createAnimalUseCase = CreateAnimalFactory(csrfToken);
      const response = await createAnimalUseCase.create({ ...data, farm_id: farmId });

      if (response?.success) {
        toast.success(response.message);
        onClose();
        methods.reset();
      } else {
        toast.error('Erro ao criar animal');
      }
    } catch (error) {
      toast.error(`Erro ao adicionar animal: ${error instanceof Error ? error.message : 'Erro desconhecido'}`);
    }
  };

  const handleCancel = () => {
    onClose();
    methods.reset();
  };

  const animalFields: FieldType[] = [
    {
      name: 'animal_name',
      label: t('animalTable.animalName'),
      type: 'text',
      placeholder: t('animalTable.animalNamePlaceholder'),
      isRequired: true,
      colSpan: 24
    },
    {
      name: 'ear_tag_number_local',
      label: t('animalTable.earringNumberLocal'),
      type: 'number',
      placeholder: t('animalTable.earringNumberLocalPlaceholder'),
      isRequired: true,
      colSpan: 12
    },
    {
      name: 'ear_tag_number_register',
      label: t('animalTable.earringNumberGlobal'),
      type: 'number',
      placeholder: t('animalTable.earringNumberGlobalPlaceholder'),
      isRequired: true,
      colSpan: 12
    },
    {
      name: 'breed',
      label: t('animalTable.breed'),
      type: 'text',
      placeholder: t('animalTable.breedPlaceholder'),
      isRequired: true,
      colSpan: 12
    }
  ];

  return (
    <Modal
      title={t('animalTable.newCattle')}
      open={isOpen}
      onCancel={handleCancel}
      footer={[
        <Button key="back" onClick={handleCancel}>
          {t('animalTable.cancel')}
        </Button>,
        <Button 
          key="submit" 
          type="primary" 
          onClick={methods.handleSubmit(onSubmit)}
          loading={methods.formState.isSubmitting}
        >
          {t('animalTable.addCattle')}
        </Button>,
      ]}
      width={600}
    >
      <CustomForm<AnimalForm>
        onSubmit={onSubmit}
        fields={animalFields}
        methods={methods}
      >
        <Row gutter={16}>
          <Col span={12}>
            <Form.Item
              label={t('animalTable.type')}
              validateStatus={methods.formState.errors.type ? 'error' : ''}
              help={methods.formState.errors.type?.message}
            >
              <Select
                value={methods.watch('type')}
                placeholder={t('animalTable.selectAnimalType')}
                onChange={(value) => methods.setValue('type', value)}
              >
                <Option value="vaca">{t('animalTable.cow')}</Option>
                <Option value="bezerro">{t('animalTable.calf')}</Option>
                <Option value="touro">{t('animalTable.bull')}</Option>
                <Option value="novilho">{t('animalTable.heifer')}</Option>
              </Select>
            </Form.Item>
          </Col>
          <Col span={12}>
            <Form.Item
              label={t('animalTable.sex')}
              validateStatus={methods.formState.errors.sex ? 'error' : ''}
              help={methods.formState.errors.sex?.message}
            >
              <Select
                value={methods.watch('sex')}
                placeholder={t('animalTable.selectAnimalSex')}
                onChange={(value) => methods.setValue('sex', value)}
              >
                <Option value={AnimalSex.MALE}>{t('animalTable.male')}</Option>
                <Option value={AnimalSex.FEMALE}>{t('animalTable.female')}</Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>
        <Row gutter={16}>
          <Col span={12}>
            <Form.Item
              label={t('animalTable.birthDate')}
              validateStatus={methods.formState.errors.birth_date ? 'error' : ''}
              help={methods.formState.errors.birth_date?.message}
            >
              <input
                type="date"
                style={{ width: '100%', padding: '8px', border: '1px solid #d9d9d9', borderRadius: '6px' }}
                onChange={(e) => methods.setValue('birth_date', e.target.value)}
                placeholder={t('animalTable.birthDatePlaceholder')}
              />
            </Form.Item>
          </Col>
        </Row>
      </CustomForm>
    </Modal>
  );
};

export { CreateAnimalModal }; 