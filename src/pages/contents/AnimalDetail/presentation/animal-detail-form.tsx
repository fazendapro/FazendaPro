import React, { useState, useEffect } from 'react';
import { Form, Input, Select, DatePicker, Switch, Button, Upload, message, Spin } from 'antd';
import { UploadOutlined, SaveOutlined, CloseOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useAnimalDetailContext } from '../hooks';
import { AnimalDetailFormData, SEX_OPTIONS, ANIMAL_TYPE_OPTIONS, STATUS_OPTIONS, PURPOSE_OPTIONS } from '../types';
import dayjs from 'dayjs';

interface AnimalFormValues {
  animal_name: string;
  ear_tag_number_local: number;
  ear_tag_number_register: number;
  type: string;
  sex: number;
  breed: string;
  birth_date?: dayjs.Dayjs | null;
  animal_type: number;
  status: number;
  confinement: boolean;
  fertilization: boolean;
  castrated: boolean;
  purpose: number;
  current_batch: number;
  father_id?: number;
  mother_id?: number;
}

const { Option } = Select;

interface AnimalDetailFormProps {
  onSave: (data: AnimalDetailFormData, photoFile: File | null) => void;
  onCancel: () => void;
}

export const AnimalDetailForm: React.FC<AnimalDetailFormProps> = ({ onSave, onCancel }) => {
  const { t } = useTranslation();
  const { animal, fathers, mothers, loadingParents, uploadPhoto } = useAnimalDetailContext();
  const [form] = Form.useForm();
  const [photoFile, setPhotoFile] = useState<File | null>(null);

  useEffect(() => {
    if (animal) {

      form.setFieldsValue({
        animal_name: animal.animal_name,
        ear_tag_number_local: animal.ear_tag_number_local,
        ear_tag_number_register: animal.ear_tag_number_register,
        type: animal.type,
        sex: animal.sex,
        breed: animal.breed,
        birth_date: animal.birth_date ? dayjs(animal.birth_date) : null,
        animal_type: animal.animal_type,
        status: animal.status,
        confinement: Boolean(animal.confinement),
        fertilization: Boolean(animal.fertilization),
        castrated: Boolean(animal.castrated),
        purpose: animal.purpose,
        current_batch: animal.current_batch,
        father_id: animal.father_id,
        mother_id: animal.mother_id
      });
      
    }
  }, [animal, form]);

  const handleSave = () => {
    form.validateFields().then((values: AnimalFormValues) => {
      const formData: AnimalDetailFormData = {
        animal_name: values.animal_name,
        ear_tag_number_local: values.ear_tag_number_local,
        ear_tag_number_register: values.ear_tag_number_register,
        type: values.type,
        sex: values.sex,
        breed: values.breed,
        birth_date: values.birth_date ? values.birth_date.format('YYYY-MM-DD') : undefined,
        animal_type: values.animal_type,
        status: values.status,
        confinement: values.confinement,
        fertilization: values.fertilization,
        castrated: values.castrated,
        purpose: values.purpose,
        current_batch: values.current_batch,
        father_id: values.father_id,
        mother_id: values.mother_id
      };
      
      onSave(formData, photoFile);
    });
  };

  const handlePhotoUpload = async (file: File) => {
    setPhotoFile(file);
    try {
      const photoUrl = await uploadPhoto(file);
      if (photoUrl) {
        message.success(t('animalDetail.photoSelected'));
      } else {
        message.error(t('animalDetail.photoUploadError'));
      }
    } catch {
      message.error(t('animalDetail.photoUploadError'));
    }
    return false;
  };

  const uploadProps = {
    beforeUpload: handlePhotoUpload,
    showUploadList: false,
    accept: 'image/*',
  };

  if (loadingParents) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <Form form={form} layout="vertical" size="large">
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '16px' }}>
        <Form.Item
          label={t('animalDetail.name')}
          name="animal_name"
          rules={[{ required: true, message: t('animalDetail.nameRequired') }]}
        >
          <Input placeholder={t('animalDetail.namePlaceholder')} />
        </Form.Item>

        <Form.Item
          label={t('animalDetail.localEarTag')}
          name="ear_tag_number_local"
          rules={[{ required: true, message: t('animalDetail.localEarTagRequired') }]}
        >
          <Input type="number" placeholder={t('animalDetail.localEarTagPlaceholder')} />
        </Form.Item>

        <Form.Item
          label={t('animalDetail.registerEarTag')}
          name="ear_tag_number_register"
        >
          <Input type="number" placeholder={t('animalDetail.registerEarTagPlaceholder')} />
        </Form.Item>

        <Form.Item
          label={t('animalDetail.type')}
          name="type"
          rules={[{ required: true, message: t('animalDetail.typeRequired') }]}
        >
          <Input placeholder={t('animalDetail.typePlaceholder')} />
        </Form.Item>

        <Form.Item
          label={t('animalDetail.sex')}
          name="sex"
          rules={[{ required: true, message: t('animalDetail.sexRequired') }]}
        >
          <Select placeholder={t('animalDetail.selectSex')}>
            {SEX_OPTIONS.map(option => (
              <Option key={option.value} value={option.value}>
                {option.label}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label={t('animalDetail.breed')}
          name="breed"
          rules={[{ required: true, message: t('animalDetail.breedRequired') }]}
        >
          <Input placeholder={t('animalDetail.breedPlaceholder')} />
        </Form.Item>

        <Form.Item
          label={t('animalDetail.birthDate')}
          name="birth_date"
        >
          <DatePicker style={{ width: '100%' }} placeholder={t('animalDetail.birthDatePlaceholder')} />
        </Form.Item>

        <Form.Item
          label={t('animalDetail.animalType')}
          name="animal_type"
          rules={[{ required: true, message: t('animalDetail.animalTypeRequired') }]}
        >
          <Select placeholder={t('animalDetail.selectAnimalType')}>
            {ANIMAL_TYPE_OPTIONS.map(option => (
              <Option key={option.value} value={option.value}>
                {option.label}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label={t('animalDetail.status')}
          name="status"
          rules={[{ required: true, message: t('animalDetail.statusRequired') }]}
        >
          <Select placeholder={t('animalDetail.selectStatus')}>
            {STATUS_OPTIONS.map(option => (
              <Option key={option.value} value={option.value}>
                {option.label}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label={t('animalDetail.purpose')}
          name="purpose"
          rules={[{ required: true, message: t('animalDetail.purposeRequired') }]}
        >
          <Select placeholder={t('animalDetail.selectPurpose')}>
            {PURPOSE_OPTIONS.map(option => (
              <Option key={option.value} value={option.value}>
                {option.label}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label={t('animalDetail.currentBatch')}
          name="current_batch"
        >
          <Input type="number" placeholder={t('animalDetail.currentBatchPlaceholder')} />
        </Form.Item>

        <Form.Item
          label={t('animalDetail.father')}
          name="father_id"
        >
          <Select placeholder={t('animalDetail.selectFather')} allowClear>
            {fathers.map(father => (
              <Option key={father.id} value={father.id}>
                {father.animal_name} - {t('animalDetail.localEarTag')}: {father.ear_tag_number_local}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label={t('animalDetail.mother')}
          name="mother_id"
        >
          <Select placeholder={t('animalDetail.selectMother')} allowClear>
            {mothers.map(mother => (
              <Option key={mother.id} value={mother.id}>
                {mother.animal_name} - {t('animalDetail.localEarTag')}: {mother.ear_tag_number_local}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label={t('animalDetail.photo')}
          name="photoUpload"
          valuePropName="file"
          getValueFromEvent={() => undefined}
        >
          <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
            {photoFile && (
              <div style={{ width: '80px', height: '80px', border: '1px solid #d9d9d9', borderRadius: '4px', overflow: 'hidden' }}>
                <img 
                  src={URL.createObjectURL(photoFile)} 
                  alt="Preview" 
                  style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                />
              </div>
            )}
            <Upload {...uploadProps}>
              <Button icon={<UploadOutlined />}>
                {photoFile ? t('animalDetail.changePhoto') : t('animalDetail.selectPhoto')}
              </Button>
            </Upload>
          </div>
        </Form.Item>
      </div>

      <div style={{ display: 'flex', gap: '16px', marginTop: '24px' }}>
        <Form.Item name="confinement" valuePropName="checked">
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <Switch 
              key={`confinement-${animal?.id}-${animal?.confinement}`}
              checked={Boolean(animal?.confinement)}
              onChange={(checked: boolean) => form.setFieldValue('confinement', checked)} 
            />
            <span style={{ marginLeft: '8px' }}>{t('animalDetail.confinement')}</span>
          </div>
        </Form.Item>

        <Form.Item name="fertilization" valuePropName="checked">
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <Switch 
              key={`fertilization-${animal?.id}-${animal?.fertilization}`}
              checked={Boolean(animal?.fertilization)}
              onChange={(checked: boolean) => form.setFieldValue('fertilization', checked)} 
            />
            <span style={{ marginLeft: '8px' }}>{t('animalDetail.fertilization')}</span>
          </div>
        </Form.Item>

        <Form.Item name="castrated" valuePropName="checked">
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <Switch 
              key={`castrated-${animal?.id}-${animal?.castrated}`}
              checked={Boolean(animal?.castrated)}
              onChange={(checked: boolean) => form.setFieldValue('castrated', checked)} 
            />
            <span style={{ marginLeft: '8px' }}>{t('animalDetail.castrated')}</span>
          </div>
        </Form.Item>
      </div>

      <div style={{ display: 'flex', gap: '16px', justifyContent: 'flex-end', marginTop: '24px' }}>
        <Button onClick={onCancel} icon={<CloseOutlined />}>
          {t('common.cancel')}
        </Button>
        <Button type="primary" onClick={handleSave} icon={<SaveOutlined />}>
          {t('common.save')}
        </Button>
      </div>
    </Form>
  );
};
