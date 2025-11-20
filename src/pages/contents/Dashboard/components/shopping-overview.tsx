import React from 'react';
import { Card, Row, Col } from 'antd';
import { ShoppingCartOutlined, HomeOutlined, RiseOutlined } from '@ant-design/icons';
import { DashboardInfoCard } from '../../../../components';
import { useTranslation } from 'react-i18next';

const ShoppingOverview: React.FC = () => {
  const { t } = useTranslation();
  return (
    <Card title={t('dashboard.shoppingOverview')} style={{ marginBottom: 16, borderRadius: 8 }}>
    <Row>
      <Col span={6}>
       <DashboardInfoCard title={t('dashboard.purchases')} value="82" icon={<ShoppingCartOutlined style={{ fontSize: 24, color: '#52c41a' }} />} isLast={false} />
      </Col>
      <Col span={6}>
        <DashboardInfoCard title={t('dashboard.cost')} value="R$ 13,573" icon={<HomeOutlined style={{ fontSize: 24, color: '#52c41a' }} />} isLast={false} />
      </Col>
      <Col span={6}>
        <DashboardInfoCard title={t('dashboard.slaughter')} value="5" icon={<ShoppingCartOutlined style={{ fontSize: 24, color: '#1890ff' }} />} isLast={false} />
      </Col>
      <Col span={6}>
        <DashboardInfoCard title={t('dashboard.return')} value="R$ 17,432" icon={<RiseOutlined style={{ fontSize: 24, color: '#faad14' }} />} isLast={true} />
      </Col>
    </Row>
  </Card>
  );
};

export {ShoppingOverview};