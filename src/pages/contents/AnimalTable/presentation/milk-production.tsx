import { useRef, useState } from 'react'
import { MilkProductionTable } from './milk-production-table'
import { CreateMilkProductionModal } from './create-milk-production-modal'

interface MilkProductionRef {
  refetch: () => void
}

const MilkProduction = () => {
  const tableRef = useRef<MilkProductionRef>(null)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [preselectedAnimalId, setPreselectedAnimalId] = useState<number | undefined>()

  const handleAddProduction = () => {
    setPreselectedAnimalId(undefined)
    setIsModalVisible(true)
  }

  const handleAddProductionForAnimal = (animalId: number) => {
    setPreselectedAnimalId(animalId)
    setIsModalVisible(true)
  }

  const handleModalCancel = () => {
    setIsModalVisible(false)
    setPreselectedAnimalId(undefined)
  }

  const handleModalSuccess = () => {
    setIsModalVisible(false)
    setPreselectedAnimalId(undefined)
    if (tableRef.current) {
      tableRef.current.refetch()
    }
  }

  return (
    <div>
      <MilkProductionTable
        ref={tableRef}
        onAddProduction={handleAddProduction}
        onAddProductionForAnimal={handleAddProductionForAnimal}
      />
      
      <CreateMilkProductionModal
        visible={isModalVisible}
        onCancel={handleModalCancel}
        onSuccess={handleModalSuccess}
        preselectedAnimalId={preselectedAnimalId}
      />
    </div>
  )
}

export { MilkProduction }
