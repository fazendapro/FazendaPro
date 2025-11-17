import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { useForm } from 'react-hook-form'
import { Form } from '../form'
import { FieldType } from '../../../../types/field-types'

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key
  })
}))

vi.mock('antd', async (importOriginal) => {
  const antd = await importOriginal<typeof import('antd')>()
  return {
    ...antd,
    Form: {
      ...antd.Form,
      Item: ({ children, label, validateStatus, help }: { children?: React.ReactNode; label?: string; validateStatus?: string; help?: string }) => (
        <div data-testid="form-item" data-status={validateStatus} data-help={help}>
          {label && <label data-testid="form-label">{label}</label>}
          {children}
          {help && <div data-testid="form-help">{help}</div>}
        </div>
      )
    },
    Input: {
      ...antd.Input,
      Password: ({ value, onChange, placeholder, ...props }: { value?: string; onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void; placeholder?: string; [key: string]: unknown }) => (
        <input
          type="password"
          value={value}
          onChange={onChange}
          placeholder={placeholder}
          data-testid="password-input"
          {...props}
        />
      )
    },
    InputNumber: ({ value, onChange, placeholder, ...props }: { value?: number; onChange?: (value: number | null) => void; placeholder?: string; [key: string]: unknown }) => (
      <input
        type="number"
        value={value}
        onChange={(e) => onChange?.(e.target.value ? Number(e.target.value) : null)}
        placeholder={placeholder}
        data-testid="number-input"
        {...props}
      />
    ),
    Checkbox: ({ checked, onChange, children }: { checked?: boolean; onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void; children?: React.ReactNode }) => (
      <label>
        <input
          type="checkbox"
          checked={checked}
          onChange={onChange}
          data-testid="checkbox-input"
        />
        {children}
      </label>
    ),
    Row: ({ children, gutter }: { children?: React.ReactNode; gutter?: number | [number, number] }) => (
      <div data-testid="form-row" data-gutter={gutter}>
        {children}
      </div>
    ),
    Col: ({ children, span }: { children?: React.ReactNode; span?: number }) => (
      <div data-testid="form-col" data-span={span}>
        {children}
      </div>
    ),
    Typography: {
      Link: ({ children, href, style }: { children?: React.ReactNode; href?: string; style?: React.CSSProperties }) => (
        <a href={href} style={style} data-testid="form-link">
          {children}
        </a>
      )
    }
  }
})

const TestWrapper: React.FC<{
  fields: FieldType[]
  onSubmit: (data: Record<string, unknown>) => void
  children?: React.ReactNode
}> = ({ fields, onSubmit, children }) => {
  const methods = useForm({
    defaultValues: {
      name: '',
      email: '',
      password: '',
      age: 0,
      terms: false
    }
  })

  return (
    <Form
      fields={fields}
      onSubmit={onSubmit}
      methods={methods}
    >
      {children}
    </Form>
  )
}

describe('Form', () => {
  const mockOnSubmit = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
  })

  const basicFields: FieldType[] = [
    {
      name: 'name',
      label: 'Nome',
      type: 'text',
      placeholder: 'Digite seu nome'
    },
    {
      name: 'email',
      label: 'Email',
      type: 'text',
      placeholder: 'Digite seu email'
    }
  ]

  it('deve renderizar o formulário com campos básicos', () => {
    render(<TestWrapper fields={basicFields} onSubmit={mockOnSubmit} />)

    expect(screen.getByTestId('form-row')).toBeInTheDocument()
    expect(screen.getAllByTestId('form-col')).toHaveLength(2)
    expect(screen.getAllByTestId('form-item')).toHaveLength(2)
  })

  it('deve renderizar campos de texto corretamente', () => {
    render(<TestWrapper fields={basicFields} onSubmit={mockOnSubmit} />)

    const nameInput = screen.getByPlaceholderText('Digite seu nome')
    const emailInput = screen.getByPlaceholderText('Digite seu email')
    expect(nameInput).toBeInTheDocument()
    expect(emailInput).toBeInTheDocument()
    expect(nameInput).toHaveAttribute('placeholder', 'Digite seu nome')
    expect(emailInput).toHaveAttribute('placeholder', 'Digite seu email')
  })

  it('deve renderizar campo de senha', () => {
    const passwordFields: FieldType[] = [
      {
        name: 'password',
        label: 'Senha',
        type: 'password',
        placeholder: 'Digite sua senha'
      }
    ]

    render(<TestWrapper fields={passwordFields} onSubmit={mockOnSubmit} />)

    const passwordInput = screen.getByTestId('password-input')
    expect(passwordInput).toBeInTheDocument()
    expect(passwordInput).toHaveAttribute('type', 'password')
    expect(passwordInput).toHaveAttribute('placeholder', 'Digite sua senha')
  })

  it('deve renderizar campo numérico', () => {
    const numberFields: FieldType[] = [
      {
        name: 'age',
        label: 'Idade',
        type: 'number',
        placeholder: 'Digite sua idade'
      }
    ]

    render(<TestWrapper fields={numberFields} onSubmit={mockOnSubmit} />)

    const numberInput = screen.getByTestId('number-input')
    expect(numberInput).toBeInTheDocument()
    expect(numberInput).toHaveAttribute('type', 'number')
    expect(numberInput).toHaveAttribute('placeholder', 'Digite sua idade')
  })

  it('deve renderizar checkbox', () => {
    const checkboxFields: FieldType[] = [
      {
        name: 'terms',
        label: 'Aceito os termos',
        type: 'checkbox'
      }
    ]

    render(<TestWrapper fields={checkboxFields} onSubmit={mockOnSubmit} />)

    const checkbox = screen.getByTestId('checkbox-input')
    expect(checkbox).toBeInTheDocument()
    expect(checkbox).toHaveAttribute('type', 'checkbox')
  })

  it('deve renderizar link', () => {
    const linkFields: FieldType[] = [
      {
        name: 'forgot',
        label: 'Esqueci minha senha',
        type: 'link'
      }
    ]

    render(<TestWrapper fields={linkFields} onSubmit={mockOnSubmit} />)

    const link = screen.getByTestId('form-link')
    expect(link).toBeInTheDocument()
    expect(link).toHaveTextContent('Esqueci minha senha')
  })

  it('deve renderizar children quando fornecidos', () => {
    render(
      <TestWrapper fields={basicFields} onSubmit={mockOnSubmit}>
        <button type="submit" data-testid="submit-button">
          Enviar
        </button>
      </TestWrapper>
    )

    expect(screen.getByTestId('submit-button')).toBeInTheDocument()
  })

  it('deve usar colSpan personalizado', () => {
    const fieldsWithColSpan: FieldType[] = [
      {
        name: 'name',
        label: 'Nome',
        type: 'text',
        colSpan: 12
      }
    ]

    render(<TestWrapper fields={fieldsWithColSpan} onSubmit={mockOnSubmit} />)

    const cols = screen.getAllByTestId('form-col')
    const colWithSpan12 = cols.find(col => col.getAttribute('data-span') === '12')
    expect(colWithSpan12).toBeInTheDocument()
    expect(colWithSpan12).toHaveAttribute('data-span', '12')
  })

  it('deve usar colSpan padrão quando não especificado', () => {
    const { container } = render(<TestWrapper fields={basicFields} onSubmit={mockOnSubmit} />)

    const form = container.querySelector('form')
    const cols = form ? Array.from(form.querySelectorAll('[data-testid="form-col"]')) : []
    expect(cols.length).toBeGreaterThanOrEqual(2)
    cols.forEach(col => {
      expect(col).toHaveAttribute('data-span', '24')
    })
  })

  it('deve renderizar labels corretamente', () => {
    render(<TestWrapper fields={basicFields} onSubmit={mockOnSubmit} />)

    expect(screen.getByText('Nome')).toBeInTheDocument()
    expect(screen.getByText('Email')).toBeInTheDocument()
  })

  it('deve renderizar placeholders corretamente', () => {
    render(<TestWrapper fields={basicFields} onSubmit={mockOnSubmit} />)

    const inputs = screen.getAllByRole('textbox')
    expect(inputs[0]).toHaveAttribute('placeholder', 'Digite seu nome')
    expect(inputs[1]).toHaveAttribute('placeholder', 'Digite seu email')
  })

  it('deve permitir interação com campos', async () => {
    render(<TestWrapper fields={basicFields} onSubmit={mockOnSubmit} />)

    const nameInput = screen.getAllByRole('textbox')[0]
    
    fireEvent.change(nameInput, { target: { value: 'João Silva' } })
    
    await waitFor(() => {
      expect(nameInput).toHaveValue('João Silva')
    })
  })

  it('deve permitir interação com checkbox', async () => {
    const checkboxFields: FieldType[] = [
      {
        name: 'terms',
        label: 'Aceito os termos',
        type: 'checkbox'
      }
    ]

    render(<TestWrapper fields={checkboxFields} onSubmit={mockOnSubmit} />)

    const checkbox = screen.getByTestId('checkbox-input')
    
    fireEvent.click(checkbox)
    
    await waitFor(() => {
      expect(checkbox).toBeChecked()
    })
  })

  it('deve renderizar com campos vazios sem erro', () => {
    const { container } = render(<TestWrapper fields={[]} onSubmit={mockOnSubmit} />)

    const formRow = container.querySelector('[data-testid="form-row"]')
    expect(formRow).toBeInTheDocument()
    const formCols = container.querySelectorAll('[data-testid="form-col"]')
    expect(formCols.length).toBe(0)
  })

  it('deve renderizar múltiplos tipos de campo', () => {
    const mixedFields: FieldType[] = [
      {
        name: 'name',
        label: 'Nome',
        type: 'text'
      },
      {
        name: 'password',
        label: 'Senha',
        type: 'password'
      },
      {
        name: 'age',
        label: 'Idade',
        type: 'number'
      },
      {
        name: 'terms',
        label: 'Aceito os termos',
        type: 'checkbox'
      },
      {
        name: 'forgot',
        label: 'Esqueci minha senha',
        type: 'link'
      }
    ]

    render(<TestWrapper fields={mixedFields} onSubmit={mockOnSubmit} />)

    expect(screen.getByRole('textbox')).toBeInTheDocument()
    expect(screen.getByTestId('password-input')).toBeInTheDocument()
    expect(screen.getByTestId('number-input')).toBeInTheDocument()
    expect(screen.getByTestId('checkbox-input')).toBeInTheDocument()
    expect(screen.getByTestId('form-link')).toBeInTheDocument()
  })
})