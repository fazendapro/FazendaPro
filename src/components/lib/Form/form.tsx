import React from 'react';
import { FormProvider, useFormContext, Controller, FieldValues } from 'react-hook-form';
import { Form as AntForm, Checkbox, Row, Col, Typography, Input } from 'antd';
import { SubmitHandler, UseFormReturn } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { FieldType } from '../../../types/field-types';

const { Link } = Typography;

interface FormProps<T extends FieldValues> {
  onSubmit: SubmitHandler<T>;
  fields: FieldType[];
  children?: React.ReactNode;
  methods: UseFormReturn<T>;
}

const FieldRenderer: React.FC<{ fields: FieldType[] }> = ({ fields }) => {
  const { control, formState: { errors } } = useFormContext();
  const { t } = useTranslation();

  return (
    <>
      {fields.map((field, index) => (
        <Col span={field.colSpan || 24} key={`field-${index}`}>
          {field.type === 'checkbox' ? (
            <div style={{ marginBottom: 16 }}>
              <Controller
                name={field.name}
                control={control}
                render={({ field: { onChange, value } }) => (
                  <Checkbox checked={value} onChange={e => onChange(e.target.checked)}>
                    {t(field.label)}
                  </Checkbox>
                )}
              />
            </div>
          ) : field.type === 'link' ? (
            <Link href="#" style={{ float: 'right' }}>
              {t(field.label)}
            </Link>
          ) : (
            <AntForm.Item
              label={t(field.label)}
              validateStatus={errors[field.name] ? 'error' : ''}
              help={errors[field.name]?.message as string}
            >
              <Controller
                name={field.name}
                control={control}
                render={({ field: { onChange, value } }) => (
                  field.type === 'password' ? (
                    <Input.Password
                      value={value}
                      onChange={onChange}
                      placeholder={field.placeholder ? t(field.placeholder) : undefined}
                    />
                  ) : (
                    <Input
                      value={value}
                      onChange={onChange}
                      placeholder={field.placeholder ? t(field.placeholder) : undefined}
                    />
                  )
                )}
              />
            </AntForm.Item>
          )}
        </Col>
      ))}
    </>
  );
};

export const Form = <T extends FieldValues>({ onSubmit, fields, children, methods }: FormProps<T>) => {
  const { handleSubmit } = methods;

  return (
    <FormProvider {...methods}>
      <AntForm onFinish={handleSubmit(onSubmit)} layout="vertical">
        <Row gutter={16}>
          <FieldRenderer fields={fields} />
        </Row>
        {children}
      </AntForm>
    </FormProvider>
  );
};