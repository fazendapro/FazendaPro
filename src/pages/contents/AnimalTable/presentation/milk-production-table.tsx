import { forwardRef, useImperativeHandle, useState } from 'react'
import { Table, Spin, Alert, Button, Space, DatePicker, Select } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { useMilkProduction } from '../hooks/useMilkProduction'
import { useFarm } from '../../../../hooks/useFarm'
import { MilkProduction } from '../types/milk-production'
import { MilkProductionFilters } from '../types/milk-production'

const { RangePicker } = DatePicker
const { Option } = Select

interface MilkProductionTableRef {
  refetch: () => void
}

interface MilkProductionTableProps {
  onAddProduction?: () => void
}

const MilkProductionTable = forwardRef<MilkProductionTableRef, MilkProductionTableProps>((props, ref) => {
  const { onAddProduction } = props
  const { t } = useTranslation()
  const { farm } = useFarm()
  const [filters, setFilters] = useState<MilkProductionFilters>({ period: 'all' })
  
  const { milkProductions, loading, error, refetch } = useMilkProduction(farm.id, filters)

  useImperativeHandle(ref, () => ({
    refetch
  }))

  const columns = [
    {
      title: t('animalTable.milkProductionContainer.animalName'),
      dataIndex: ['animal', 'name'],
      key: 'animalName',
      sorter: (a: MilkProduction, b: MilkProduction) => 
        a.animal.name.localeCompare(b.animal.name)
    },
    {
      title: t('animalTable.milkProductionContainer.earTag'),
      dataIndex: ['animal', 'earTagNumberLocal'],
      key: 'earTag',
      sorter: (a: MilkProduction, b: MilkProduction) => 
        a.animal.earTagNumberLocal.localeCompare(b.animal.earTagNumberLocal)
    },
    {
      title: t('animalTable.milkProductionContainer.liters'),
      dataIndex: 'liters',
      key: 'liters',
      sorter: (a: MilkProduction, b: MilkProduction) => a.liters - b.liters,
      render: (liters: number) => `${liters.toFixed(1)} L`
    },
    {
      title: t('animalTable.milkProductionContainer.date'),
      dataIndex: 'date',
      key: 'date',
      sorter: (a: MilkProduction, b: MilkProduction) => 
        new Date(a.date).getTime() - new Date(b.date).getTime(),
      render: (date: string) => new Date(date).toLocaleDateString('pt-BR')
    }
  ]

  const handlePeriodChange = (value: 'week' | 'month' | 'all') => {
    setFilters(prev => ({ ...prev, period: value }))
  }

  const handleDateRangeChange = (dates: any) => {
    if (dates && dates.length === 2) {
      setFilters(prev => ({
        ...prev,
        startDate: dates[0].format('YYYY-MM-DD'),
        endDate: dates[1].format('YYYY-MM-DD')
      }))
    } else {
      setFilters(prev => ({
        ...prev,
        startDate: undefined,
        endDate: undefined
      }))
    }
  }

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    )
  }

  if (error) {
    return (
      <Alert
        message="Erro"
        description={error}
        type="error"
        showIcon
        style={{ marginBottom: '16px' }}
      />
    )
  }

  return (
    <div>
      <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Space>
          <Select
            value={filters.period}
            onChange={handlePeriodChange}
            style={{ width: 120 }}
          >
            <Option value="all">{t('animalTable.milkProductionContainer.filters.all')}</Option>
            <Option value="week">{t('animalTable.milkProductionContainer.filters.week')}</Option>
            <Option value="month">{t('animalTable.milkProductionContainer.filters.month')}</Option>
          </Select>
          
          <RangePicker
            onChange={handleDateRangeChange}
            placeholder={[t('animalTable.milkProductionContainer.filters.startDate'), t('animalTable.milkProductionContainer.filters.endDate')]}
          />
        </Space>

        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={onAddProduction}
        >
          {t('animalTable.milkProductionContainer.addProduction')}
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={milkProductions}
        rowKey="id"
        pagination={{
          pageSize: 10,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total, range) =>
            `${range[0]}-${range[1]} de ${total} registros`
        }}
      />
    </div>
  )
})

MilkProductionTable.displayName = 'MilkProductionTable'

export { MilkProductionTable }
