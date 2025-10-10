import { ColumnType } from 'antd/es/table';
import { Animal } from '../../types/type';
import { ActionButton } from './action-button';

export interface AnimalColumn {
  key: string;
  title: string;
  dataIndex: string;
  defaultVisible: boolean;
  render?: (value: unknown, record: Animal) => React.ReactNode;
}

export class AnimalColumnBuilder {
  private columns: AnimalColumn[] = [];

  constructor() {
    this.initializeColumns();
  }

  private initializeColumns() {
    this.columns = [
      {
        key: 'animal_name',
        title: 'Nome do Animal',
        dataIndex: 'animal_name',
        defaultVisible: true
      },
      {
        key: 'type',
        title: 'Tipo',
        dataIndex: 'type',
        defaultVisible: true
      },
      {
        key: 'ear_tag_number_local',
        title: 'Número do Brinco Local',
        dataIndex: 'ear_tag_number_local',
        defaultVisible: true
      },
      {
        key: 'ear_tag_number_register',
        title: 'Número do Brinco de Registro',
        dataIndex: 'ear_tag_number_register',
        defaultVisible: true
      },
      {
        key: 'breed',
        title: 'Raça',
        dataIndex: 'breed',
        defaultVisible: true
      },
      {
        key: 'birth_date',
        title: 'Data de Nascimento',
        dataIndex: 'birth_date',
        defaultVisible: true,
        render: (date: unknown) => {
          if (!date || typeof date !== 'string') return '-';
          return new Date(date).toLocaleDateString('pt-BR');
        }
      },
      {
        key: 'sex',
        title: 'Sexo',
        dataIndex: 'sex',
        defaultVisible: true,
        render: (sex: unknown) => {
          if (typeof sex !== 'number') return '-';
          return sex === 0 ? 'Macho' : 'Fêmea';
        }
      },
      {
        key: 'age',
        title: 'Idade',
        dataIndex: 'age',
        defaultVisible: false,
        render: (age: unknown) => {
          if (typeof age !== 'number') return '-';
          return age ? `${age} anos` : '-';
        }
      },
      {
        key: 'weight',
        title: 'Peso',
        dataIndex: 'weight',
        defaultVisible: false,
        render: (weight: unknown) => {
          if (typeof weight !== 'number') return '-';
          return weight ? `${weight} kg` : '-';
        }
      },
      {
        key: 'milk_production',
        title: 'Produção de Leite',
        dataIndex: 'milk_production',
        defaultVisible: false,
        render: (production: unknown) => {
          if (typeof production !== 'number') return '-';
          return production ? `${production} L/dia` : '-';
        }
      },
      {
        key: 'last_insemination',
        title: 'Última Inseminação',
        dataIndex: 'last_insemination',
        defaultVisible: false,
        render: (date: unknown) => {
          if (!date || typeof date !== 'string') return '-';
          return new Date(date).toLocaleDateString('pt-BR');
        }
      },
      {
        key: 'pregnancy_status',
        title: 'Status de Gestação',
        dataIndex: 'pregnancy_status',
        defaultVisible: false,
        render: (status: unknown) => {
          if (typeof status !== 'number') return '-';
          const statusMap = {
            0: 'Não gestante',
            1: 'Gestante',
            2: 'Lactante'
          };
          return statusMap[status as keyof typeof statusMap] || '-';
        }
      },
      {
        key: 'health_status',
        title: 'Status de Saúde',
        dataIndex: 'health_status',
        defaultVisible: false,
        render: (status: unknown) => {
          if (typeof status !== 'number') return '-';
          const statusMap = {
            0: 'Saudável',
            1: 'Em tratamento',
            2: 'Doente'
          };
          return statusMap[status as keyof typeof statusMap] || '-';
        }
      },
      {
        key: 'location',
        title: 'Localização',
        dataIndex: 'location',
        defaultVisible: false
      },
      {
        key: 'owner',
        title: 'Proprietário',
        dataIndex: 'owner',
        defaultVisible: false
      },
      {
        key: 'castrated',
        title: 'Castrado',
        dataIndex: 'castrated',
        defaultVisible: false,
        render: (castrated: unknown) => {
          if (typeof castrated !== 'boolean') return '-';
          return castrated ? 'Sim' : 'Não';
        }
      },
      {
        key: 'confinement',
        title: 'Confinamento',
        dataIndex: 'confinement',
        defaultVisible: false,
        render: (confinement: unknown) => {
          if (typeof confinement !== 'boolean') return '-';
          return confinement ? 'Sim' : 'Não';
        }
      },
      {
        key: 'fertilization',
        title: 'Fertilização',
        dataIndex: 'fertilization',
        defaultVisible: false,
        render: (fertilization: unknown) => {
          if (typeof fertilization !== 'boolean') return '-';
          return fertilization ? 'Sim' : 'Não';
        }
      },
      {
        key: 'status',
        title: 'Status',
        dataIndex: 'status',
        defaultVisible: false,
        render: (status: unknown) => {
          if (typeof status !== 'number') return '-';
          const statusMap = {
            0: 'Ativo',
            1: 'Inativo',
            2: 'Vendido',
            3: 'Abatido'
          };
          return statusMap[status as keyof typeof statusMap] || 'Desconhecido';
        }
      },
      {
        key: 'purpose',
        title: 'Propósito',
        dataIndex: 'purpose',
        defaultVisible: false,
        render: (purpose: unknown) => {
          if (typeof purpose !== 'number') return '-';
          const purposeMap = {
            0: 'Leite',
            1: 'Carne',
            2: 'Reprodução',
            3: 'Misto'
          };
          return purposeMap[purpose as keyof typeof purposeMap] || 'Desconhecido';
        }
      },
      {
        key: 'current_batch',
        title: 'Lote Atual',
        dataIndex: 'current_batch',
        defaultVisible: false
      },
      {
        key: 'created_at',
        title: 'Criado em',
        dataIndex: 'createdAt',
        defaultVisible: false,
        render: (date: unknown) => {
          if (!date || typeof date !== 'string') return '-';
          return new Date(date).toLocaleDateString('pt-BR');
        }
      },
      {
        key: 'actions',
        title: 'Ações',
        dataIndex: 'actions',
        defaultVisible: true,
        render: (_, record: Animal) => {
          return <ActionButton animalId={record.id} />;
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
    
    return selectedColumns.map(col => ({
      title: col.title,
      dataIndex: col.dataIndex,
      key: col.key,
      render: col.render
    }));
  }
}

export const useAnimalColumnBuilder = () => {
  const builder = new AnimalColumnBuilder();
  
  return {
    getAllColumns: () => builder.getAllColumns(),
    getDefaultColumns: () => builder.getDefaultColumns(),
    getDefaultColumnKeys: () => builder.getDefaultColumnKeys(),
    getColumnsByKeys: (keys: string[]) => builder.getColumnsByKeys(keys),
    getColumnOptions: () => builder.getColumnOptions(),
    buildTableColumns: (selectedKeys: string[]) => builder.buildTableColumns(selectedKeys)
  };
}; 