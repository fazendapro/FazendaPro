import { vi } from 'vitest'
import React from 'react'

export const mockAntd = {
  message: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
  notification: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
  },
  modal: {
    confirm: vi.fn(),
    info: vi.fn(),
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
  },
}

interface TableColumn {
  title?: string
  dataIndex?: string
  render?: (value: unknown, record: Record<string, unknown>) => React.ReactNode
}

interface TabItem {
  label?: string
  children?: React.ReactNode
  disabled?: boolean
}

export const mockAntdComponents = {
  Button: ({ children, onClick, ...props }: { children?: React.ReactNode; onClick?: () => void; [key: string]: unknown }) => 
    React.createElement('button', { onClick, ...props }, children),
  
  Input: ({ onChange, value, ...props }: { onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void; value?: string; [key: string]: unknown }) => 
    React.createElement('input', { onChange, value, ...props }),
  
  Table: ({ dataSource, columns, ...props }: { dataSource?: Record<string, unknown>[]; columns?: TableColumn[]; [key: string]: unknown }) => 
    React.createElement('table', props,
      React.createElement('thead', null,
        React.createElement('tr', null,
          columns?.map((col: TableColumn, index: number) => 
            React.createElement('th', { key: index }, col.title)
          )
        )
      ),
      React.createElement('tbody', null,
        dataSource?.map((row: Record<string, unknown>, rowIndex: number) => 
          React.createElement('tr', { key: rowIndex },
            columns?.map((col: TableColumn, colIndex: number) => 
              React.createElement('td', { key: colIndex },
                col.render ? (col.render(row[col.dataIndex || ''], row) as React.ReactNode) : (row[col.dataIndex || ''] as React.ReactNode)
              )
            )
          )
        )
      )
    ),
  
  Tabs: ({ items, activeKey, onChange, ...props }: { items?: TabItem[]; activeKey?: string; onChange?: (key: string) => void; [key: string]: unknown }) => 
    React.createElement('div', props,
      React.createElement('div', { role: 'tablist' },
        items?.map((item: TabItem, index: number) => 
          React.createElement('button', {
            key: index,
            role: 'tab',
            'aria-selected': activeKey === index.toString(),
            onClick: () => onChange?.(index.toString()),
            disabled: item.disabled
          }, item.label)
        )
      ),
      React.createElement('div', null,
        items?.[parseInt(activeKey || '0')]?.children
      )
    ),
  
  Tag: ({ children, color, ...props }: { children?: React.ReactNode; color?: string; [key: string]: unknown }) => 
    React.createElement('span', { style: { backgroundColor: color }, ...props }, children),
  
  Space: ({ children, ...props }: { children?: React.ReactNode; [key: string]: unknown }) => 
    React.createElement('div', { style: { display: 'flex', gap: '8px' }, ...props }, children),

  Popconfirm: ({ children, onConfirm, title, ...props }: { children?: React.ReactNode; onConfirm?: () => void; title?: string; [key: string]: unknown }) => 
    React.createElement('div', props,
      children,
      React.createElement('div', { 'data-testid': 'popconfirm' },
        React.createElement('p', null, title),
        React.createElement('button', { onClick: onConfirm }, 'Confirm')
      )
    ),
}