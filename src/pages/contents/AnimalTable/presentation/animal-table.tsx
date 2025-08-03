import React from 'react';
import { Table, Tag } from 'antd';

const columns = [
  { title: 'Nomes', dataIndex: 'name', key: 'name' },
  { title: 'Preço de Venda', dataIndex: 'price', key: 'price' },
  { title: 'Peso Atual', dataIndex: 'currentWeight', key: 'currentWeight' },
  { title: 'Peso Ideal', dataIndex: 'idealWeight', key: 'idealWeight' },
  { title: 'Última Atualização', dataIndex: 'lastUpdate', key: 'lastUpdate' },
  {
    title: 'Produção de Leite',
    dataIndex: 'milkProduction',
    key: 'milkProduction',
    render: (text: string) => {
      let color = 'green';
      if (text === 'Baixa') color = 'red';
      else if (text === 'Média') color = 'orange';
      return <Tag color={color}>{text}</Tag>;
    },
  },
];

const data = [
  { key: '1', name: 'Maggi', price: 'R$430', currentWeight: '43 Kg', idealWeight: '12 Kg', lastUpdate: '11/12/22', milkProduction: 'Alta' },
  { key: '2', name: 'Bru', price: 'R$257', currentWeight: '22 Kg', idealWeight: '12 Kg', lastUpdate: '21/12/22', milkProduction: 'Baixa' },
  { key: '3', name: 'Red Bull', price: 'R$405', currentWeight: '36 Kg', idealWeight: '9 Kg', lastUpdate: '5/12/22', milkProduction: 'Alta' },
  { key: '4', name: 'Bourn Vita', price: 'R$502', currentWeight: '14 Kg', idealWeight: '6 Kg', lastUpdate: '8/12/22', milkProduction: 'Baixa' },
  { key: '5', name: 'Horlicks', price: 'R$530', currentWeight: '5 Kg', idealWeight: '5 Kg', lastUpdate: '9/12/23', milkProduction: 'Alta' },
  { key: '6', name: 'Harpic', price: 'R$605', currentWeight: '10 Kg', idealWeight: '5 Kg', lastUpdate: '9/12/23', milkProduction: 'Alta' },
  { key: '7', name: 'Ariel', price: 'R$408', currentWeight: '23 Kg', idealWeight: '7 Kg', lastUpdate: '15/12/23', milkProduction: 'Baixa' },
  { key: '8', name: 'Scotch Brite', price: 'R$359', currentWeight: '43 Kg', idealWeight: '8 Kg', lastUpdate: '6/12/23', milkProduction: 'Alta' },
  { key: '9', name: 'Coca Cola', price: 'R$205', currentWeight: '41 Kg', idealWeight: '10 Kg', lastUpdate: '11/12/22', milkProduction: 'Média' },
];

const AnimalTable: React.FC = () => {
    return (
      <Table columns={columns} dataSource={data} pagination={{ showSizeChanger: true }} />
  );
};

export { AnimalTable };