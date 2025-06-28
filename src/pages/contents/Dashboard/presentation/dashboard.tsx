import { useRef, useEffect } from 'react';
import { Layout, Row, Col, Card } from 'antd';
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, LineElement, PointElement } from 'chart.js';
import { Bar } from 'react-chartjs-2';
import { Line } from 'react-chartjs-2';
import { Overview, MilkProduction, NextToCalve, Rations, CattleQuantity, ShoppingOverview } from '../components';
import { useTranslation } from 'react-i18next';

const { Content } = Layout;

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, LineElement, PointElement);

const Dashboard = () => {
  const { t } = useTranslation();
  const barChartRef = useRef<ChartJS<"bar", number[], unknown>>(null);
  const lineChartRef = useRef<ChartJS<"line", number[], unknown>>(null); // Updated ref type for line chart

  const barData = {
    labels: ['Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun', 'Jul', 'Ago', 'Mai', 'Jun'],
    datasets: [
      {
        label: 'Compras',
        data: [60000, 50000, 35000, 40000, 45000, 30000, 40000, 35000, 25000, 20000],
        backgroundColor: '#A5B4FC',
      },
      {
        label: 'Vendas',
        data: [50000, 40000, 30000, 35000, 40000, 25000, 35000, 30000, 20000, 15000],
        backgroundColor: '#34D399',
      },
    ],
  };

  const barOptions = {
    plugins: {
      legend: { position: 'bottom' as const },
    },
    scales: {
      x: { title: { display: true, text: 'Meses' } },
      y: { title: { display: true, text: 'Valor' }, beginAtZero: true },
    },
    barPercentage: 0.5,
    categoryPercentage: 0.8,
    borderRadius: 10,
  };

  const lineData = {
    labels: ['Jan', 'Fev', 'Mar', 'Abr', 'Mai'],
    datasets: [
      {
        label: 'Semeadas',
        data: [4000, 3000, 3500, 2000, 4000],
        backgroundColor: '#F59E0B',
        tension: 0.4,
      },
      {
        label: 'Não semeadas',
        data: [2500, 2000, 3000, 1500, 3500],
        backgroundColor: '#93C5FD',
        tension: 0.4,
      },
    ],
  };

  const lineOptions = {
    plugins: {
      legend: { position: 'bottom' as const },
    },
    scales: {
      x: { title: { display: true, text: 'Meses' } },
      y: { title: { display: true, text: 'Valor' }, beginAtZero: true },
    },
  };

  useEffect(() => {
    if (barChartRef.current) {
      const chart = barChartRef.current;
      const gradient1 = chart.ctx.createLinearGradient(0, 0, 0, chart.chartArea.height);
      gradient1.addColorStop(0, '#A5B4FC');
      gradient1.addColorStop(1, '#93C5FD');
      const gradient2 = chart.ctx.createLinearGradient(0, 0, 0, chart.chartArea.height);
      gradient2.addColorStop(0, '#34D399');
      gradient2.addColorStop(1, '#6EE7B7');
      chart.data.datasets[0].backgroundColor = gradient1;
      chart.data.datasets[1].backgroundColor = gradient2;
      chart.update();
    }
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
              <Bar ref={barChartRef} data={barData} options={barOptions} />
            </Card>
          </Col>
          <Col span={12}>
            <Card title={t('dashboard.semenation')}>
              <Line ref={lineChartRef} data={lineData} options={lineOptions} />
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

export { Dashboard };