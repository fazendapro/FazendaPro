import React, { forwardRef, ForwardRefRenderFunction, useEffect, useState } from 'react';
import { Input as AntInput, Form, Typography, Space, Tooltip } from 'antd';
import type { InputRef } from 'antd/es/input';
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
  const watchedValue = watch ? (watch(name as string) as unknown) : null;
  const value = typeof watchedValue === 'string' ? watchedValue : null;

  const typedBold = typeof bold === 'boolean' ? bold : undefined;
  const typedMtLabel = typeof mtLabel === 'number' ? mtLabel : 0;
  const typedMbField = typeof mbField === 'number' ? mbField : 0;
  const typedMax = typeof max === 'string' || typeof max === 'number' ? max : undefined;
  const typedIconLeft = (iconLeft && iconLeft !== false) ? (iconLeft as React.ReactNode) : undefined;
  const typedLabel = typeof label === 'string' ? label : (label as string | undefined);
  const typedTooltip = typeof tooltip === 'string' ? tooltip : (tooltip as string | undefined);
  const typedPlaceholder = typeof placeholder === 'string' ? placeholder : (placeholder as string | undefined);
  const typedName = typeof name === 'string' ? name : String(name);
  const typedAutoComplete = typeof autoComplete === 'string' ? autoComplete : 'on';

  useEffect(() => {
    if (value && typeof value === 'string') {
      setSize(value.length);
    } else {
      setSize(0);
    }
  }, [value]);

  return (
    <Form.Item
      label={
        typedLabel ? (
          <Space>
            <Text strong={typedBold} style={{ fontSize: 14, marginTop: typedMtLabel }}>
              {typedLabel}
            </Text>
            {isRequired ? (
              <Text style={{ color: 'red', fontSize: 14 }}>*</Text>
            ) : null}
            {typedTooltip && (
              <QuestionMarkToolTip
                tipText={typedTooltip}
                iconSize={16}
                style={{ marginLeft: 4 }}
              />
            )}
          </Space>
        ) : undefined
      }
      validateStatus={error || isInvalid ? 'error' : ''}
      help={error && typeof error === 'object' && 'message' in error && error.message ? String(error.message) : undefined}
      style={typedMbField > 0 ? { marginBottom: `${typedMbField}px` } : undefined}
    >
      <div style={{ position: 'relative' }}>
        <AntInput
          prefix={typedIconLeft}
          name={typedName}
          id={typedName}
          placeholder={typedPlaceholder || undefined}
          maxLength={typedMax ? Number(typedMax) : undefined}
          ref={ref as React.Ref<InputRef>}
          autoComplete={typedAutoComplete}
          {...(rest as Record<string, unknown>)}
        />
        {value && typedMax && Number(typedMax) > 0 && (
          <div style={{ position: 'absolute', right: 8, bottom: -20, fontSize: 12, color: '#8c8c8c' }}>
            {size} / {String(typedMax)}
          </div>
        )}
      </div>
    </Form.Item>
  );
};

export const Input = forwardRef(InputBase);