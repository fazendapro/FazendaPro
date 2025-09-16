import React, { useRef, useState } from 'react'
import { ReproductionTable } from './reproduction-table'
import { CreateReproductionModal } from './create-reproduction-modal'
import { Reproduction as ReproductionType } from '../../domain/model/reproduction'

interface ReproductionRef {
  refetch: () => void
}

const Reproduction = () => {
  const tableRef = useRef<ReproductionRef>(null)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [preselectedAnimalId, setPreselectedAnimalId] = useState<number | undefined>()
  const [editingReproduction, setEditingReproduction] = useState<ReproductionType | undefined>()

  const handleAddReproduction = () => {
    setPreselectedAnimalId(undefined)
    setEditingReproduction(undefined)
    setIsModalVisible(true)
  }

  const handleEditReproduction = (reproduction: ReproductionType) => {
    setEditingReproduction(reproduction)
    setPreselectedAnimalId(undefined)
    setIsModalVisible(true)
  }

  const handleModalCancel = () => {
    setIsModalVisible(false)
    setPreselectedAnimalId(undefined)
    setEditingReproduction(undefined)
  }

  const handleModalSuccess = () => {
    setIsModalVisible(false)
    setPreselectedAnimalId(undefined)
    setEditingReproduction(undefined)
    if (tableRef.current) {
      tableRef.current.refetch()
    }
  }

  return (
    <div>
      <ReproductionTable
        ref={tableRef}
        onAddReproduction={handleAddReproduction}
        onEditReproduction={handleEditReproduction}
      />
      
      <CreateReproductionModal
        visible={isModalVisible}
        onCancel={handleModalCancel}
        onSuccess={handleModalSuccess}
        preselectedAnimalId={preselectedAnimalId}
        editingReproduction={editingReproduction}
      />
    </div>
  )
}

export { Reproduction }
