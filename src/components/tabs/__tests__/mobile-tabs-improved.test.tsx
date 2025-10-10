import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { MobileTabs } from '../mobile-tabs'
import { Tab } from '../types'

vi.mock('@ant-design/icons', () => ({
  UnorderedListOutlined: () => <div data-testid="list-icon">List</div>,
  ExperimentOutlined: () => <div data-testid="experiment-icon">Experiment</div>,
  HeartOutlined: () => <div data-testid="heart-icon">Heart</div>
}))

describe('MobileTabs - Improved Layout', () => {
  const mockTabs: Tab[] = [
    {
      title: 'Lista de Animais',
      name: 'animals-list',
      component: <div data-testid="animals-content">Animals Content</div>,
      isDisabled: false
    },
    {
      title: 'Produção de Leite',
      name: 'milk-production',
      component: <div data-testid="milk-content">Milk Content</div>,
      isDisabled: false
    },
    {
      title: 'Reprodução',
      name: 'reproduction',
      component: <div data-testid="reproduction-content">Reproduction Content</div>,
      isDisabled: false
    }
  ]

  it('deve renderizar ícones lado a lado em layout horizontal', () => {
    render(<MobileTabs tabs={mockTabs} />)
    
    expect(screen.getByTestId('list-icon')).toBeInTheDocument()
    expect(screen.getByTestId('experiment-icon')).toBeInTheDocument()
    expect(screen.getByTestId('heart-icon')).toBeInTheDocument()
  })

  it('deve mostrar tooltip com nome completo da aba ao passar o mouse', () => {
    render(<MobileTabs tabs={mockTabs} />)

    expect(screen.getByText('Lista')).toBeInTheDocument()
    expect(screen.getByText('Produção')).toBeInTheDocument()
    expect(screen.getByText('Reprodução')).toBeInTheDocument()
  })

  it('deve mostrar apenas a primeira palavra do título no botão', () => {
    render(<MobileTabs tabs={mockTabs} />)
    
    expect(screen.getByText('Lista')).toBeInTheDocument()
    expect(screen.getByText('Produção')).toBeInTheDocument()
    expect(screen.getByText('Reprodução')).toBeInTheDocument()
  })

  it('deve destacar a aba selecionada com estilo primary', () => {
    render(<MobileTabs tabs={mockTabs} defaultTabIndex={1} />)
    
    const buttons = screen.getAllByRole('button')
    expect(buttons[1]).toHaveClass('ant-btn-primary')
    expect(buttons[0]).toHaveClass('ant-btn-text')
    expect(buttons[2]).toHaveClass('ant-btn-text')
  })

  it('deve trocar de aba ao clicar em um ícone', () => {
    render(<MobileTabs tabs={mockTabs} />)
    
    expect(screen.getByTestId('animals-content')).toBeInTheDocument()
    
    const buttons = screen.getAllByRole('button')
    fireEvent.click(buttons[1])
    
    expect(screen.getByTestId('milk-content')).toBeInTheDocument()
    expect(screen.queryByTestId('animals-content')).not.toBeInTheDocument()
  })

  it('deve chamar onChange e onTabSelect ao trocar de aba', () => {
    const onChange = vi.fn()
    const onTabSelect = vi.fn()
    
    render(
      <MobileTabs 
        tabs={mockTabs} 
        onChange={onChange}
        onTabSelect={onTabSelect}
      />
    )
    
    const buttons = screen.getAllByRole('button')
    fireEvent.click(buttons[2])
    
    expect(onChange).toHaveBeenCalledWith(2)
    expect(onTabSelect).toHaveBeenCalledWith('reproduction')
  })

  it('deve aplicar estilos responsivos corretos', () => {
    render(<MobileTabs tabs={mockTabs} />)
    
    const buttons = screen.getAllByRole('button')
    expect(buttons).toHaveLength(3)
    
    buttons.forEach(button => {
      expect(button).toHaveStyle({
        minWidth: '60px',
        height: '50px',
        borderRadius: '8px'
      })
    })
  })

  it('deve desabilitar aba quando isDisabled for true', () => {
    const disabledTabs: Tab[] = [
      ...mockTabs.slice(0, 2),
      { ...mockTabs[2], isDisabled: true }
    ]
    
    render(<MobileTabs tabs={disabledTabs} />)
    
    const buttons = screen.getAllByRole('button')
    expect(buttons[2]).toBeDisabled()
  })

  it('deve renderizar ícone correto para cada tipo de aba', () => {
    render(<MobileTabs tabs={mockTabs} />)
    
    expect(screen.getByTestId('list-icon')).toBeInTheDocument()
    expect(screen.getByTestId('experiment-icon')).toBeInTheDocument()
    expect(screen.getByTestId('heart-icon')).toBeInTheDocument()
  })

  it('deve ter layout flexível com espaçamento adequado', () => {
    render(<MobileTabs tabs={mockTabs} />)
    
    const buttons = screen.getAllByRole('button')
    buttons.forEach(button => {
      expect(button).toHaveStyle({
        minWidth: '60px',
        height: '50px',
        borderRadius: '8px'
      })
    })
  })
})
