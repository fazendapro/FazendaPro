import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { useForm, FormProvider } from 'react-hook-form'
import { Input } from '../input'

vi.mock('antd', async (importOriginal) => {
  const antd = await importOriginal<typeof import('antd')>()
  return {
    ...antd,
    Form: {
      Item: ({ children, label, validateStatus, help, style }: any) => (
        <div 
          data-testid="form-item" 
          data-status={validateStatus} 
          data-help={help}
          style={style}
        >
          {label && <div data-testid="form-label">{label}</div>}
          {children}
          {help && <div data-testid="form-help">{help}</div>}
        </div>
      )
    },
    Input: ({ prefix, name, id, placeholder, maxLength, autoComplete, ...props }: any) => (
      <input
        prefix={prefix}
        name={name}
        id={id}
        placeholder={placeholder}
        maxLength={maxLength}
        autoComplete={autoComplete}
        data-testid="ant-input"
        {...props}
      />
    ),
    Typography: {
      Text: ({ children, strong, style }: any) => (
        <span data-testid="text" data-strong={strong} style={style}>
          {children}
        </span>
      )
    },
    Space: ({ children }: any) => (
      <div data-testid="space">{children}</div>
    ),
    Tooltip: ({ children, title }: any) => (
      <div data-testid="tooltip" title={title}>
        {children}
      </div>
    )
  }
})

vi.mock('@ant-design/icons', () => ({
  QuestionCircleOutlined: ({ style }: any) => (
    <span data-testid="question-icon" style={style}>?</span>
  )
}))

const TestWrapper: React.FC<{
  children: React.ReactNode
  defaultValues?: any
}> = ({ children, defaultValues = {} }) => {
  const methods = useForm({
    defaultValues: {
      testField: '',
      ...defaultValues
    },
    mode: 'onChange'
  })

  return (
    <FormProvider {...methods}>
      {children}
    </FormProvider>
  )
}

describe('Input', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('deve renderizar input básico', () => {
    render(
      <TestWrapper>
        <Input name="testField" />
      </TestWrapper>
    )

    const input = screen.getByTestId('ant-input')
    expect(input).toBeInTheDocument()
    expect(input).toHaveAttribute('name', 'testField')
    expect(input).toHaveAttribute('id', 'testField')
  })

  it('deve renderizar com label', () => {
    render(
      <TestWrapper>
        <Input name="testField" label="Nome" />
      </TestWrapper>
    )

    expect(screen.getByTestId('form-label')).toBeInTheDocument()
    expect(screen.getByText('Nome')).toBeInTheDocument()
  })

  it('deve renderizar com placeholder', () => {
    render(
      <TestWrapper>
        <Input name="testField" placeholder="Digite seu nome" />
      </TestWrapper>
    )

    const input = screen.getByTestId('ant-input')
    expect(input).toHaveAttribute('placeholder', 'Digite seu nome')
  })

  it('deve renderizar como obrigatório', () => {
    render(
      <TestWrapper>
        <Input name="testField" label="Nome" isRequired />
      </TestWrapper>
    )

    expect(screen.getByText('*')).toBeInTheDocument()
  })

  it('deve renderizar com tooltip', () => {
    render(
      <TestWrapper>
        <Input name="testField" label="Nome" tooltip="Informação adicional" />
      </TestWrapper>
    )

    expect(screen.getByTestId('tooltip')).toBeInTheDocument()
    expect(screen.getByTestId('question-icon')).toBeInTheDocument()
  })

  it('deve renderizar com ícone à esquerda', () => {
    const icon = <span data-testid="custom-icon">@</span>
    
    render(
      <TestWrapper>
        <Input name="testField" iconLeft={icon} />
      </TestWrapper>
    )

    const input = screen.getByTestId('ant-input')
    expect(input).toHaveAttribute('prefix')
  })

  it('deve renderizar com erro', () => {
    const error = { message: 'Campo obrigatório' }
    
    render(
      <TestWrapper>
        <Input name="testField" error={error} />
      </TestWrapper>
    )

    const formItem = screen.getByTestId('form-item')
    expect(formItem).toHaveAttribute('data-status', 'error')
    expect(screen.getByTestId('form-help')).toHaveTextContent('Campo obrigatório')
  })

  it('deve renderizar como inválido', () => {
    render(
      <TestWrapper>
        <Input name="testField" isInvalid />
      </TestWrapper>
    )

    const formItem = screen.getByTestId('form-item')
    expect(formItem).toHaveAttribute('data-status', 'error')
  })

  it('deve renderizar com maxLength', () => {
    render(
      <TestWrapper>
        <Input name="testField" max="50" />
      </TestWrapper>
    )

    const input = screen.getByTestId('ant-input')
    expect(input).toHaveAttribute('maxLength', '50')
  })

  it('deve renderizar contador de caracteres quando maxLength é definido', () => {
    render(
      <TestWrapper defaultValues={{ testField: 'test' }}>
        <Input name="testField" max="50" />
      </TestWrapper>
    )

    const counterElement = screen.getByText(/4.*50|50.*4/)
    expect(counterElement).toBeInTheDocument()
  })

  it('deve renderizar contador com valor inicial correto', () => {
    render(
      <TestWrapper defaultValues={{ testField: 'test' }}>
        <Input name="testField" max="50" />
      </TestWrapper>
    )

    const counterElement = screen.getByText(/4.*50|50.*4/)
    expect(counterElement).toBeInTheDocument()
  })

  it('deve renderizar contador quando há valor e maxLength', () => {
    render(
      <TestWrapper defaultValues={{ testField: 'abc' }}>
        <Input name="testField" max="100" />
      </TestWrapper>
    )

    const counterElement = screen.getByText(/3.*100|100.*3/)
    expect(counterElement).toBeInTheDocument()
  })

  it('deve renderizar com autoComplete off', () => {
    render(
      <TestWrapper>
        <Input name="testField" autoComplete="off" />
      </TestWrapper>
    )

    const input = screen.getByTestId('ant-input')
    expect(input).toHaveAttribute('autoComplete', 'off')
  })

  it('deve renderizar com autoComplete on (padrão)', () => {
    render(
      <TestWrapper>
        <Input name="testField" />
      </TestWrapper>
    )

    const input = screen.getByTestId('ant-input')
    expect(input).toHaveAttribute('autoComplete', 'on')
  })

  it('deve renderizar com estilo personalizado', () => {
    render(
      <TestWrapper>
        <Input name="testField" label="Nome" bold />
      </TestWrapper>
    )

    const text = screen.getByTestId('text')
    expect(text).toHaveAttribute('data-strong', 'true')
  })

  it('deve renderizar com margem personalizada', () => {
    render(
      <TestWrapper>
        <Input name="testField" mtLabel={10} mbField={20} />
      </TestWrapper>
    )

    const formItem = screen.getByTestId('form-item')
    expect(formItem).toHaveStyle('margin-bottom: 20px')
  })

  it('deve permitir interação com o input', async () => {
    render(
      <TestWrapper>
        <Input name="testField" />
      </TestWrapper>
    )

    const input = screen.getByTestId('ant-input')
    
    fireEvent.change(input, { target: { value: 'João Silva' } })
    
    await waitFor(() => {
      expect(input).toHaveValue('João Silva')
    })
  })

  it('deve renderizar sem label quando não fornecido', () => {
    render(
      <TestWrapper>
        <Input name="testField" />
      </TestWrapper>
    )

    expect(screen.queryByTestId('form-label')).not.toBeInTheDocument()
  })

  it('deve renderizar sem tooltip quando não fornecido', () => {
    render(
      <TestWrapper>
        <Input name="testField" label="Nome" />
      </TestWrapper>
    )

    expect(screen.queryByTestId('tooltip')).not.toBeInTheDocument()
  })

  it('deve renderizar sem ícone quando iconLeft é false', () => {
    render(
      <TestWrapper>
        <Input name="testField" iconLeft={false} />
      </TestWrapper>
    )

    const input = screen.getByTestId('ant-input')
    expect(input).not.toHaveAttribute('prefix')
  })

  it('deve renderizar com todas as props customizadas', () => {
    const icon = <span data-testid="custom-icon">@</span>
    const error = { message: 'Erro de validação' }
    
    render(
      <TestWrapper>
        <Input
          name="testField"
          label="Email"
          placeholder="Digite seu email"
          isRequired
          tooltip="Digite um email válido"
          iconLeft={icon}
          error={error}
          max="100"
          autoComplete="off"
          bold
          mtLabel={5}
          mbField={10}
        />
      </TestWrapper>
    )

    expect(screen.getByText('Email')).toBeInTheDocument()
    expect(screen.getByText('*')).toBeInTheDocument()
    expect(screen.getByTestId('tooltip')).toBeInTheDocument()
    expect(screen.getByTestId('form-help')).toHaveTextContent('Erro de validação')
    
    const input = screen.getByTestId('ant-input')
    expect(input).toHaveAttribute('prefix')
    expect(input).toHaveAttribute('placeholder', 'Digite seu email')
    expect(input).toHaveAttribute('maxLength', '100')
    expect(input).toHaveAttribute('autoComplete', 'off')
  })

  it('deve renderizar sem contador quando maxLength não é definido', () => {
    render(
      <TestWrapper defaultValues={{ testField: 'test' }}>
        <Input name="testField" />
      </TestWrapper>
    )

    expect(screen.queryByText(/\/ \d+/)).not.toBeInTheDocument()
  })

  it('deve renderizar sem contador quando maxLength é 0', () => {
    render(
      <TestWrapper defaultValues={{ testField: 'test' }}>
        <Input name="testField" max="0" />
      </TestWrapper>
    )

    expect(screen.queryByText(/\/ \d+/)).not.toBeInTheDocument()
  })
})
