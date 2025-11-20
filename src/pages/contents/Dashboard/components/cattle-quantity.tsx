import React from 'react';
import { Card, Row, Col } from 'antd';
import { UserOutlined, UserSwitchOutlined } from '@ant-design/icons';
import { DashboardInfoCard } from '../../../../components';
import { useTranslation } from 'react-i18next';

const CattleQuantity: React.FC = () => {
  const { t } = useTranslation();

  return (
    <Card title={t('cattleQuantity')} style={{ marginBottom: 16, borderRadius: 8 }}>
      <Row>
        <Col span={12}>
          <DashboardInfoCard
            title={t('dashboard.machos')}
            value="868"
            icon={<UserOutlined style={{ fontSize: 24, color: '#faad14' }} />} 
            isLast={false}
          />
        </Col>
        <Col span={12}>
          <DashboardInfoCard
            title={t('dashboard.femeas')}
            value="200"
            icon={<UserSwitchOutlined style={{ fontSize: 24, color: '#d3adf7' }} />}
            isLast={true}
          />
        </Col>
      </Row>
    </Card>
  );
};

export { CattleQuantity };