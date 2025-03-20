import React from 'react';
import { FormProvider, useFormContext } from 'react-hook-form';
import { Form as AntForm, Checkbox, Row, Col, Typography, Input } from 'antd';
import { FieldValues, SubmitHandler, UseFormReturn } from 'react-hook-form';

const { Link } = Typography;

type FieldType = {
  name: string;
  label: string;
  type: 'text' | 'password' | 'checkbox' | 'link';
  placeholder?: string;
  colSpan?: number;
  isRequired?: boolean;
};

interface FormProps<T extends FieldValues> {
  onSubmit: SubmitHandler<T>;
  fields: FieldType[];
  children?: React.ReactNode;
  methods: UseFormReturn<T>;
}

const FieldRenderer: React.FC<{ fields: FieldType[] }> = ({ fields }) => {
  const { register, formState: { errors } } = useFormContext();

  return (
    <>
      {fields.map((field, index) => (
        <Col span={field.colSpan || 24} key={`field-${index}`}>
          {field.type === 'checkbox' ? (
            <div style={{ marginBottom: 16 }}>
              <Checkbox {...register(field.name)}>
                {field.label}
              </Checkbox>
            </div>
          ) : field.type === 'link' ? (
            <Link href="#" style={{ float: 'right' }}>
              {field.label}
            </Link>
          ) : (
            <AntForm.Item
              label={field.label}
              validateStatus={errors[field.name] ? 'error' : ''}
              help={errors[field.name]?.message as string}
            >
              {field.type === 'password' ? (
                <Input.Password
                  placeholder={field.placeholder}
                  {...register(field.name)}
                />
              ) : (
                <Input
                  placeholder={field.placeholder}
                  {...register(field.name)}
                />
              )}
            </AntForm.Item>
          )}
        </Col>
      ))}
    </>
  );
};

export const Form = <T extends FieldValues>({ onSubmit, fields, children, methods }: FormProps<T>) => {
  return (
    <FormProvider {...methods}>
      <AntForm onFinish={methods.handleSubmit(onSubmit)} layout="vertical">
        <Row gutter={16}>
          <FieldRenderer fields={fields} />
        </Row>
        {children}
      </AntForm>
    </FormProvider>
  );
};