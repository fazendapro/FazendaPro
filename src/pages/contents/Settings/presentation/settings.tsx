import React, { useState, useEffect, useCallback } from 'react';
import { Card, Form, Input, Button, Upload, message, Row, Col, Avatar, Typography, Select } from 'antd';
import { UploadOutlined, UserOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import i18n from '../../../../locale/i18n';
import dayjs from '../../../../config/dayjs';
import { UpdateFarmFactory, GetFarmFactory } from '../factories';
import { FarmData, BackendFarmData, UpdateFarmParams } from '../types/farm-types';
import { useSelectedFarm } from '../../../../hooks/useSelectedFarm';
import { useFarm } from '../../../../hooks/useFarm';

const { Title } = Typography;

const Settings: React.FC = () => {
  const { t } = useTranslation();
  const { farmId, farmLogo } = useSelectedFarm();
  const { farm } = useFarm();
  const [form] = Form.useForm();
  const [farmData, setFarmData] = useState<FarmData | null>(null);

  const loadFarmData = useCallback(async () => {
    if (!farmId) return;
    
    try {
      const getFarmUseCase = GetFarmFactory.create();
      const response = await getFarmUseCase.get(farmId);

      if (response.data) {
        const backendData = response.data as BackendFarmData;

        const newFarmData = {
          id: backendData.ID,
          logo: farmLogo || backendData.Logo || '',
          language: backendData.Language || 'pt',
          company_id: backendData.CompanyID,
          company: backendData.Company ? {
            id: backendData.Company.ID,
            company_name: backendData.Company.CompanyName,
            location: backendData.Company.Location,
            farm_cnpj: backendData.Company.FarmCNPJ,
          } : {
            id: 1,
            company_name: farm?.name || 'Fazenda',
            location: '',
            farm_cnpj: '',
          },
          created_at: backendData.CreatedAt,
          updated_at: backendData.UpdatedAt,
        };

        setFarmData(newFarmData);

        const farmLanguage = backendData.Language || 'pt';
        form.setFieldsValue({
          name: backendData.Company?.CompanyName || farm?.name || 'Fazenda',
          language: farmLanguage,
        });

        // Aplicar idioma no i18n se for diferente do atual
        if (i18n.language !== farmLanguage) {
          i18n.changeLanguage(farmLanguage);
        }

        // Aplicar locale do dayjs
        const dayjsLocaleMap: Record<string, string> = {
          'pt': 'pt-br',
          'en': 'en',
          'es': 'es',
        };
        const dayjsLocale = dayjsLocaleMap[farmLanguage] || 'pt-br';
        dayjs.locale(dayjsLocale);
      }
    } catch (error) {
      console.error('Error loading farm data:', error);
      setFarmData({
        id: farmId,
        logo: '',
        language: 'pt',
        company_id: 1,
        company: {
          id: 1,
          company_name: farm?.name || 'Fazenda',
          location: '',
          farm_cnpj: '',
        },
        created_at: farm?.created_at || new Date().toISOString(),
        updated_at: farm?.updated_at || new Date().toISOString(),
      });
      form.setFieldsValue({
          name: farm?.name || 'Fazenda',
          language: 'pt',
      });
    }
  }, [farmId, farm, form, farmLogo]);

  useEffect(() => {
    loadFarmData();
  }, [farmId, loadFarmData]);

  const handleSubmit = async () => {
    if (!farmId) return;

    try {
      const updateFarmUseCase = UpdateFarmFactory.create();
      const formValues = form.getFieldsValue();
      const params: UpdateFarmParams = {
        logo: farmData?.logo || '',
        language: formValues.language || farmData?.language || 'pt',
      };

      await updateFarmUseCase.update(farmId, params);

      if (farmData) {
        setFarmData({
          ...farmData,
          language: params.language || 'pt',
        });
      }

      message.success(t('farmUpdatedSuccessfully'));
    } catch (error) {
      message.error(t('errorUpdatingFarm'));
      console.error('Error updating farm:', error);
    }
  };

  const handleLanguageChange = async (language: string) => {
    if (!farmId) return;

    try {
      await i18n.changeLanguage(language);

      const dayjsLocaleMap: Record<string, string> = {
        'pt': 'pt-br',
        'en': 'en',
        'es': 'es',
      };
      const dayjsLocale = dayjsLocaleMap[language] || 'pt-br';
      dayjs.locale(dayjsLocale);

      form.setFieldsValue({ language });

      const updateFarmUseCase = UpdateFarmFactory.create();
      const params: UpdateFarmParams = {
        logo: farmData?.logo || '',
        language: language,
      };

      await updateFarmUseCase.update(farmId, params);

      // Atualizar farmData localmente
      if (farmData) {
        setFarmData({
          ...farmData,
          language: language,
        });
      }

      // Não recarregar os dados para evitar resetar o idioma
      message.success(t('farmUpdatedSuccessfully'));
    } catch (error) {
      console.error('Error updating language:', error);
      message.error(t('errorUpdatingFarm'));
    }
  };

  const handleLogoUpload = async (file: File) => {
    if (!farmId) return;

    const reader = new FileReader();
    reader.onload = async (e) => {
      const logoUrl = e.target?.result as string;
      
      try {
        const updateFarmUseCase = UpdateFarmFactory.create();
        const params: UpdateFarmParams = {
          logo: logoUrl,
        };
        
        await updateFarmUseCase.update(farmId, params);
        
        await loadFarmData();
        message.success(t('logoUpdatedSuccessfully'));
      } catch (error) {
        message.error(t('errorUpdatingLogo'));
        console.error('Error updating logo:', error);
      }
    };
    reader.readAsDataURL(file);
    
    return false;
  };

  const uploadProps = {
    beforeUpload: handleLogoUpload,
    showUploadList: false,
    accept: 'image/*',
  };

  if (!farmId) {
    return <div>{t('noFarmSelected')}</div>;
  }

  return (
    <div style={{ padding: '24px' }}>
      <Title level={2}>{t('title')}</Title>
      
      <Row gutter={[24, 24]}>
        <Col xs={24} lg={12}>
          <Card title={t('basicInfo')}>
            <Form
              form={form}
              layout="vertical"
              onFinish={handleSubmit}
              initialValues={{
                name: farmData?.company?.company_name || '',
              }}
            >
              <Form.Item
                label={t('farmName')}
                name="name"
              >
                <Input 
                  placeholder={t('farmNamePlaceholder')} 
                  disabled 
                  style={{ backgroundColor: '#f5f5f5' }}
                />
              </Form.Item>
              <div style={{ color: '#666', fontSize: '12px', marginTop: '-8px', marginBottom: '16px' }}>
                {t('farmNameDisabled')}
              </div>
              
              <Form.Item
                label={t('language')}
                name="language"
              >
                <Select
                  placeholder={t('languagePlaceholder')}
                  onChange={handleLanguageChange}
                  value={farmData?.language || 'pt'}
                  options={[
                    { label: t('languagePt'), value: 'pt' },
                    { label: t('languageEn'), value: 'en' },
                    { label: t('languageEs'), value: 'es' },
                  ]}
                />
              </Form.Item>
              <div style={{ color: '#666', fontSize: '12px', marginTop: '-8px', marginBottom: '16px' }}>
                {t('languageDescription')}
              </div>
            </Form>
          </Card>
        </Col>

        <Col xs={24} lg={12}>
          <Card title={t('farmLogo')}>
            <div style={{ textAlign: 'center' }}>
              <div style={{ marginBottom: 16 }}>
                {farmData?.logo ? (
                  <Avatar
                    size={120}
                    src={farmData.logo}
                    shape="square"
                    style={{ margin: '0 auto' }}
                  />
                ) : (
                  <Avatar
                    size={120}
                    shape="square"
                    icon={<UserOutlined />}
                    style={{ margin: '0 auto' }}
                  />
                )}
              </div>
              
              <Upload {...uploadProps}>
                <Button icon={<UploadOutlined />}>
                  {farmData?.logo ? t('changeLogo') : t('addLogo')}
                </Button>
              </Upload>
              
              <div style={{ marginTop: 8, color: '#666', fontSize: '12px' }}>
                {t('logoFormats')}
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export { Settings };
