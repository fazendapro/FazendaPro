import { useRef, useState } from 'react'
import { MilkProductionTable } from './milk-production-table'
import { CreateMilkProductionModal } from './create-milk-production-modal'

interface MilkProductionRef {
  refetch: () => void
}

const MilkProduction = () => {
  const tableRef = useRef<MilkProductionRef>(null)
  const [isModalVisible, setIsModalVisible] = useState(false)

  const handleAddProduction = () => {
    setIsModalVisible(true)
  }

  const handleModalCancel = () => {
    setIsModalVisible(false)
  }

  const handleModalSuccess = () => {
    setIsModalVisible(false)
    if (tableRef.current) {
      tableRef.current.refetch()
    }
  }

  return (
    <div>
      <MilkProductionTable
        ref={tableRef}
        onAddProduction={handleAddProduction}
      />
      
      <CreateMilkProductionModal
        visible={isModalVisible}
        onCancel={handleModalCancel}
        onSuccess={handleModalSuccess}
      />
    </div>
  )
}

export { MilkProduction }
