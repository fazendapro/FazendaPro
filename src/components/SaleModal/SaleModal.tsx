import React, { useEffect } from 'react';
import { Modal, Form, Input, InputNumber, DatePicker, Button, message, Select } from 'antd';
import { useTranslation } from 'react-i18next';
import { useSaleForm } from '../../hooks/useSale';
import { CreateSaleRequest, UpdateSaleRequest, Sale } from '../../types/sale';
import { useAnimals } from '../../pages/contents/AnimalTable/hooks/useAnimals';
import { useFarm } from '../../hooks/useFarm';
import dayjs from 'dayjs';

interface SaleModalProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
  sale?: Sale;
  mode: 'create' | 'edit';
}

export const SaleModal: React.FC<SaleModalProps> = ({
  visible,
  onCancel,
  onSuccess,
  sale,
  mode
}) => {
  const { t } = useTranslation();
  const { createSale, updateSale, loading } = useSaleForm();
  const { farm } = useFarm();
  const { animals, loading: animalsLoading } = useAnimals(farm?.id);
  const [form] = Form.useForm();

  useEffect(() => {
    if (visible) {
      if (mode === 'edit' && sale) {
        form.setFieldsValue({
          animal_id: sale.animal_id,
          buyer_name: sale.buyer_name,
          price: sale.price,
          sale_date: dayjs(sale.sale_date),
          notes: sale.notes,
        });
      } else {
        form.resetFields();
      }
    }
  }, [visible, mode, sale, form]);

  const handleSubmit = async (values: {
    animal_id: number;
    buyer_name: string;
    price: number;
    sale_date: dayjs.Dayjs;
    notes?: string;
  }) => {
    try {
      if (mode === 'create') {
        const saleData: CreateSaleRequest = {
          animal_id: values.animal_id,
          buyer_name: values.buyer_name,
          price: values.price,
          sale_date: values.sale_date.format('YYYY-MM-DD'),
          notes: values.notes,
        };
        await createSale(saleData);
      } else {
        const saleData: UpdateSaleRequest = {
          buyer_name: values.buyer_name,
          price: values.price,
          sale_date: values.sale_date.format('YYYY-MM-DD'),
          notes: values.notes,
        };
        await updateSale(sale!.id, saleData);
      }
      
      message.success(
        mode === 'create' 
          ? t('saleModal.success.created') 
          : t('saleModal.success.updated')
      );
      form.resetFields();
      onSuccess();
    } catch {
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title={
        mode === 'create' 
          ? t('saleModal.title.create') 
          : t('saleModal.title.edit')
      }
      open={visible}
      onCancel={handleCancel}
      footer={null}
      width={600}
      destroyOnHidden
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        initialValues={{
          sale_date: dayjs(),
        }}
      >
        <Form.Item
          name="animal_id"
          label={t('saleModal.fields.animal')}
          rules={[
            { required: true, message: t('saleModal.validation.animalRequired') }
          ]}
        >
          <Select
            placeholder={t('saleModal.placeholders.animal')}
            loading={animalsLoading}
            disabled={mode === 'edit'}
            showSearch
            optionFilterProp="children"
            filterOption={(input, option) =>
              (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
            }
          >
            {animals
              .filter(animal => animal.status === 0)
              .map(animal => (
                <Select.Option key={animal.id} value={animal.id}>
                  {animal.animal_name} - Brinco: {animal.ear_tag_number_local}
                </Select.Option>
              ))
            }
          </Select>
        </Form.Item>

        <Form.Item
          name="buyer_name"
          label={t('saleModal.fields.buyerName')}
          rules={[
            { required: true, message: t('saleModal.validation.buyerNameRequired') }
          ]}
        >
          <Input placeholder={t('saleModal.placeholders.buyerName')} />
        </Form.Item>

        <Form.Item
          name="price"
          label={t('saleModal.fields.price')}
          rules={[
            { required: true, message: t('saleModal.validation.priceRequired') },
            { type: 'number', min: 0.01, message: t('saleModal.validation.priceMin') }
          ]}
        >
          <InputNumber
            placeholder={t('saleModal.placeholders.price')}
            style={{ width: '100%' }}
            formatter={(value) => `R$ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
            parser={(value) => value!.replace(/R\$\s?|(,*)/g, '')}
            precision={2}
          />
        </Form.Item>

        <Form.Item
          name="sale_date"
          label={t('saleModal.fields.saleDate')}
          rules={[
            { required: true, message: t('saleModal.validation.saleDateRequired') }
          ]}
        >
          <DatePicker
            style={{ width: '100%' }}
            format="DD/MM/YYYY"
            placeholder={t('saleModal.placeholders.saleDate')}
          />
        </Form.Item>

        <Form.Item
          name="notes"
          label={t('saleModal.fields.notes')}
        >
          <Input.TextArea
            placeholder={t('saleModal.placeholders.notes')}
            rows={3}
          />
        </Form.Item>

        <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
          <Button onClick={handleCancel} style={{ marginRight: 8 }}>
            {t('common.cancel')}
          </Button>
          <Button type="primary" htmlType="submit" loading={loading}>
            {mode === 'create' ? t('common.create') : t('common.save')}
          </Button>
        </Form.Item>
      </Form>
    </Modal>
  );
};
