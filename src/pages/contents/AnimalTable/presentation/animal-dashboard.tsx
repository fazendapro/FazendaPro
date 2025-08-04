import React from 'react';
import { Button, Input, Select, Modal, Card, Row, Col, Statistic, Space, Form, DatePicker } from 'antd';
import { useModal } from '../../../../hooks';
import { useTranslation } from 'react-i18next';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { animalSchema } from './animal-schema';
import { toast } from 'react-toastify';

const { Search } = Input;
const { Option } = Select;

interface AnimalForm {
  animalName: string;
  earringNumber: string;
  earringNumberGlobal: string;
  type: 'vaca' | 'bezerro' | 'touro' | 'novilho';
  sex: 'macho' | 'fêmea';
  breed: string;
  birthDate: string;
}

const AnimalDashboard: React.FC = () => {
  const { isOpen, onOpen, onClose } = useModal();
  const { t } = useTranslation();

  const methods = useForm<AnimalForm>({
    resolver: yupResolver(animalSchema),
    defaultValues: {
      animalName: '',
      earringNumber: '',
      earringNumberGlobal: '',
      type: 'vaca',
      sex: 'macho',
      breed: '',
      birthDate: ''
    },
  });

  const onSubmit = async (data: AnimalForm) => {
    try {
      console.log('Dados do animal:', data);
      toast.success('Animal adicionado com sucesso!');
      onClose();
      methods.reset();
    } catch (error) {
      toast.error(`Erro ao adicionar animal: ${error}`);
    }
  };

  const handleCancel = () => {
    onClose();
    methods.reset();
  };

  return (
    <div style={{ }}>
      <Card style={{ marginBottom: '16px' }}>
        <Row gutter={16}>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.categories')} 
              value={14} 
              valueStyle={{ color: '#1890ff' }}
            />
          </Col>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.totalAnimals')} 
              value={868} 
              valueStyle={{ color: '#faad14' }}
            />
          </Col>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.revenue')} 
              value={25000} 
              prefix="R$"
              valueStyle={{ color: '#000' }}
            />
          </Col>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.bestSales')} 
              value={5} 
              valueStyle={{ color: '#722ed1' }}
            />
          </Col>
        </Row>
        <Row gutter={16} style={{ marginTop: '16px' }}>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.cost')} 
              value={2500} 
              prefix="R$"
              valueStyle={{ color: '#000' }}
            />
          </Col>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.lowProduction')} 
              value={12} 
              valueStyle={{ color: '#ff4d4f' }}
            />
          </Col>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.inseminated')} 
              value={2} 
              valueStyle={{ color: '#000' }}
            />
          </Col>
          <Col span={6}>
            <Statistic 
              title={t('animalTable.notInseminated')} 
              value={2} 
              valueStyle={{ color: '#000' }}
            />
          </Col>
        </Row>
      </Card>

      <Card style={{ marginBottom: '16px' }}>
        <Space>
          <Button type="primary" onClick={onOpen}>
            {t('animalTable.createCow')}
          </Button>
          <Select defaultValue="filtro" style={{ width: 120 }}>
            <Option value="filtro">{t('animalTable.filter')}</Option>
          </Select>
          <Search 
            placeholder={t('animalTable.search')} 
            style={{ width: 200 }} 
          />
        </Space>
      </Card>

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
        <Form layout="vertical">
          <Row gutter={16}>
            <Col span={24}>
              <Form.Item
                label={t('animalTable.animalName')}
                validateStatus={methods.formState.errors.animalName ? 'error' : ''}
                help={methods.formState.errors.animalName?.message}
              >
                <Input
                  {...methods.register('animalName')}
                  placeholder={t('animalTable.animalNamePlaceholder')}
                />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label={t('animalTable.earringNumberLocal')}
                validateStatus={methods.formState.errors.earringNumber ? 'error' : ''}
                help={methods.formState.errors.earringNumber?.message}
              >
                <Input
                  {...methods.register('earringNumber')}
                  placeholder={t('animalTable.earringNumberLocalPlaceholder')}
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label={t('animalTable.earringNumberGlobal')}
                validateStatus={methods.formState.errors.earringNumberGlobal ? 'error' : ''}
                help={methods.formState.errors.earringNumberGlobal?.message}
              >
                <Input
                  {...methods.register('earringNumberGlobal')}
                  placeholder={t('animalTable.earringNumberGlobalPlaceholder')}
                />
              </Form.Item>
            </Col>
          </Row>

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
                  <Option value="macho">{t('animalTable.male')}</Option>
                  <Option value="fêmea">{t('animalTable.female')}</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label={t('animalTable.breed')}
                validateStatus={methods.formState.errors.breed ? 'error' : ''}
                help={methods.formState.errors.breed?.message}
              >
                <Input
                  {...methods.register('breed')}
                  placeholder={t('animalTable.breedPlaceholder')}
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label={t('animalTable.birthDate')}
                validateStatus={methods.formState.errors.birthDate ? 'error' : ''}
                help={methods.formState.errors.birthDate?.message}
              >
                <DatePicker
                  style={{ width: '100%' }}
                  onChange={(date, dateString) => {
                    if (typeof dateString === 'string') {
                      methods.setValue('birthDate', dateString);
                    }
                  }}
                  placeholder={t('animalTable.birthDatePlaceholder')}
                />
              </Form.Item>
            </Col>
          </Row>

        </Form>
      </Modal>
    </div>
  );
};

export { AnimalDashboard };