import React from 'react';
import { Card, Button, Image, Tag, Space, Spin, Avatar, Row, Col, Typography, Divider } from 'antd';
import { EditOutlined, CameraOutlined, UserOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useAnimalDetailContext } from '../hooks';
import { SEX_OPTIONS } from '../types';
import { AnimalHistoryExport } from '../../../../components/AnimalHistoryExport/AnimalHistoryExport';

const { Title, Text } = Typography;

interface AnimalDetailDisplayProps {
  onEdit: () => void;
}

export const AnimalDetailDisplay: React.FC<AnimalDetailDisplayProps> = ({ onEdit }) => {
  const { t } = useTranslation();
  const { animal, loading } = useAnimalDetailContext();

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!animal) {
    return (
      <Card>
        <div style={{ textAlign: 'center', padding: '50px' }}>
          <p>{t('animalDetail.notFound')}</p>
        </div>
      </Card>
    );
  }

  const getSexLabel = (sex: number) => {
    return SEX_OPTIONS.find(option => option.value === sex)?.label || t('common.notInformed');
  };

  return (
    <div style={{ padding: '24px', backgroundColor: '#f5f5f5', minHeight: '100vh' }}>
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center', 
        marginBottom: '24px',
        backgroundColor: 'white',
        padding: '24px',
        borderRadius: '8px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
      }}>
        <Title level={1} style={{ margin: 0, color: '#262626' }}>
          {animal.animal_name}
        </Title>
        <Space>
          <Button 
            type="primary" 
            icon={<EditOutlined />} 
            onClick={onEdit}
            style={{ 
              backgroundColor: '#1890ff',
              borderColor: '#1890ff',
              borderRadius: '6px'
            }}
          >
            {t('common.edit')}
          </Button>
          <AnimalHistoryExport
            animal={animal}
            sales={[]}
            milkCollections={[]}
            reproductions={[]}
          />
        </Space>
      </div>

      <Row gutter={[24, 24]}>
        <Col xs={24} lg={16}>
          <Card 
            title={
              <Title level={4} style={{ margin: 0, color: '#262626' }}>
                {t('animalDetail.information')}
              </Title>
            }
            style={{ 
              borderRadius: '8px',
              boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
              border: 'none'
            }}
            styles={{ body: { padding: '24px' } }}
          >
            <Row gutter={[16, 16]}>
              <Col xs={24} sm={12}>
                <div style={{ marginBottom: '16px' }}>
                  <Text strong style={{ color: '#595959', fontSize: '14px' }}>
                    {t('animalDetail.name')}
                  </Text>
                  <div style={{ fontSize: '16px', color: '#262626', marginTop: '4px' }}>
                    {animal.animal_name}
                  </div>
                </div>
              </Col>
              <Col xs={24} sm={12}>
                <div style={{ marginBottom: '16px' }}>
                  <Text strong style={{ color: '#595959', fontSize: '14px' }}>
                    {t('animalDetail.localEarTag')}
                  </Text>
                  <div style={{ fontSize: '16px', color: '#262626', marginTop: '4px' }}>
                    {animal.ear_tag_number_local}
                  </div>
                </div>
              </Col>
              <Col xs={24} sm={12}>
                <div style={{ marginBottom: '16px' }}>
                  <Text strong style={{ color: '#595959', fontSize: '14px' }}>
                    {t('animalDetail.type')}
                  </Text>
                  <div style={{ fontSize: '16px', color: '#262626', marginTop: '4px' }}>
                    {animal.type}
                  </div>
                </div>
              </Col>
              <Col xs={24} sm={12}>
                <div style={{ marginBottom: '16px' }}>
                  <Text strong style={{ color: '#595959', fontSize: '14px' }}>
                    {t('animalDetail.sex')}
                  </Text>
                  <div style={{ marginTop: '4px' }}>
                    <Tag color={animal.sex === 0 ? 'blue' : 'pink'} style={{ borderRadius: '4px' }}>
                      {getSexLabel(animal.sex)}
                    </Tag>
                  </div>
                </div>
              </Col>
              <Col xs={24} sm={12}>
                <div style={{ marginBottom: '16px' }}>
                  <Text strong style={{ color: '#595959', fontSize: '14px' }}>
                    {t('animalDetail.breed')}
                  </Text>
                  <div style={{ fontSize: '16px', color: '#262626', marginTop: '4px' }}>
                    {animal.breed}
                  </div>
                </div>
              </Col>
              <Col xs={24} sm={12}>
                <div style={{ marginBottom: '16px' }}>
                  <Text strong style={{ color: '#595959', fontSize: '14px' }}>
                    {t('animalDetail.birthDate')}
                  </Text>
                  <div style={{ fontSize: '16px', color: '#262626', marginTop: '4px' }}>
                    {animal.birth_date ? new Date(animal.birth_date).toLocaleDateString('pt-BR') : t('common.notInformed')}
                  </div>
                </div>
              </Col>
            </Row>

            <Divider style={{ margin: '24px 0' }} />

            <div>
              <Title level={5} style={{ color: '#262626', marginBottom: '16px' }}>
                {t('animalDetail.characteristics')}
              </Title>
              <Space wrap>
                {animal.confinement && (
                  <Tag color="orange" style={{ borderRadius: '4px', padding: '4px 8px' }}>
                    {t('animalDetail.confinement')}
                  </Tag>
                )}
                {animal.fertilization && (
                  <Tag color="purple" style={{ borderRadius: '4px', padding: '4px 8px' }}>
                    {t('animalDetail.fertilization')}
                  </Tag>
                )}
                {animal.castrated && (
                  <Tag color="red" style={{ borderRadius: '4px', padding: '4px 8px' }}>
                    {t('animalDetail.castrated')}
                  </Tag>
                )}
              </Space>
            </div>
          </Card>

          <Card 
            title={
              <Title level={4} style={{ margin: 0, color: '#262626' }}>
                {t('animalDetail.parents')}
              </Title>
            }
            style={{ 
              borderRadius: '8px',
              boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
              border: 'none',
              marginTop: '24px'
            }}
            styles={{ body: { padding: '24px' } }}
          >
            <Row gutter={[24, 24]}>
              <Col xs={24} sm={12}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                  <Avatar 
                    size={48} 
                    icon={<UserOutlined />}
                    style={{ backgroundColor: '#1890ff' }}
                  />
                  <div>
                    <Text strong style={{ color: '#262626', fontSize: '16px' }}>
                      {animal.father?.animal_name || t('common.notInformed')}
                    </Text>
                    <div style={{ color: '#8c8c8c', fontSize: '14px' }}>
                      {t('animalDetail.father')}
                    </div>
                  </div>
                </div>
              </Col>
              <Col xs={24} sm={12}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                  <Avatar 
                    size={48} 
                    icon={<UserOutlined />}
                    style={{ backgroundColor: '#f759ab' }}
                  />
                  <div>
                    <Text strong style={{ color: '#262626', fontSize: '16px' }}>
                      {animal.mother?.animal_name || t('common.notInformed')}
                    </Text>
                    <div style={{ color: '#8c8c8c', fontSize: '14px' }}>
                      {t('animalDetail.mother')}
                    </div>
                  </div>
                </div>
              </Col>
            </Row>
          </Card>
        </Col>

        <Col xs={24} lg={8}>
          <Card 
            style={{ 
              borderRadius: '8px',
              boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
              border: 'none',
              textAlign: 'center'
            }}
            styles={{ body: { padding: '24px' } }}
          >
            {animal.photo ? (
              <Image
                src={animal.photo}
                alt={t('animalDetail.photoAlt', { name: animal.animal_name })}
                style={{ 
                  width: '100%', 
                  height: '300px', 
                  objectFit: 'cover',
                  borderRadius: '8px'
                }}
                placeholder={
                  <div style={{ 
                    width: '100%', 
                    height: '300px', 
                    display: 'flex', 
                    alignItems: 'center', 
                    justifyContent: 'center',
                    backgroundColor: '#f5f5f5',
                    borderRadius: '8px'
                  }}>
                    <CameraOutlined style={{ fontSize: '48px', color: '#ccc' }} />
                  </div>
                }
              />
            ) : (
              <div style={{ 
                width: '100%', 
                height: '300px', 
                display: 'flex', 
                alignItems: 'center', 
                justifyContent: 'center',
                backgroundColor: '#f5f5f5',
                border: '2px dashed #d9d9d9',
                borderRadius: '8px'
              }}>
                <div style={{ textAlign: 'center' }}>
                  <CameraOutlined style={{ fontSize: '48px', color: '#ccc' }} />
                  <p style={{ marginTop: '16px', color: '#8c8c8c' }}>
                    {t('animalDetail.noPhoto')}
                  </p>
                </div>
              </div>
            )}
          </Card>
        </Col>
      </Row>
    </div>
  );
};
