import { ColumnType } from 'antd/es/table';
import { Animal } from '../../types/type';
import { ActionButton } from './action-button';
import { TFunction } from 'i18next';
import { useTranslation } from 'react-i18next';

export interface AnimalColumn {
  key: string;
  title: string;
  dataIndex: string;
  defaultVisible: boolean;
  render?: (value: unknown, record: Animal) => React.ReactNode;
}

export class AnimalColumnBuilder {
  private columns: AnimalColumn[] = [];
  private t: TFunction;

  constructor(t: TFunction) {
    this.t = t;
    this.initializeColumns();
  }

  private initializeColumns() {
    const t = this.t;
    this.columns = [
      {
        key: 'animal_name',
        title: t('animalTable.columns.animalName'),
        dataIndex: 'animal_name',
        defaultVisible: true
      },
      {
        key: 'actions',
        title: t('animalTable.columns.actions'),
        dataIndex: 'actions',
        defaultVisible: true,
        render: (_, record: Animal) => {
          return <ActionButton animalId={record.id} />;
        }
      },
      {
        key: 'type',
        title: t('animalTable.columns.type'),
        dataIndex: 'type',
        defaultVisible: true
      },
      {
        key: 'ear_tag_number_local',
        title: t('animalTable.columns.earTagNumberLocal'),
        dataIndex: 'ear_tag_number_local',
        defaultVisible: true
      },
      {
        key: 'ear_tag_number_register',
        title: t('animalTable.columns.earTagNumberRegister'),
        dataIndex: 'ear_tag_number_register',
        defaultVisible: true
      },
      {
        key: 'breed',
        title: t('animalTable.columns.breed'),
        dataIndex: 'breed',
        defaultVisible: true
      },
      {
        key: 'birth_date',
        title: t('animalTable.columns.birthDate'),
        dataIndex: 'birth_date',
        defaultVisible: true,
        render: (date: unknown) => {
          if (!date || typeof date !== 'string') return '-';
          return new Date(date).toLocaleDateString('pt-BR');
        }
      },
      {
        key: 'sex',
        title: t('animalTable.columns.sex'),
        dataIndex: 'sex',
        defaultVisible: true,
        render: (sex: unknown) => {
          if (typeof sex !== 'number') return '-';
          return sex === 0 ? t('animalTable.columnValues.male') : t('animalTable.columnValues.female');
        }
      },
      {
        key: 'age',
        title: t('animalTable.columns.age'),
        dataIndex: 'age',
        defaultVisible: false,
        render: (age: unknown) => {
          if (typeof age !== 'number') return '-';
          return age ? `${age} ${t('animalTable.columnValues.years')}` : '-';
        }
      },
      {
        key: 'weight',
        title: t('animalTable.columns.weight'),
        dataIndex: 'weight',
        defaultVisible: false,
        render: (weight: unknown) => {
          if (weight === null || weight === undefined) return '-';
          if (typeof weight !== 'number') return '-';
          if (weight === 0) return '-';
          return `${weight.toFixed(2)} ${t('animalTable.columnValues.kg')}`;
        }
      },
      {
        key: 'milk_production',
        title: t('animalTable.columns.milkProduction'),
        dataIndex: 'milk_production',
        defaultVisible: false,
        render: (production: unknown) => {
          if (typeof production !== 'number') return '-';
          return production ? `${production} ${t('animalTable.columnValues.litersPerDay')}` : '-';
        }
      },
      {
        key: 'last_insemination',
        title: t('animalTable.columns.lastInsemination'),
        dataIndex: 'last_insemination',
        defaultVisible: false,
        render: (date: unknown) => {
          if (!date || typeof date !== 'string') return '-';
          return new Date(date).toLocaleDateString('pt-BR');
        }
      },
      {
        key: 'pregnancy_status',
        title: t('animalTable.columns.pregnancyStatus'),
        dataIndex: 'pregnancy_status',
        defaultVisible: false,
        render: (status: unknown) => {
          if (typeof status !== 'number') return '-';
          const statusMap: Record<number, string> = {
            0: t('animalTable.columnValues.notPregnant'),
            1: t('animalTable.columnValues.pregnant'),
            2: t('animalTable.columnValues.lactating')
          };
          return statusMap[status] || '-';
        }
      },
      {
        key: 'health_status',
        title: t('animalTable.columns.healthStatus'),
        dataIndex: 'health_status',
        defaultVisible: false,
        render: (status: unknown) => {
          if (typeof status !== 'number') return '-';
          const statusMap: Record<number, string> = {
            0: t('animalTable.columnValues.healthy'),
            1: t('animalTable.columnValues.inTreatment'),
            2: t('animalTable.columnValues.sick')
          };
          return statusMap[status] || '-';
        }
      },
      {
        key: 'location',
        title: t('animalTable.columns.location'),
        dataIndex: 'location',
        defaultVisible: false
      },
      {
        key: 'owner',
        title: t('animalTable.columns.owner'),
        dataIndex: 'owner',
        defaultVisible: false
      },
      {
        key: 'castrated',
        title: t('animalTable.columns.castrated'),
        dataIndex: 'castrated',
        defaultVisible: false,
        render: (castrated: unknown) => {
          if (typeof castrated !== 'boolean') return '-';
          return castrated ? t('animalTable.columnValues.yes') : t('animalTable.columnValues.no');
        }
      },
      {
        key: 'confinement',
        title: t('animalTable.columns.confinement'),
        dataIndex: 'confinement',
        defaultVisible: false,
        render: (confinement: unknown) => {
          if (typeof confinement !== 'boolean') return '-';
          return confinement ? t('animalTable.columnValues.yes') : t('animalTable.columnValues.no');
        }
      },
      {
        key: 'fertilization',
        title: t('animalTable.columns.fertilization'),
        dataIndex: 'fertilization',
        defaultVisible: false,
        render: (fertilization: unknown) => {
          if (typeof fertilization !== 'boolean') return '-';
          return fertilization ? t('animalTable.columnValues.yes') : t('animalTable.columnValues.no');
        }
      },
      {
        key: 'status',
        title: t('animalTable.columns.status'),
        dataIndex: 'status',
        defaultVisible: false,
        render: (status: unknown) => {
          if (typeof status !== 'number') return '-';
          const statusMap: Record<number, string> = {
            0: t('animalTable.columnValues.active'),
            1: t('animalTable.columnValues.inactive'),
            2: t('animalTable.columnValues.sold'),
            3: t('animalTable.columnValues.slaughtered')
          };
          return statusMap[status] || t('animalTable.columnValues.unknown');
        }
      },
      {
        key: 'purpose',
        title: t('animalTable.columns.purpose'),
        dataIndex: 'purpose',
        defaultVisible: false,
        render: (purpose: unknown) => {
          if (typeof purpose !== 'number') return '-';
          const purposeMap: Record<number, string> = {
            0: t('animalTable.columnValues.milk'),
            1: t('animalTable.columnValues.meat'),
            2: t('animalTable.columnValues.reproduction'),
            3: t('animalTable.columnValues.mixed')
          };
          return purposeMap[purpose] || t('animalTable.columnValues.unknown');
        }
      },
      {
        key: 'current_batch',
        title: t('animalTable.columns.currentBatch'),
        dataIndex: 'current_batch',
        defaultVisible: false
      },
      {
        key: 'created_at',
        title: t('animalTable.columns.createdAt'),
        dataIndex: 'createdAt',
        defaultVisible: false,
        render: (date: unknown) => {
          if (!date || typeof date !== 'string') return '-';
          return new Date(date).toLocaleDateString('pt-BR');
        }
      }
    ];
  }

  public getAllColumns(): AnimalColumn[] {
    return this.columns;
  }

  public getDefaultColumns(): AnimalColumn[] {
    return this.columns.filter(col => col.defaultVisible);
  }

  public getDefaultColumnKeys(): string[] {
    return this.getDefaultColumns().map(col => col.key);
  }

  public getColumnsByKeys(selectedKeys: string[]): AnimalColumn[] {
    return this.columns.filter(col => selectedKeys.includes(col.key));
  }

  public getColumnOptions(): { key: string; label: string; defaultVisible: boolean }[] {
    return this.columns.map(col => ({
      key: col.key,
      label: col.title,
      defaultVisible: col.defaultVisible
    }));
  }

  public buildTableColumns(selectedKeys: string[]): ColumnType<Animal>[] {
    const selectedColumns = this.getColumnsByKeys(selectedKeys);
    
    // Define larguras específicas para colunas que precisam de controle
    const columnWidths: Record<string, number> = {
      'ear_tag_number_local': 140,
      'ear_tag_number_register': 180,
      'animal_name': 180,
      'type': 100,
      'breed': 120,
      'sex': 80,
      'birth_date': 140,
      'actions': 100
    };
    
    return selectedColumns.map(col => ({
      title: col.title,
      dataIndex: col.dataIndex,
      key: col.key,
      render: col.render,
      width: columnWidths[col.key] || undefined,
      ellipsis: col.key === 'ear_tag_number_local' || col.key === 'ear_tag_number_register' ? true : undefined
    }));
  }
}

export const useAnimalColumnBuilder = () => {
  const { t } = useTranslation();
  const builder = new AnimalColumnBuilder(t);
  
  return {
    getAllColumns: () => builder.getAllColumns(),
    getDefaultColumns: () => builder.getDefaultColumns(),
    getDefaultColumnKeys: () => builder.getDefaultColumnKeys(),
    getColumnsByKeys: (keys: string[]) => builder.getColumnsByKeys(keys),
    getColumnOptions: () => builder.getColumnOptions(),
    buildTableColumns: (selectedKeys: string[]) => builder.buildTableColumns(selectedKeys)
  };
}; 