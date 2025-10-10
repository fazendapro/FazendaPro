import { forwardRef, useImperativeHandle, useState } from 'react'
import { Table, Spin, Alert, Button, Space, DatePicker, Select } from 'antd'
import { PlusOutlined, EditOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { useMilkProduction } from '../../hooks/useMilkProduction'
import { useFarm } from '../../../../../hooks/useFarm'
import { useResponsive } from '../../../../../hooks'
import { MilkProduction, MilkProductionFilters } from '../../domain/model/milk-production'
import { CustomPagination } from '../../../../../components/lib/Pagination/custom-pagination'
import dayjs from 'dayjs'

const { RangePicker } = DatePicker
const { Option } = Select

interface MilkProductionTableRef {
  refetch: () => void
}

interface MilkProductionTableProps {
  onAddProduction?: () => void
  onEditProduction?: (production: MilkProduction) => void
}

const MilkProductionTable = forwardRef<MilkProductionTableRef, MilkProductionTableProps>((props, ref) => {
  const { onAddProduction, onEditProduction } = props
  const { t } = useTranslation()
  const { farm } = useFarm()
  const { isMobile, isTablet } = useResponsive()
  const [filters, setFilters] = useState<MilkProductionFilters>({ period: 'all' })
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize, setPageSize] = useState(isMobile ? 5 : isTablet ? 8 : 10)
  
  const { milkProductions, loading, error, refetch } = useMilkProduction(farm?.id || 0, filters)

  useImperativeHandle(ref, () => ({
    refetch
  }))

  const columns = [
    {
      title: t('animalTable.milkProductionContainer.animalName'),
      dataIndex: ['animal', 'animal_name'],
      key: 'animalName',
      sorter: (a: MilkProduction, b: MilkProduction) => 
        a.animal.animal_name.localeCompare(b.animal.animal_name)
    },
    {
      title: t('animalTable.milkProductionContainer.earTag'),
      dataIndex: ['animal', 'ear_tag_number_local'],
      key: 'earTag',
      sorter: (a: MilkProduction, b: MilkProduction) => 
        a.animal.ear_tag_number_local - b.animal.ear_tag_number_local
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
    },
    {
      title: t('animalTable.milkProductionContainer.actions'),
      key: 'actions',
      width: isMobile ? 80 : 120,
      render: (record: MilkProduction) => (
        <Button
          type="primary"
          size={isMobile ? 'small' : 'small'}
          icon={<EditOutlined />}
          onClick={() => onEditProduction?.(record)}
          title={t('animalTable.milkProductionContainer.editProduction')}
        >
          {!isMobile && t('animalTable.milkProductionContainer.edit')}
        </Button>
      )
    }
  ]

  const handlePeriodChange = (value: 'week' | 'month' | 'all') => {
    setFilters(prev => ({ ...prev, period: value }))
  }

  const handleDateRangeChange = (dates: [dayjs.Dayjs | null, dayjs.Dayjs | null] | null) => {
    if (dates && dates.length === 2 && dates[0] && dates[1]) {
      setFilters(prev => ({
        ...prev,
        startDate: dates[0]!.format('YYYY-MM-DD'),
        endDate: dates[1]!.format('YYYY-MM-DD')
      }))
    } else {
      setFilters(prev => ({
        ...prev,
        startDate: undefined,
        endDate: undefined
      }))
    }
  }

  const handlePageChange = (page: number, size: number) => {
    setCurrentPage(page)
    setPageSize(size)
  }

  const handleShowSizeChange = (_: number, size: number) => {
    setCurrentPage(1)
    setPageSize(size)
  }

  const startIndex = (currentPage - 1) * pageSize
  const endIndex = startIndex + pageSize
  const paginatedData = (milkProductions || []).slice(startIndex, endIndex)

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
      <div style={{ 
        marginBottom: '16px', 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        flexDirection: isMobile ? 'column' : 'row',
        gap: isMobile ? '12px' : '0'
      }}>
        <Space direction={isMobile ? 'vertical' : 'horizontal'} style={{ width: isMobile ? '100%' : 'auto' }}>
          <Select
            value={filters.period}
            onChange={handlePeriodChange}
            style={{ width: isMobile ? '100%' : 120 }}
            size={'middle'}
          >
            <Option value="all">{t('animalTable.milkProductionContainer.filters.all')}</Option>
            <Option value="week">{t('animalTable.milkProductionContainer.filters.week')}</Option>
            <Option value="month">{t('animalTable.milkProductionContainer.filters.month')}</Option>
          </Select>
          
          <RangePicker
            onChange={handleDateRangeChange}
            placeholder={[t('animalTable.milkProductionContainer.filters.startDate'), t('animalTable.milkProductionContainer.filters.endDate')]}
            style={{ width: isMobile ? '100%' : 'auto' }}
            size={'middle'}
          />
        </Space>

        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={onAddProduction}
          size={'middle'}
          block={isMobile}
        >
          {t('animalTable.milkProductionContainer.addProduction')}
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={paginatedData}
        rowKey="id"
        pagination={false}
        scroll={{
          x: isMobile ? 600 : isTablet ? 700 : 800,
          y: isMobile ? 400 : undefined
        }}
        size={isMobile ? 'small' : 'middle'}
        style={{
          fontSize: isMobile ? '12px' : '14px'
        }}
      />

      <CustomPagination
        current={currentPage}
        total={(milkProductions || []).length}
        pageSize={pageSize}
        onChange={handlePageChange}
        onShowSizeChange={handleShowSizeChange}
        showSizeChanger={!isMobile}
        showTotal={!isMobile}
      />
    </div>
  )
})

MilkProductionTable.displayName = 'MilkProductionTable'

export { MilkProductionTable }
