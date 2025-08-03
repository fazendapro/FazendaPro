import React from 'react';
import { Button, Input, Select, Modal } from 'antd';
import { useModal } from '../../../../hooks';
const { Search } = Input;
const { Option } = Select;

const AnimalDashboard: React.FC = () => {
  const { isOpen, onOpen, onClose } = useModal();

  return (
    <div style={{ background: '#f5f5f5', borderRadius: '8px' }}>
      <div style={{ marginBottom: '20px', background: '#fff', padding: '10px', borderRadius: '4px' }}>
        <span style={{ color: '#1890ff' }}>Categorias: 14</span> | <span style={{ color: '#faad14' }}>Total de animais: 868</span> | <span style={{ color: '#000' }}>Receita: R$25000</span>
        <span style={{ marginLeft: '20px', color: '#722ed1' }}>Melhores Vendas: 5</span> | <span style={{ color: '#000' }}>Custo: R$2500</span>
        <span style={{ marginLeft: '20px', color: '#ff4d4f' }}>Em baixa produção: 12</span> | <span style={{ color: '#000' }}>Semindas: 2</span> | <span style={{ color: '#000' }}>Não semindas: 2</span>
      </div>
      <div style={{ marginBottom: '20px', background: '#fff', padding: '10px', borderRadius: '4px', display: 'flex', alignItems: 'center' }}>
        <Button type="primary" style={{ marginRight: '10px' }} onClick={onOpen}>Criar Vaca</Button>
        <Select defaultValue="Filtro" style={{ width: 120, marginRight: '10px' }}>
          <Option value="filtro">Filtro</Option>
        </Select>
        <Search placeholder="Importar csv" style={{ width: 200 }} />
      </div>
      <Modal
        title="Novo Gado"
        open={isOpen}
        onOk={onClose}
        onCancel={onClose}
        footer={[
          <Button key="back" onClick={onClose}>
            Abatear
          </Button>,
          <Button key="submit" type="primary" onClick={onClose}>
            Adicionar gado
          </Button>,
        ]}
      >
        <div style={{ padding: '20px' }}>
          <div style={{ marginBottom: '10px' }}>
            <label>Nome do Animal</label>
            <Input placeholder="Coloque o Nome do Animal" />
          </div>
          <div style={{ marginBottom: '10px' }}>
            <label>Número do Brinco</label>
            <Input placeholder="Coloque o número do brinco" />
          </div>
          <div style={{ marginBottom: '10px' }}>
            <label>Categoria</label>
            <Select placeholder="Selecione o tipo do animal">
              <Option value="vaca">Vaca</Option>
            </Select>
          </div>
          <div style={{ marginBottom: '10px' }}>
            <label>Preço da Compra</label>
            <Input placeholder="Coloque o preço da compra" />
          </div>
          <div style={{ marginBottom: '10px' }}>
            <label>Data de Nascimento</label>
            <Input placeholder="Coloque a data de nascimento" />
          </div>
          <div style={{ marginBottom: '10px' }}>
            <label>Peso</label>
            <Input placeholder="Coloque o peso de nascimento" />
          </div>
          <div style={{ marginBottom: '10px' }}>
            <label>Pais</label>
            <Input placeholder="Enter threshold value" />
          </div>
        </div>
      </Modal>
    </div>
  );
};

export { AnimalDashboard };