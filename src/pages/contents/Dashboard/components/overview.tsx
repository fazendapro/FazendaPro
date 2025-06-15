import React from 'react';
import { Card, Row, Col } from 'antd';
import { DollarOutlined, RiseOutlined, BarChartOutlined, HomeOutlined } from '@ant-design/icons';
import { InfoCard } from '../../../../components';
import { useTranslation } from 'react-i18next';

const Overview: React.FC = () => {
  const { t } = useTranslation();
  return (
  <Card title={t('dashboard.salesOverview')} style={{ marginBottom: 16, borderRadius: 8, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
    <Row>
      <Col span={6}>
        <InfoCard title={t('dashboard.sales')} value="R$ 832" icon={<DollarOutlined style={{ fontSize: 24, color: '#1890ff' }} />} isLast={false} />
      </Col>
      <Col span={6}>
        <InfoCard title={t('dashboard.revenue')} value="R$ 18,300" icon={<RiseOutlined style={{ fontSize: 24, color: '#1890ff' }} />} isLast={false} />
      </Col>
      <Col span={6}>
        <InfoCard title={t('dashboard.profit')} value="R$ 868" icon={<BarChartOutlined style={{ fontSize: 24, color: '#faad14' }} />} isLast={false} />
      </Col>
      <Col span={6}>
        <InfoCard title={t('dashboard.cost')} value="R$ 17,432" icon={<HomeOutlined style={{ fontSize: 24, color: '#52c41a' }} />} isLast={true} />
      </Col>
    </Row>
  </Card>
  );
};

export {Overview};