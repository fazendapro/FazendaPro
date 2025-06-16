import { useRef, useEffect } from 'react';
import { Layout, Row, Col, Card } from 'antd';
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend } from 'chart.js';
import { Bar } from 'react-chartjs-2';
import { Overview, MilkProduction, NextToCalve, Rations, CattleQuantity, ShoppingOverview } from '../components';
import { useTranslation } from 'react-i18next';

const { Content } = Layout;

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

const Dashboard = () => {
  const { t } = useTranslation();
  const barChartRef = useRef<ChartJS<"bar", number[], unknown>>(null);
  const lineChartRef = useRef<ChartJS<"bar", number[], unknown>>(null);

  const barData = {
    labels: ['Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun', 'Jul', 'Ago'],
    datasets: [
      { label: 'Compras', data: [50000, 60000, 30000, 40000, 50000, 40000, 50000, 30000], backgroundColor: '#6B7280' },
      { label: 'Vendas', data: [40000, 50000, 20000, 30000, 40000, 30000, 40000, 20000], backgroundColor: '#10B981' },
    ],
  };

  const lineData = {
    labels: ['Jan', 'Fev', 'Mar', 'Abr', 'Mai'],
    datasets: [
      { label: 'Semeadas', data: [3000, 3500, 2000, 2500, 4000], borderColor: '#F59E0B', fill: false },
      { label: 'Não semeadas', data: [2000, 2500, 1500, 2000, 3000], borderColor: '#93C5FD', fill: false },
    ],
  };

  useEffect(() => {
    return () => {
      if (barChartRef.current) barChartRef.current.destroy();
      if (lineChartRef.current) lineChartRef.current.destroy();
    };
  }, []);

  return (
    <Layout style={{ padding: '24px' }}>
      <Content>
        <Row gutter={16}>
          <Col span={16}><Overview /></Col>
          <Col span={8}><CattleQuantity /></Col>
        </Row>
        <Row gutter={16}>
          <Col span={16}><ShoppingOverview /></Col>
          <Col span={8}><Rations /></Col>
        </Row>
        <Row gutter={16}>
          <Col span={12}>
            <Card title={t('dashboard.salesAndPurchases')} extra={<span>{t('dashboard.weekly')}</span>}>
              <Bar ref={barChartRef} data={barData} />
            </Card>
          </Col>
          <Col span={12}>
            <Card title={t('dashboard.semenation')}>
              <Bar ref={lineChartRef} data={lineData} options={{ indexAxis: 'y' }} />
            </Card>
          </Col>
        </Row>
        <Row gutter={16} style={{ marginTop: 16 }}>
          <Col span={12}><MilkProduction /></Col>
          <Col span={12}><NextToCalve /></Col>
        </Row>
      </Content>
    </Layout>
  );
};

export {Dashboard};