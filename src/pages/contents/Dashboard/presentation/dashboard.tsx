import { Row, Col } from 'antd';
import { DashboardMilkProduction, NextToCalve, SalesAndPurchases } from '../components';

const Dashboard = () => {
  return (
    <div>
      {/* <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} sm={24} md={16} lg={16} xl={16}>
          <Overview />
        </Col>
        <Col xs={24} sm={24} md={8} lg={8} xl={8}>
          <CattleQuantity />
        </Col>
      </Row>
      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} sm={24} md={16} lg={16} xl={16}>
          <ShoppingOverview />
        </Col>
        <Col xs={24} sm={24} md={8} lg={8} xl={8}>
          <Rations />
        </Col>
      </Row> */}
      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} sm={24} md={12} lg={16} xl={16}>
          <DashboardMilkProduction />
        </Col>
        <Col xs={24} sm={24} md={12} lg={8} xl={8}>
          <NextToCalve />
        </Col>
      </Row>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={24} md={24} lg={24} xl={24}>
          <SalesAndPurchases />
        </Col>
      </Row>
    </div>
  );
};

export { Dashboard };