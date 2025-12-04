import { useRef, useState } from 'react'
import { VaccinesTable } from './vaccines-table'
import { CreateVaccineModal } from './create-vaccine-modal'
import { CreateVaccineApplicationModal } from './create-vaccine-application-modal'
import { VaccinesListModal } from './vaccines-list-modal'
import { VaccineApplication as VaccineApplicationType } from '../../domain/model/vaccine-application'

interface VaccinesRef {
  refetch: () => void
}

const Vaccines = () => {
  const tableRef = useRef<VaccinesRef>(null)
  const [isVaccineModalVisible, setIsVaccineModalVisible] = useState(false)
  const [isApplicationModalVisible, setIsApplicationModalVisible] = useState(false)
  const [isVaccinesListModalVisible, setIsVaccinesListModalVisible] = useState(false)
  const [preselectedAnimalId, setPreselectedAnimalId] = useState<number | undefined>()
  const [editingApplication, setEditingApplication] = useState<VaccineApplicationType | undefined>()

  const handleAddVaccine = () => {
    setIsVaccineModalVisible(true)
  }

  const handleAddApplication = () => {
    setPreselectedAnimalId(undefined)
    setEditingApplication(undefined)
    setIsApplicationModalVisible(true)
  }

  const handleEditApplication = (application: VaccineApplicationType) => {
    setEditingApplication(application)
    setPreselectedAnimalId(undefined)
    setIsApplicationModalVisible(true)
  }

  const handleListVaccines = () => {
    setIsVaccinesListModalVisible(true)
  }

  const handleVaccineModalCancel = () => {
    setIsVaccineModalVisible(false)
  }

  const handleVaccineModalSuccess = () => {
    setIsVaccineModalVisible(false)
  }

  const handleVaccinesListModalCancel = () => {
    setIsVaccinesListModalVisible(false)
  }

  const handleApplicationModalCancel = () => {
    setIsApplicationModalVisible(false)
    setPreselectedAnimalId(undefined)
    setEditingApplication(undefined)
  }

  const handleApplicationModalSuccess = () => {
    setIsApplicationModalVisible(false)
    setPreselectedAnimalId(undefined)
    setEditingApplication(undefined)
    if (tableRef.current) {
      tableRef.current.refetch()
    }
  }

  return (
    <div>
      <VaccinesTable
        ref={tableRef}
        onAddVaccine={handleAddVaccine}
        onAddApplication={handleAddApplication}
        onEditApplication={handleEditApplication}
        onListVaccines={handleListVaccines}
      />
      
      <CreateVaccineModal
        visible={isVaccineModalVisible}
        onCancel={handleVaccineModalCancel}
        onSuccess={handleVaccineModalSuccess}
      />

      <CreateVaccineApplicationModal
        visible={isApplicationModalVisible}
        onCancel={handleApplicationModalCancel}
        onSuccess={handleApplicationModalSuccess}
        preselectedAnimalId={preselectedAnimalId}
        editingApplication={editingApplication}
      />

      <VaccinesListModal
        visible={isVaccinesListModalVisible}
        onCancel={handleVaccinesListModalCancel}
      />
    </div>
  )
}

export { Vaccines }

