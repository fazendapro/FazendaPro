import React, { useState } from 'react';
import { Modal, Checkbox, Button, Space, Divider } from 'antd';
import { useTranslation } from 'react-i18next';
import { useAnimalColumnBuilder } from './column-builder';

export interface ColumnOption {
  key: string;
  label: string;
  defaultVisible: boolean;
}

interface FilterModalProps {
  isOpen: boolean;
  onClose: () => void;
  onApplyFilters: (selectedColumns: string[]) => void;
  currentColumns: string[];
}

const FilterModal: React.FC<FilterModalProps> = ({ 
  isOpen, 
  onClose, 
  onApplyFilters, 
  currentColumns 
}) => {
  const { t } = useTranslation();
  const [selectedColumns, setSelectedColumns] = useState<string[]>(currentColumns);
  const { getColumnOptions, getDefaultColumnKeys } = useAnimalColumnBuilder();

  const availableColumns: ColumnOption[] = getColumnOptions();

  const handleColumnToggle = (columnKey: string, checked: boolean) => {
    if (checked) {
      setSelectedColumns(prev => [...prev, columnKey]);
    } else {
      setSelectedColumns(prev => prev.filter(col => col !== columnKey));
    }
  };

  const handleSelectAll = () => {
    setSelectedColumns(availableColumns.map(col => col.key));
  };

  const handleDeselectAll = () => {
    setSelectedColumns([]);
  };

  const handleResetToDefault = () => {
    const defaultColumns = getDefaultColumnKeys();
    setSelectedColumns(defaultColumns);
  };

  const handleApply = () => {
    onApplyFilters(selectedColumns);
    onClose();
  };

  const handleCancel = () => {
    setSelectedColumns(currentColumns);
    onClose();
  };

  return (
    <Modal
      title={t('animalTable.columnFilter')}
      open={isOpen}
      onCancel={handleCancel}
      footer={[
        <Button key="cancel" onClick={handleCancel}>
          {t('animalTable.cancel')}
        </Button>,
        <Button key="apply" type="primary" onClick={handleApply}>
          {t('animalTable.apply')}
        </Button>,
      ]}
      width={500}
    >
      <div style={{ marginBottom: '16px' }}>
        <Space>
          <Button size="small" onClick={handleSelectAll}>
            {t('animalTable.selectAll')}
          </Button>
          <Button size="small" onClick={handleDeselectAll}>
            {t('animalTable.deselectAll')}
          </Button>
          <Button size="small" onClick={handleResetToDefault}>
            {t('animalTable.resetToDefault')}
          </Button>
        </Space>
      </div>

      <Divider />

      <div style={{ maxHeight: '400px', overflowY: 'auto' }}>
        <Space direction="vertical" style={{ width: '100%' }}>
          {availableColumns.map((column) => (
            <Checkbox
              key={column.key}
              checked={selectedColumns.includes(column.key)}
              onChange={(e) => handleColumnToggle(column.key, e.target.checked)}
            >
              {column.label}
            </Checkbox>
          ))}
        </Space>
      </div>

      <Divider />

      <div style={{ fontSize: '12px', color: '#666' }}>
        <p>{t('animalTable.selectedColumns')}: {selectedColumns.length}</p>
      </div>
    </Modal>
  );
};

export { FilterModal }; 