import React, { forwardRef, ForwardRefRenderFunction, useEffect, useState } from 'react';
import { Input as AntInput, Form, Typography, Space, Tooltip, InputRef } from 'antd';
import { FieldError, FieldErrors, Merge, useFormContext } from 'react-hook-form';
import { QuestionCircleOutlined } from '@ant-design/icons';

const { Text } = Typography;

interface InputProps {
  name: string;
  label?: string;
  error?: FieldError | Merge<FieldError, FieldErrors<Date>> | null | undefined;
  isRequired?: boolean;
  isInvalid?: boolean;
  iconLeft?: React.ReactNode | boolean;
  tooltip?: string;
  bold?: boolean;
  placeholder?: string;
  mtLabel?: number;
  maxNumber?: string | number;
  mbField?: number;
  autoComplete?: 'off' | 'on';
  [key: string]: unknown;
}

const QuestionMarkToolTip: React.FC<{ tipText: string; iconSize?: number; style?: React.CSSProperties }> = ({
  tipText,
  iconSize = 14,
  style,
}) => (
  <Tooltip title={tipText}>
    <QuestionCircleOutlined style={{ fontSize: iconSize, marginLeft: 4, ...style }} />
  </Tooltip>
);

const InputBase: ForwardRefRenderFunction<HTMLInputElement, Omit<InputProps, "ref">> = (
  {
    name,
    isRequired = false,
    label,
    max,
    error = null,
    isInvalid = false,
    iconLeft = false,
    tooltip,
    bold,
    placeholder = '',
    mtLabel = 0,
    mbField = 0,
    autoComplete = 'on',
    ...rest
  },
  ref
) => {
  const [size, setSize] = useState(0);
  const formContext = useFormContext();
  const watch = formContext?.watch;
  const value = watch ? watch(name) : null;

  useEffect(() => {
    if (value) {
      setSize(value.length);
    } else {
      setSize(0);
    }
  }, [value]);

  return (
    <Form.Item
      label={
        label && (
          <Space>
            <Text strong={bold} style={{ fontSize: 14, marginTop: mtLabel }}>
              {label}
            </Text>
            {isRequired && (
              <Text style={{ color: 'red', fontSize: 14 }}>*</Text>
            )}
            {tooltip && (
              <QuestionMarkToolTip
                tipText={tooltip}
                iconSize={16}
                style={{ marginLeft: 4 }}
              />
            )}
          </Space>
        )
      }
      validateStatus={error || isInvalid ? 'error' : ''}
      help={error?.message as string}
      style={mbField > 0 ? { marginBottom: `${mbField}px` } : undefined}
    >
      <div style={{ position: 'relative' }}>
        <AntInput
          prefix={iconLeft || null}
          name={name}
          id={name}
          placeholder={placeholder}
          maxLength={max ? Number(max) : undefined}
          ref={ref as React.Ref<InputRef>}
          autoComplete={autoComplete}
          {...rest}
        />
        {value && max && Number(max) > 0 && (
          <div style={{ position: 'absolute', right: 8, bottom: -20, fontSize: 12, color: '#8c8c8c' }}>
            {size} / {max}
          </div>
        )}
      </div>
    </Form.Item>
  );
};

export const Input = forwardRef(InputBase);