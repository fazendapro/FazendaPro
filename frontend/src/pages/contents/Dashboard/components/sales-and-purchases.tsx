import React, { useMemo, useRef } from 'react';
import { Card, Spin, Tooltip } from 'antd';
import { InfoCircleOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { Bar } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip as ChartTooltip,
  Legend,
} from 'chart.js';
import { useMonthlySalesAndPurchases } from '../hooks/useMonthlySalesAndPurchases';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  ChartTooltip,
  Legend
);

const SalesAndPurchases: React.FC = () => {
  const { t } = useTranslation();
  const { data, loading, error } = useMonthlySalesAndPurchases(undefined, 12);
  const barChartRef = useRef<ChartJS<'bar'>>(null);

  const chartData = useMemo(() => {
    const labels = data.sales.map(item => item.month);
    const salesValues = data.sales.map(item => item.sales || 0);

    return {
      labels,
      datasets: [
        {
          label: t('dashboard.sales') || 'Vendas',
          data: salesValues,
          backgroundColor: 'rgba(75, 192, 192, 0.6)',
          borderColor: 'rgba(75, 192, 192, 1)',
          borderWidth: 1,
        },
      ],
    };
  }, [data, t]);

  const chartOptions = useMemo(() => ({
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'bottom' as const,
      },
      tooltip: {
        mode: 'index' as const,
        intersect: false,
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          callback: function(value: string | number) {
            const numValue = typeof value === 'string' ? parseFloat(value) : value;
            return new Intl.NumberFormat('pt-BR', {
              minimumFractionDigits: 0,
              maximumFractionDigits: 0,
            }).format(numValue);
          },
        },
      },
    },
  }), []);

  if (loading) {
    return (
      <Card 
        title={
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            {t('dashboard.salesAndPurchases')}
            <Tooltip 
              title={t('dashboard.salesAndPurchasesTooltip') || 'Gráfico mostrando o valor total de vendas realizadas nos últimos 12 meses, baseado nas vendas de animais registradas no sistema.'}
              placement="top"
            >
              <InfoCircleOutlined style={{ color: '#1890ff', cursor: 'help' }} />
            </Tooltip>
          </div>
        } 
        style={{ height: '100%' }}
      >
        <div style={{ display: 'flex', justifyContent: 'center', padding: '20px' }}>
          <Spin size="large" />
        </div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card 
        title={
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            {t('dashboard.salesAndPurchases')}
            <Tooltip 
              title={t('dashboard.salesAndPurchasesTooltip') || 'Gráfico mostrando o valor total de vendas realizadas nos últimos 12 meses, baseado nas vendas de animais registradas no sistema.'}
              placement="top"
            >
              <InfoCircleOutlined style={{ color: '#1890ff', cursor: 'help' }} />
            </Tooltip>
          </div>
        } 
        style={{ height: '100%' }}
      >
        <div style={{ textAlign: 'center', padding: '20px', color: 'red' }}>
          {error}
        </div>
      </Card>
    );
  }

  return (
    <Card
      title={
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          {t('dashboard.salesAndPurchases')}
          <Tooltip 
            title={t('dashboard.salesAndPurchasesTooltip') || 'Gráfico mostrando o valor total de vendas realizadas nos últimos 12 meses, baseado nas vendas de animais registradas no sistema.'}
            placement="top"
          >
            <InfoCircleOutlined style={{ color: '#1890ff', cursor: 'help' }} />
          </Tooltip>
        </div>
      }
      style={{ height: '100%' }}
    >
      <div style={{ height: '300px', position: 'relative' }}>
        <Bar ref={barChartRef} data={chartData} options={chartOptions} />
      </div>
    </Card>
  );
};

export { SalesAndPurchases };

