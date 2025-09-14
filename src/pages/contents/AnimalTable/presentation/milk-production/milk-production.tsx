import { useRef, useState } from 'react'
import { MilkProductionTable } from './milk-production-table'
import { CreateMilkProductionModal } from './create-milk-production-modal'
import { MilkProduction as MilkProductionType } from '../../domain/model/milk-production'

interface MilkProductionRef {
  refetch: () => void
}

const MilkProduction = () => {
  const tableRef = useRef<MilkProductionRef>(null)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [preselectedAnimalId, setPreselectedAnimalId] = useState<number | undefined>()
  const [editingProduction, setEditingProduction] = useState<MilkProductionType | undefined>()

  const handleAddProduction = () => {
    setPreselectedAnimalId(undefined)
    setEditingProduction(undefined)
    setIsModalVisible(true)
  }

  const handleEditProduction = (production: MilkProductionType) => {
    setEditingProduction(production)
    setPreselectedAnimalId(undefined)
    setIsModalVisible(true)
  }

  const handleModalCancel = () => {
    setIsModalVisible(false)
    setPreselectedAnimalId(undefined)
    setEditingProduction(undefined)
  }

  const handleModalSuccess = () => {
    setIsModalVisible(false)
    setPreselectedAnimalId(undefined)
    setEditingProduction(undefined)
    if (tableRef.current) {
      tableRef.current.refetch()
    }
  }

  return (
    <div>
      <MilkProductionTable
        ref={tableRef}
        onAddProduction={handleAddProduction}
        onEditProduction={handleEditProduction}
      />
      
      <CreateMilkProductionModal
        visible={isModalVisible}
        onCancel={handleModalCancel}
        onSuccess={handleModalSuccess}
        preselectedAnimalId={preselectedAnimalId}
        editingProduction={editingProduction}
      />
    </div>
  )
}

export { MilkProduction }
