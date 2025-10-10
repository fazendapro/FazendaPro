import { Row, Col } from 'antd';
import { DashboardMilkProduction, NextToCalve } from '../components';

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
      </Row>
      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} sm={24} md={12} lg={12} xl={12}>
          <Card title={t('dashboard.salesAndPurchases')} extra={<span>{t('dashboard.weekly')}</span>}>
            <Bar ref={barChartRef} data={barData} options={barOptions} />
          </Card>
        </Col>
        <Col xs={24} sm={24} md={12} lg={12} xl={12}>
          <Card title={t('dashboard.semenation')}>
            <Line ref={lineChartRef} data={lineData} options={lineOptions} />
          </Card>
        </Col>
      </Row> */}
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={24} md={12} lg={12} xl={12}>
          <DashboardMilkProduction />
        </Col>
        <Col xs={24} sm={24} md={12} lg={12} xl={12}>
          <NextToCalve />
        </Col>
      </Row>
    </div>
  );
};

export { Dashboard };