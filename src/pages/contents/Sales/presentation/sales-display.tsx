import React, { useState, useEffect } from 'react';
import { 
  Card, 
  Table, 
  Button, 
  Space, 
  Typography, 
  Row, 
  Col, 
  DatePicker, 
  Modal,
  message,
  Tooltip,
  Tag
} from 'antd';
import { 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined, 
  SearchOutlined,
  FilterOutlined,
  ExportOutlined
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useSaleList } from '../../../../hooks/useSale';
import { SaleModal } from '../../../../components/SaleModal/SaleModal';
import { Sale, SaleFilters } from '../../../../types/sale';

const { Title, Text } = Typography;
const { RangePicker } = DatePicker;

export const SalesDisplay: React.FC = () => {
  const { t } = useTranslation();
  const {
    sales,
    loading,
    error,
    getSalesHistory,
    getSalesByDateRange,
    deleteSale,
    clearError
  } = useSaleList();

  const [modalVisible, setModalVisible] = useState(false);
  const [editingSale, setEditingSale] = useState<Sale | null>(null);
  const [filters, setFilters] = useState<SaleFilters>({});
  const [filteredSales, setFilteredSales] = useState<Sale[]>([]);

  useEffect(() => {
    getSalesHistory();
  }, []);

  useEffect(() => {
    setFilteredSales(sales);
  }, [sales]);

  useEffect(() => {
    if (sales.length > 0) {
      console.log('Sales data:', sales);
      console.log('First sale animal:', sales[0]?.animal);
    }
  }, [sales]);

  useEffect(() => {
    if (error) {
      message.error(error);
      clearError();
    }
  }, [error, clearError]);

  const handleCreateSale = () => {
    setEditingSale(null);
    setModalVisible(true);
  };

  const handleEditSale = (sale: Sale) => {
    setEditingSale(sale);
    setModalVisible(true);
  };

  const handleDeleteSale = (sale: Sale) => {
    Modal.confirm({
      title: t('sales.confirmDelete.title'),
      content: t('sales.confirmDelete.content', { buyer: sale.buyer_name }),
      okText: t('common.yes'),
      cancelText: t('common.no'),
      onOk: async () => {
        try {
          await deleteSale(sale.id);
        } catch {
        }
      },
    });
  };

  const handleModalSuccess = () => {
    setModalVisible(false);
    setEditingSale(null);
    getSalesHistory();
  };

  const handleFilter = () => {
    if (filters.start_date && filters.end_date) {
      getSalesByDateRange(filters).then(filteredData => {
        setFilteredSales(filteredData);
      });
    } else {
      setFilteredSales(sales);
    }
  };

  const handleClearFilters = () => {
    setFilters({});
    setFilteredSales(sales);
  };

  const handleExport = () => {
    // TODO: Implement export functionality
    message.info(t('sales.export.comingSoon'));
  };

  const columns = [
    {
      title: t('sales.table.animal'),
      dataIndex: ['animal', 'animal_name'],
      key: 'animal_name',
      render: (_: string, record: Sale) => {
        console.log('Animal data for record:', record.id, record.animal);
        return (
          <div>
            <Text strong>{record.animal?.animal_name || 'N/A'}</Text>
            <br />
            <Text type="secondary" style={{ fontSize: '12px' }}>
              {t('sales.table.earTag')}: {record.animal?.ear_tag_number_local || 'N/A'}
            </Text>
          </div>
        );
      },
    },
    {
      title: t('sales.table.status'),
      dataIndex: ['animal', 'status'],
      key: 'status',
      render: (status: number, record: Sale) => {
        const animalStatus = record.animal?.status ?? status;
        const statusConfig = {
          0: { text: 'Ativo', color: '#52c41a' },
          1: { text: 'Vendido', color: '#ff4d4f' },
          2: { text: 'Falecido', color: '#8c8c8c' }
        };
        const config = statusConfig[animalStatus as keyof typeof statusConfig] || { text: 'Desconhecido', color: '#8c8c8c' };
        return (
          <Tag color={config.color}>
            {config.text}
          </Tag>
        );
      },
    },
    {
      title: t('sales.table.buyer'),
      dataIndex: 'buyer_name',
      key: 'buyer_name',
    },
    {
      title: t('sales.table.price'),
      dataIndex: 'price',
      key: 'price',
      render: (price: number) => (
        <Text strong style={{ color: '#52c41a' }}>
          R$ {price.toFixed(2)}
        </Text>
      ),
    },
    {
      title: t('sales.table.saleDate'),
      dataIndex: 'sale_date',
      key: 'sale_date',
      render: (date: string) => new Date(date).toLocaleDateString('pt-BR'),
    },
    {
      title: t('sales.table.notes'),
      dataIndex: 'notes',
      key: 'notes',
      render: (notes: string) => notes || '-',
    },
    {
      title: t('common.actions'),
      key: 'actions',
      width: 120,
      fixed: 'right' as const,
      render: (_: unknown, record: Sale) => (
        <Space size="small">
          <Tooltip title={t('common.edit')}>
            <Button
              type="text"
              icon={<EditOutlined />}
              onClick={() => handleEditSale(record)}
              size="small"
            />
          </Tooltip>
          <Tooltip title={t('common.delete')}>
            <Button
              type="text"
              danger
              icon={<DeleteOutlined />}
              onClick={() => handleDeleteSale(record)}
              size="small"
            />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const totalValue = filteredSales.reduce((sum, sale) => sum + sale.price, 0);

  return (
    <div style={{ padding: '16px', backgroundColor: '#f5f5f5', minHeight: '100vh' }}>
      <div style={{ 
        marginBottom: '24px',
        backgroundColor: 'white',
        padding: '16px',
        borderRadius: '8px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
      }}>
        <Row gutter={[16, 16]} align="middle">
          <Col xs={24} sm={12}>
            <Title level={1} style={{ margin: 0, color: '#262626', fontSize: '24px' }}>
              {t('sales.title')}
            </Title>
          </Col>
          <Col xs={24} sm={12}>
            <Space direction="vertical" size="middle" style={{ width: '100%' }}>
              <Row gutter={[8, 8]}>
                <Col xs={12}>
                  <Button 
                    icon={<ExportOutlined />}
                    onClick={handleExport}
                    style={{ 
                      backgroundColor: '#f0f0f0',
                      borderColor: '#d9d9d9',
                      color: '#262626',
                      borderRadius: '6px',
                      width: '100%'
                    }}
                    size="middle"
                  >
                    {t('sales.export.title')}
                  </Button>
                </Col>
                <Col xs={12}>
                  <Button 
                    type="primary" 
                    icon={<PlusOutlined />} 
                    onClick={handleCreateSale}
                    style={{ 
                      backgroundColor: '#1890ff',
                      borderColor: '#1890ff',
                      borderRadius: '6px',
                      width: '100%'
                    }}
                    size="middle"
                  >
                    {t('sales.createSale')}
                  </Button>
                </Col>
              </Row>
            </Space>
          </Col>
        </Row>
      </div>

      <Row gutter={[24, 24]}>
        <Col xs={24} lg={18}>
          <Card 
            title={
              <Title level={4} style={{ margin: 0, color: '#262626' }}>
                {t('sales.history')}
              </Title>
            }
            style={{ 
              borderRadius: '8px',
              boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
              border: 'none'
            }}
            bodyStyle={{ padding: '16px' }}
          >
            <div style={{ marginBottom: '16px' }}>
              <Row gutter={[8, 8]}>
                <Col xs={24} sm={12} md={8}>
                  <RangePicker
                    style={{ width: '100%' }}
                    placeholder={[t('sales.filters.startDate'), t('sales.filters.endDate')]}
                    onChange={(dates) => {
                      if (dates) {
                        setFilters({
                          ...filters,
                          start_date: dates[0]?.format('YYYY-MM-DD'),
                          end_date: dates[1]?.format('YYYY-MM-DD'),
                        });
                      } else {
                        setFilters({
                          ...filters,
                          start_date: undefined,
                          end_date: undefined,
                        });
                      }
                    }}
                  />
                </Col>
                <Col xs={12} sm={6} md={4}>
                  <Button 
                    icon={<SearchOutlined />} 
                    onClick={handleFilter}
                    style={{ width: '100%' }}
                    size="middle"
                  >
                    {t('sales.filters.apply')}
                  </Button>
                </Col>
                <Col xs={12} sm={6} md={4}>
                  <Button 
                    icon={<FilterOutlined />} 
                    onClick={handleClearFilters}
                    style={{ width: '100%' }}
                    size="middle"
                  >
                    {t('sales.filters.clear')}
                  </Button>
                </Col>
              </Row>
            </div>

            <div style={{ overflowX: 'auto' }}>
              <Table
                columns={columns}
                dataSource={filteredSales}
                rowKey="id"
                loading={loading}
                scroll={{ x: 800 }}
                pagination={{
                  pageSize: 10,
                  showSizeChanger: true,
                  showQuickJumper: true,
                  showTotal: (total, range) => 
                    `${range[0]}-${range[1]} de ${total} ${t('sales.table.items')}`,
                }}
              />
            </div>
          </Card>
        </Col>

        <Col xs={24} lg={6}>
          <Card 
            title={
              <Title level={4} style={{ margin: 0, color: '#262626' }}>
                {t('sales.summary')}
              </Title>
            }
            style={{ 
              borderRadius: '8px',
              boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
              border: 'none'
            }}
            bodyStyle={{ padding: '16px' }}
          >
            <div style={{ textAlign: 'center', marginBottom: '16px' }}>
              <Text type="secondary" style={{ fontSize: '14px' }}>
                {t('sales.summaryDetails.totalSales')}
              </Text>
              <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#52c41a' }}>
                {filteredSales.length}
              </div>
            </div>

            <div style={{ textAlign: 'center' }}>
              <Text type="secondary" style={{ fontSize: '14px' }}>
                {t('sales.summaryDetails.totalValue')}
              </Text>
              <div style={{ fontSize: '20px', fontWeight: 'bold', color: '#1890ff' }}>
                R$ {totalValue.toFixed(2)}
              </div>
            </div>
          </Card>
        </Col>
      </Row>

      <SaleModal
        visible={modalVisible}
        onCancel={() => setModalVisible(false)}
        onSuccess={handleModalSuccess}
        sale={editingSale || undefined}
        mode={editingSale ? 'edit' : 'create'}
      />
    </div>
  );
};
