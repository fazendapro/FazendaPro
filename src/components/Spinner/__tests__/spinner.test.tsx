import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Spinner } from '../spinner'

vi.mock('antd', async (importOriginal) => {
  const antd = await importOriginal<typeof import('antd')>()
  return {
    ...antd,
    Spin: ({ size, ...props }: any) => (
      <div data-testid="ant-spin" data-size={size} {...props}>
        <div data-testid="spinner-content">Loading...</div>
      </div>
    )
  }
})

describe('Spinner', () => {
  it('deve renderizar o spinner corretamente', () => {
    render(<Spinner />)

    const spinner = screen.getByTestId('ant-spin')
    expect(spinner).toBeInTheDocument()
    expect(spinner).toHaveAttribute('data-size', 'large')
  })

  it('deve renderizar com tamanho large', () => {
    render(<Spinner />)

    const spinner = screen.getByTestId('ant-spin')
    expect(spinner).toHaveAttribute('data-size', 'large')
  })

  it('deve renderizar com conteúdo de loading', () => {
    render(<Spinner />)

    expect(screen.getByTestId('spinner-content')).toBeInTheDocument()
    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })

  it('deve ter estilo de container centralizado', () => {
    render(<Spinner />)

    const spinner = screen.getByTestId('ant-spin')
    const container = spinner.parentElement

    expect(container).toHaveStyle({
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100vh'
    })
  })

  it('deve renderizar sem erros', () => {
    expect(() => render(<Spinner />)).not.toThrow()
  })

  it('deve ser acessível', () => {
    render(<Spinner />)

    const spinner = screen.getByTestId('ant-spin')
    expect(spinner).toBeInTheDocument()
  })

  it('deve ocupar toda a altura da viewport', () => {
    render(<Spinner />)

    const container = screen.getByTestId('ant-spin').parentElement
    expect(container).toHaveStyle('height: 100vh')
  })

  it('deve centralizar horizontalmente', () => {
    render(<Spinner />)

    const container = screen.getByTestId('ant-spin').parentElement
    expect(container).toHaveStyle('justify-content: center')
  })

  it('deve centralizar verticalmente', () => {
    render(<Spinner />)

    const container = screen.getByTestId('ant-spin').parentElement
    expect(container).toHaveStyle('align-items: center')
  })

  it('deve usar flexbox para layout', () => {
    render(<Spinner />)

    const container = screen.getByTestId('ant-spin').parentElement
    expect(container).toHaveStyle('display: flex')
  })
})
