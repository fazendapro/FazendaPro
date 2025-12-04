import { useEffect } from 'react';
import { Modal, Form, Select, DatePicker, Input, Switch, message } from 'antd';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';
import { useReproduction } from '../../hooks/useReproduction';
import { useAnimals } from '../../hooks/useAnimals';
import { useResponsive } from '../../../../../hooks';
import { CreateReproductionRequest, Reproduction, ReproductionPhase } from '../../domain/model/reproduction';

const { TextArea } = Input;

interface CreateReproductionModalProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
  preselectedAnimalId?: number;
  editingReproduction?: Reproduction;
}

export const CreateReproductionModal = ({ 
  visible, 
  onCancel, 
  onSuccess, 
  preselectedAnimalId,
  editingReproduction
}: CreateReproductionModalProps) => {
  const { t } = useTranslation();
  const { createReproduction, updateReproduction, loading } = useReproduction();
  const { animals = [], loading: animalsLoading } = useAnimals();
  const { isMobile, isTablet } = useResponsive();
  const [form] = Form.useForm();

  useEffect(() => {
    if (visible) {
      if (editingReproduction) {
        form.setFieldsValue({
          animal_id: editingReproduction.animal_id,
          current_phase: editingReproduction.current_phase,
          insemination_date: editingReproduction.insemination_date ? dayjs(editingReproduction.insemination_date) : undefined,
          insemination_type: editingReproduction.insemination_type,
          pregnancy_date: editingReproduction.pregnancy_date ? dayjs(editingReproduction.pregnancy_date) : undefined,
          expected_birth_date: editingReproduction.expected_birth_date ? dayjs(editingReproduction.expected_birth_date) : undefined,
          actual_birth_date: editingReproduction.actual_birth_date ? dayjs(editingReproduction.actual_birth_date) : undefined,
          lactation_start_date: editingReproduction.lactation_start_date ? dayjs(editingReproduction.lactation_start_date) : undefined,
          lactation_end_date: editingReproduction.lactation_end_date ? dayjs(editingReproduction.lactation_end_date) : undefined,
          dry_period_start_date: editingReproduction.dry_period_start_date ? dayjs(editingReproduction.dry_period_start_date) : undefined,
          veterinary_confirmation: editingReproduction.veterinary_confirmation,
          observations: editingReproduction.observations,
        });
      } else if (preselectedAnimalId) {
        form.setFieldsValue({ animal_id: preselectedAnimalId });
      }
    }
  }, [visible, preselectedAnimalId, editingReproduction, form]);

  const handleSubmit = async (values: {
    animal_id: number;
    current_phase: number;
    insemination_date?: dayjs.Dayjs;
    insemination_type?: string;
    pregnancy_date?: dayjs.Dayjs;
    expected_birth_date?: dayjs.Dayjs;
    actual_birth_date?: dayjs.Dayjs;
    lactation_start_date?: dayjs.Dayjs;
    lactation_end_date?: dayjs.Dayjs;
    dry_period_start_date?: dayjs.Dayjs;
    veterinary_confirmation?: boolean;
    observations?: string;
  }) => {
    const data: CreateReproductionRequest = {
      animal_id: values.animal_id,
      current_phase: values.current_phase,
      insemination_date: values.insemination_date?.format('YYYY-MM-DD'),
      insemination_type: values.insemination_type,
      pregnancy_date: values.pregnancy_date?.format('YYYY-MM-DD'),
      expected_birth_date: values.expected_birth_date?.format('YYYY-MM-DD'),
      actual_birth_date: values.actual_birth_date?.format('YYYY-MM-DD'),
      lactation_start_date: values.lactation_start_date?.format('YYYY-MM-DD'),
      lactation_end_date: values.lactation_end_date?.format('YYYY-MM-DD'),
      dry_period_start_date: values.dry_period_start_date?.format('YYYY-MM-DD'),
      veterinary_confirmation: values.veterinary_confirmation || false,
      observations: values.observations,
    };

    try {
      if (editingReproduction) {
        const updateData = { ...data, id: editingReproduction.id };
        const success = await updateReproduction(updateData);
        
        if (success) {
          message.success(t('animalTable.reproduction.updatedSuccessfully'));
          form.resetFields();
          onSuccess();
        } else {
          message.error(t('animalTable.reproduction.updateError'));
        }
      } else {
        const result = await createReproduction(data);
        
        if (result) {
          message.success(t('animalTable.reproduction.createdSuccessfully'));
          form.resetFields();
          onSuccess();
        } else {
          message.error('Erro ao criar registro de reprodução');
        }
      }
    } catch {
      message.error(editingReproduction ? 'Erro ao atualizar registro de reprodução' : 'Erro ao criar registro de reprodução');
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title={editingReproduction ? t('animalTable.reproduction.editTitle') : t('animalTable.reproduction.createTitle')}
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
          current_phase: ReproductionPhase.VAZIAS,
          veterinary_confirmation: false,
        }}
      >
        <Form.Item
          name="animal_id"
          label={t('animalTable.reproduction.selectAnimal')}
          rules={[{ required: true, message: t('animalTable.reproduction.animalRequired') }]}
        >
        <Select
          placeholder={t('animalTable.reproduction.selectAnimalPlaceholder')}
          showSearch
          loading={animalsLoading}
          optionFilterProp="children"
          filterOption={(input: string, option?: { children?: React.ReactNode }) =>
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
          name="current_phase"
          label={t('animalTable.reproduction.currentPhase')}
          rules={[{ required: true, message: t('animalTable.reproduction.phaseRequired') }]}
        >
          <Select>
            <Select.Option value={ReproductionPhase.LACTACAO}>
              {t('animalTable.reproduction.phases.lactacao')}
            </Select.Option>
            <Select.Option value={ReproductionPhase.SECANDO}>
              {t('animalTable.reproduction.phases.secando')}
            </Select.Option>
            <Select.Option value={ReproductionPhase.VAZIAS}>
              {t('animalTable.reproduction.phases.vazias')}
            </Select.Option>
            <Select.Option value={ReproductionPhase.PRENHAS}>
              {t('animalTable.reproduction.phases.prenhas')}
            </Select.Option>
          </Select>
        </Form.Item>

        <Form.Item
          name="insemination_date"
          label={t('animalTable.reproduction.inseminationDate')}
        >
          <DatePicker style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item
          name="insemination_type"
          label={t('animalTable.reproduction.inseminationType')}
        >
          <Select placeholder={t('animalTable.reproduction.selectInseminationType')}>
            <Select.Option value="Natural">Natural</Select.Option>
            <Select.Option value="Artificial">Inseminação Artificial</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item
          name="pregnancy_date"
          label={t('animalTable.reproduction.pregnancyDate')}
        >
          <DatePicker style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item
          name="expected_birth_date"
          label={t('animalTable.reproduction.expectedBirthDate')}
        >
          <DatePicker style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item
          name="actual_birth_date"
          label={t('animalTable.reproduction.actualBirthDate')}
        >
          <DatePicker style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item
          name="lactation_start_date"
          label={t('animalTable.reproduction.lactationStartDate')}
        >
          <DatePicker style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item
          name="lactation_end_date"
          label={t('animalTable.reproduction.lactationEndDate')}
        >
          <DatePicker style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item
          name="dry_period_start_date"
          label={t('animalTable.reproduction.dryPeriodStartDate')}
        >
          <DatePicker style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item
          name="veterinary_confirmation"
          label={t('animalTable.reproduction.veterinaryConfirmation')}
          valuePropName="checked"
        >
          <Switch />
        </Form.Item>

        <Form.Item
          name="observations"
          label={t('animalTable.reproduction.observations')}
        >
          <TextArea rows={3} placeholder={t('animalTable.reproduction.observationsPlaceholder')} />
        </Form.Item>
      </Form>
    </Modal>
  );
};
