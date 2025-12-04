import { useEffect } from 'react';
import { Modal, Form, Select, DatePicker, Input, Switch, message } from 'antd';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';
import { useReproduction } from '../../hooks/useReproduction';
import { UpdateReproductionPhaseRequest, Reproduction, ReproductionPhase } from '../../domain/model/reproduction';

const { TextArea } = Input;

interface UpdateReproductionPhaseModalProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
  reproduction: Reproduction | null;
}

export const UpdateReproductionPhaseModal = ({ 
  visible, 
  onCancel, 
  onSuccess, 
  reproduction 
}: UpdateReproductionPhaseModalProps) => {
  const { t } = useTranslation();
  const { updateReproductionPhase, loading } = useReproduction();
  const [form] = Form.useForm();

  useEffect(() => {
    if (visible && reproduction) {
      form.setFieldsValue({
        new_phase: reproduction.current_phase,
        insemination_date: reproduction.insemination_date ? dayjs(reproduction.insemination_date) : undefined,
        insemination_type: reproduction.insemination_type,
        pregnancy_date: reproduction.pregnancy_date ? dayjs(reproduction.pregnancy_date) : undefined,
        lactation_start_date: reproduction.lactation_start_date ? dayjs(reproduction.lactation_start_date) : undefined,
        lactation_end_date: reproduction.lactation_end_date ? dayjs(reproduction.lactation_end_date) : undefined,
        dry_period_start_date: reproduction.dry_period_start_date ? dayjs(reproduction.dry_period_start_date) : undefined,
        actual_birth_date: reproduction.actual_birth_date ? dayjs(reproduction.actual_birth_date) : undefined,
        veterinary_confirmation: reproduction.veterinary_confirmation,
        observations: reproduction.observations,
      });
    }
  }, [visible, reproduction, form]);

  const handleSubmit = async (values: {
    new_phase: number;
    insemination_date?: dayjs.Dayjs;
    insemination_type?: string;
    pregnancy_date?: dayjs.Dayjs;
    lactation_start_date?: dayjs.Dayjs;
    lactation_end_date?: dayjs.Dayjs;
    dry_period_start_date?: dayjs.Dayjs;
    actual_birth_date?: dayjs.Dayjs;
    veterinary_confirmation?: boolean;
    observations?: string;
  }) => {
    if (!reproduction) return;

    const additionalData: UpdateReproductionPhaseRequest['additional_data'] = {};

    if (values.insemination_date) {
      additionalData.insemination_date = values.insemination_date.format('YYYY-MM-DD');
    }
    if (values.insemination_type) {
      additionalData.insemination_type = values.insemination_type;
    }
    if (values.pregnancy_date) {
      additionalData.pregnancy_date = values.pregnancy_date.format('YYYY-MM-DD');
    }
    if (values.lactation_start_date) {
      additionalData.lactation_start_date = values.lactation_start_date.format('YYYY-MM-DD');
    }
    if (values.lactation_end_date) {
      additionalData.lactation_end_date = values.lactation_end_date.format('YYYY-MM-DD');
    }
    if (values.dry_period_start_date) {
      additionalData.dry_period_start_date = values.dry_period_start_date.format('YYYY-MM-DD');
    }
    if (values.actual_birth_date) {
      additionalData.actual_birth_date = values.actual_birth_date.format('YYYY-MM-DD');
    }
    if (values.veterinary_confirmation !== undefined) {
      additionalData.veterinary_confirmation = values.veterinary_confirmation;
    }
    if (values.observations) {
      additionalData.observations = values.observations;
    }

    const data: UpdateReproductionPhaseRequest = {
      animal_id: reproduction.animal_id,
      new_phase: values.new_phase,
      additional_data: additionalData,
    };

    const success = await updateReproductionPhase(data);
    if (success) {
      message.success(t('animalTable.reproduction.phaseUpdatedSuccessfully'));
      form.resetFields();
      onSuccess();
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title={t('animalTable.reproduction.updatePhaseTitle')}
      open={visible}
      onCancel={handleCancel}
      onOk={() => form.submit()}
      confirmLoading={loading}
      width={600}
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
      >
        <Form.Item
          name="new_phase"
          label={t('animalTable.reproduction.newPhase')}
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
          name="actual_birth_date"
          label={t('animalTable.reproduction.actualBirthDate')}
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

