import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '../../../test/test-utils'
import { DesktopTabs } from '../desktop-tabs'
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

describe('DesktopTabs', () => {
  it('deve renderizar as tabs corretamente', () => {
    render(<DesktopTabs tabs={mockTabs} />)
    
    expect(screen.getByText('Tab 1')).toBeInTheDocument()
    expect(screen.getByText('Tab 2')).toBeInTheDocument()
    expect(screen.getByText('Tab 3')).toBeInTheDocument()
  })

  it('deve mostrar o conteúdo da primeira tab por padrão', () => {
    render(<DesktopTabs tabs={mockTabs} />)
    
    expect(screen.getByTestId('tab1-content')).toBeInTheDocument()
    expect(screen.queryByTestId('tab2-content')).not.toBeInTheDocument()
    expect(screen.queryByTestId('tab3-content')).not.toBeInTheDocument()
  })

  it('deve permitir mudança de tab', () => {
    render(<DesktopTabs tabs={mockTabs} />)
    
    expect(screen.getByText('Tab 2')).toBeInTheDocument()
    
    fireEvent.click(screen.getByText('Tab 2'))
    
    expect(screen.getByTestId('tab2-content')).toBeInTheDocument()
  })

  it('deve chamar onChange quando uma tab é selecionada', () => {
    const onChange = vi.fn()
    render(<DesktopTabs tabs={mockTabs} onChange={onChange} />)
    
    fireEvent.click(screen.getByText('Tab 2'))
    
    expect(onChange).toHaveBeenCalledWith(1)
  })

  it('deve chamar onTabSelect quando uma tab é selecionada', () => {
    const onTabSelect = vi.fn()
    render(<DesktopTabs tabs={mockTabs} onTabSelect={onTabSelect} />)
    
    fireEvent.click(screen.getByText('Tab 2'))
    
    expect(onTabSelect).toHaveBeenCalledWith('tab2')
  })

  it('deve desabilitar tabs marcadas como disabled', () => {
    render(<DesktopTabs tabs={mockTabs} />)
    
    const disabledTab = screen.getByText('Tab 3').closest('.ant-tabs-tab')
    expect(disabledTab).toHaveClass('ant-tabs-tab-disabled')
  })

  it('deve usar defaultTabIndex para definir tab inicial', () => {
    render(<DesktopTabs tabs={mockTabs} defaultTabIndex={1} />)
    
    expect(screen.getByTestId('tab2-content')).toBeInTheDocument()
    expect(screen.queryByTestId('tab1-content')).not.toBeInTheDocument()
  })

  it('deve renderizar com tabs vazias sem erro', () => {
    render(<DesktopTabs tabs={[]} />)
    
    expect(screen.getByRole('tablist')).toBeInTheDocument()
  })
})
