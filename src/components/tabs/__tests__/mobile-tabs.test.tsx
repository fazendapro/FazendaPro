import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '../../../test/test-utils'
import { MobileTabs } from '../mobile-tabs'
import { Tab } from '../types'

const mockTabs: Tab[] = [
  {
    title: 'Tab 1',
    name: 'tab1',
    component: <div data-testid="tab1-content">Conteúdo da Tab 1</div>,
  },
  {
    title: 'Tab 2',
    name: 'tab2',
    component: <div data-testid="tab2-content">Conteúdo da Tab 2</div>,
  },
  {
    title: 'Tab 3',
    name: 'tab3',
    component: <div data-testid="tab3-content">Conteúdo da Tab 3</div>,
    isDisabled: true,
  },
]

describe('MobileTabs', () => {
  it('deve renderizar as tabs como botões', () => {
    render(<MobileTabs tabs={mockTabs} />)
    
    expect(screen.getByRole('button', { name: 'Tab 1' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Tab 2' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Tab 3' })).toBeInTheDocument()
  })

  it('deve mostrar o conteúdo da primeira tab por padrão', () => {
    render(<MobileTabs tabs={mockTabs} />)
    
    expect(screen.getByTestId('tab1-content')).toBeInTheDocument()
    expect(screen.queryByTestId('tab2-content')).not.toBeInTheDocument()
    expect(screen.queryByTestId('tab3-content')).not.toBeInTheDocument()
  })

  it('deve permitir mudança de tab', () => {
    render(<MobileTabs tabs={mockTabs} />)
    
    fireEvent.click(screen.getByRole('button', { name: 'Tab 2' }))
    
    expect(screen.queryByTestId('tab1-content')).not.toBeInTheDocument()
    expect(screen.getByTestId('tab2-content')).toBeInTheDocument()
    expect(screen.queryByTestId('tab3-content')).not.toBeInTheDocument()
  })

  it('deve chamar onChange quando uma tab é selecionada', () => {
    const onChange = vi.fn()
    render(<MobileTabs tabs={mockTabs} onChange={onChange} />)
    
    fireEvent.click(screen.getByRole('button', { name: 'Tab 2' }))
    
    expect(onChange).toHaveBeenCalledWith(1)
  })

  it('deve chamar onTabSelect quando uma tab é selecionada', () => {
    const onTabSelect = vi.fn()
    render(<MobileTabs tabs={mockTabs} onTabSelect={onTabSelect} />)
    
    fireEvent.click(screen.getByRole('button', { name: 'Tab 2' }))
    
    expect(onTabSelect).toHaveBeenCalledWith('tab2')
  })

  it('deve desabilitar tabs marcadas como disabled', () => {
    render(<MobileTabs tabs={mockTabs} />)
    
    const disabledButton = screen.getByRole('button', { name: 'Tab 3' })
    expect(disabledButton).toBeDisabled()
  })

  it('deve usar defaultTabIndex para definir tab inicial', () => {
    render(<MobileTabs tabs={mockTabs} defaultTabIndex={1} />)
    
    expect(screen.getByTestId('tab2-content')).toBeInTheDocument()
    expect(screen.queryByTestId('tab1-content')).not.toBeInTheDocument()
  })

  it('deve aplicar estilo primary na tab ativa', () => {
    render(<MobileTabs tabs={mockTabs} />)
    
    const activeButton = screen.getByRole('button', { name: 'Tab 1' })
    expect(activeButton).toHaveClass('ant-btn-primary')
    
    const inactiveButton = screen.getByRole('button', { name: 'Tab 2' })
    expect(inactiveButton).not.toHaveClass('ant-btn-primary')
  })

  it('deve renderizar com tabs vazias sem erro', () => {
    render(<MobileTabs tabs={[]} />)
    
    expect(screen.queryByRole('button')).not.toBeInTheDocument()
  })

  it('deve aplicar tooltip quando fornecido', () => {
    const tabsWithTooltip: Tab[] = [
      {
        title: 'Tab 1',
        name: 'tab1',
        component: <div>Content</div>,
        tooltip: 'Tooltip da Tab 1',
      },
    ]
    
    render(<MobileTabs tabs={tabsWithTooltip} />)
    
    const button = screen.getByRole('button', { name: 'Tab 1' })
    expect(button).toHaveAttribute('title', 'Tooltip da Tab 1')
  })
})
