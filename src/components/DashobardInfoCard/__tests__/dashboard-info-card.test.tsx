import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { DashboardInfoCard } from '../dashboard-info-card'
import { UserOutlined } from '@ant-design/icons'

describe('DashboardInfoCard', () => {
  it('deve renderizar título e valor corretamente', () => {
    render(
      <DashboardInfoCard
        title="Total de Animais"
        value="150"
        icon={<UserOutlined />}
        isLast={false}
      />
    )

    expect(screen.getByText('Total de Animais')).toBeInTheDocument()
    expect(screen.getByText('150')).toBeInTheDocument()
  })

  it('deve renderizar ícone corretamente', () => {
    render(
      <DashboardInfoCard
        title="Total de Animais"
        value="150"
        icon={<UserOutlined data-testid="icon" />}
        isLast={false}
      />
    )

    expect(screen.getByTestId('icon')).toBeInTheDocument()
  })

  it('deve não ter borda direita quando isLast é true', () => {
    const { container } = render(
      <DashboardInfoCard
        title="Total de Animais"
        value="150"
        icon={<UserOutlined />}
        isLast={true}
      />
    )

    const card = container.querySelector('.ant-card')
    expect(card).toHaveStyle({ borderRight: 'none' })
  })

  it('deve ter borda direita quando isLast é false', () => {
    const { container } = render(
      <DashboardInfoCard
        title="Total de Animais"
        value="150"
        icon={<UserOutlined />}
        isLast={false}
      />
    )

    const card = container.querySelector('.ant-card')
    expect(card).toHaveStyle({ borderRight: '1px solid rgba(0,0,0,0.1)' })
  })

  it('deve ter estilo de texto centralizado', () => {
    const { container } = render(
      <DashboardInfoCard
        title="Total de Animais"
        value="150"
        icon={<UserOutlined />}
        isLast={false}
      />
    )

    const card = container.querySelector('.ant-card')
    expect(card).toHaveStyle({ textAlign: 'center' })
  })
})

