import { forwardRef, useImperativeHandle, useState, useMemo } from 'react'
import { Table, Spin, Alert, Button, Space, DatePicker, Card } from 'antd'
import { PlusOutlined, EditOutlined, ClearOutlined, UnorderedListOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { useVaccineApplication } from '../../hooks/useVaccineApplication'
import { useFarm } from '../../../../../hooks/useFarm'
import { useResponsive } from '../../../../../hooks'
import { VaccineApplication, VaccineApplicationFilters } from '../../domain/model/vaccine-application'
import { CustomPagination } from '../../../../../components/lib/Pagination/custom-pagination'
import dayjs from '../../../../../config/dayjs'

const { RangePicker } = DatePicker
interface VaccinesTableRef {
  refetch: () => void
}

interface VaccinesTableProps {
  onAddVaccine?: () => void
  onAddApplication?: () => void
  onEditApplication?: (application: VaccineApplication) => void
  onListVaccines?: () => void
}

const VaccinesTable = forwardRef<VaccinesTableRef, VaccinesTableProps>((props, ref) => {
  const { onAddVaccine, onAddApplication, onEditApplication, onListVaccines } = props
  const { t } = useTranslation()
  const { farm } = useFarm()
  const { isMobile, isTablet } = useResponsive()
  const [filters, setFilters] = useState<VaccineApplicationFilters>({})
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  
  const paginationFilters = useMemo(() => ({
    ...filters,
    page: currentPage,
    limit: pageSize
  }), [filters, currentPage, pageSize])
  
  const { vaccineApplications, total, loading, error, refetch } = useVaccineApplication(farm?.id || 0, paginationFilters)

  useImperativeHandle(ref, () => ({
    refetch
  }))

  const columns = [
    {
      title: t('animalTable.vaccines.animalName'),
      dataIndex: ['animal', 'animal_name'],
      key: 'animalName',
      sorter: (a: VaccineApplication, b: VaccineApplication) => 
        a.animal.animal_name.localeCompare(b.animal.animal_name)
    },
    {
      title: t('animalTable.vaccines.earTag'),
      dataIndex: ['animal', 'ear_tag_number_local'],
      key: 'earTag',
      sorter: (a: VaccineApplication, b: VaccineApplication) => 
        a.animal.ear_tag_number_local - b.animal.ear_tag_number_local
    },
    {
      title: t('animalTable.vaccines.vaccineName'),
      dataIndex: ['vaccine', 'name'],
      key: 'vaccineName',
      sorter: (a: VaccineApplication, b: VaccineApplication) => 
        a.vaccine.name.localeCompare(b.vaccine.name)
    },
    {
      title: t('animalTable.vaccines.applicationDate'),
      dataIndex: 'application_date',
      key: 'applicationDate',
      sorter: (a: VaccineApplication, b: VaccineApplication) => 
        new Date(a.application_date).getTime() - new Date(b.application_date).getTime(),
      render: (date: string) => new Date(date).toLocaleDateString('pt-BR')
    },
    {
      title: t('animalTable.vaccines.batchNumber'),
      dataIndex: 'batch_number',
      key: 'batchNumber',
      render: (batch: string) => batch || '-'
    },
    {
      title: t('animalTable.vaccines.veterinarian'),
      dataIndex: 'veterinarian',
      key: 'veterinarian',
      render: (vet: string) => vet || '-'
    },
    {
      title: t('animalTable.vaccines.actions'),
      key: 'actions',
      width: isMobile ? 80 : 120,
      render: (record: VaccineApplication) => (
        <Button
          type="primary"
          size={isMobile ? 'small' : 'small'}
          icon={<EditOutlined />}
          onClick={() => onEditApplication?.(record)}
          title={t('animalTable.vaccines.editApplication')}
        >
          {!isMobile && t('animalTable.vaccines.edit')}
        </Button>
      )
    }
  ]

  const handleDateRangeChange = (dates: [dayjs.Dayjs | null, dayjs.Dayjs | null] | null) => {
    if (dates && dates[0] && dates[1]) {
      setFilters(prev => ({
        ...prev,
        startDate: dates[0]!.format('YYYY-MM-DD'),
        endDate: dates[1]!.format('YYYY-MM-DD')
      }))
      setCurrentPage(1)
    } else {
      setFilters(prev => ({
        ...prev,
        startDate: undefined,
        endDate: undefined
      }))
      setCurrentPage(1)
    }
  }

  const handleClearFilters = () => {
    setFilters({})
    setCurrentPage(1)
  }

  const handlePageChange = (page: number, size: number) => {
    setCurrentPage(page)
    setPageSize(size)
  }

  const handleShowSizeChange = (_: number, size: number) => {
    setCurrentPage(1)
    setPageSize(size)
  }

  const rangePickerValue: [dayjs.Dayjs, dayjs.Dayjs] | null = filters.startDate && filters.endDate 
    ? [dayjs(filters.startDate), dayjs(filters.endDate)] 
    : null

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
      <Card style={{ marginBottom: '16px' }}>
        <div style={{ 
          display: 'flex', 
          justifyContent: 'space-between', 
          alignItems: 'center',
          flexDirection: isMobile ? 'column' : 'row',
          gap: isMobile ? '12px' : '0'
        }}>
          <Space direction={isMobile ? 'vertical' : 'horizontal'} style={{ width: isMobile ? '100%' : 'auto' }}>
            <RangePicker
              value={rangePickerValue}
              onChange={handleDateRangeChange}
              placeholder={[t('animalTable.vaccines.filters.startDate'), t('animalTable.vaccines.filters.endDate')]}
              style={{ width: isMobile ? '100%' : 'auto' }}
              size={'middle'}
              format="DD/MM/YYYY"
              allowClear
              showToday
            />
          </Space>

          <Space direction={isMobile ? 'vertical' : 'horizontal'} style={{ width: isMobile ? '100%' : 'auto' }}>
            <Button
              type="default"
              icon={<UnorderedListOutlined />}
              onClick={onListVaccines}
              size={'middle'}
              block={isMobile}
            >
              {t('animalTable.vaccines.listVaccines') || 'Listar Vacinas'}
            </Button>
            <Button
              type="default"
              icon={<PlusOutlined />}
              onClick={onAddVaccine}
              size={'middle'}
              block={isMobile}
            >
              {t('animalTable.vaccines.addVaccine')}
            </Button>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={onAddApplication}
              size={'middle'}
              block={isMobile}
            >
              {t('animalTable.vaccines.addApplication')}
            </Button>
            
            {(filters.startDate || filters.endDate) && (
              <Button
                icon={<ClearOutlined />}
                onClick={handleClearFilters}
                size={'middle'}
                block={isMobile}
              >
                {t('animalTable.vaccines.clearFilters')}
              </Button>
            )}
          </Space>
        </div>
      </Card>

      <Table
        columns={columns}
        dataSource={vaccineApplications}
        rowKey="id"
        pagination={false}
        scroll={{
          x: isMobile ? 800 : isTablet ? 900 : 1000,
          y: isMobile ? 400 : undefined
        }}
        size={isMobile ? 'small' : 'middle'}
        style={{
          fontSize: isMobile ? '12px' : '14px'
        }}
      />

      <CustomPagination
        current={currentPage}
        total={total}
        pageSize={pageSize}
        onChange={handlePageChange}
        onShowSizeChange={handleShowSizeChange}
        showSizeChanger={!isMobile}
        showTotal={!isMobile}
      />
    </div>
  )
})

VaccinesTable.displayName = 'VaccinesTable'

export { VaccinesTable }

