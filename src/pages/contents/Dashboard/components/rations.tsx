import { Card, Row, Col } from 'antd';
import { UserOutlined, FileTextOutlined } from '@ant-design/icons';
import { DashboardInfoCard } from '../../../../components';
import { useTranslation } from 'react-i18next';

const Rations = () => {
  const { t } = useTranslation();
  return (
    <Card title={t('dashboard.rations')} style={{ marginBottom: 16, borderRadius: 8 }}>
    <Row>
      <Col span={12} style={{ border: 'none'}}>
        <DashboardInfoCard title={t('dashboard.suppliers')} value="31" icon={<UserOutlined style={{ fontSize: 24, color: '#1890ff' }} />} isLast={false} />
      </Col>
      <Col span={12} style={{ border: 'none'}}>
        <DashboardInfoCard title={t('dashboard.categories')} value="21" icon={<FileTextOutlined style={{ fontSize: 24, color: '#d3adf7' }} />} isLast={true} />
      </Col>
    </Row>
  </Card>
  );
};

export {Rations};