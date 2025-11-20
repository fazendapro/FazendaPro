import { Row, Col } from 'antd';
import { DashboardMilkProduction, NextToCalve, SalesAndPurchases, Overview } from '../components';

const Dashboard = () => {
  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={24} md={24} lg={24} xl={24}>
          <Overview />
        </Col>
      </Row>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={24} md={24} lg={24} xl={24}>
          <SalesAndPurchases />
        </Col>
      </Row>
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} sm={24} md={12} lg={16} xl={16}>
          <DashboardMilkProduction />
        </Col>
        <Col xs={24} sm={24} md={12} lg={8} xl={8}>
          <NextToCalve />
        </Col>
      </Row>
    </div>
  );
};

export { Dashboard };