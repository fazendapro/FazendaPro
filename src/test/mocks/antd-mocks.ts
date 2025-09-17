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

export const mockAntdComponents = {
  Button: ({ children, onClick, ...props }: any) => 
    React.createElement('button', { onClick, ...props }, children),
  
  Input: ({ onChange, value, ...props }: any) => 
    React.createElement('input', { onChange, value, ...props }),
  
  Table: ({ dataSource, columns, ...props }: any) => 
    React.createElement('table', props,
      React.createElement('thead', null,
        React.createElement('tr', null,
          columns?.map((col: any, index: number) => 
            React.createElement('th', { key: index }, col.title)
          )
        )
      ),
      React.createElement('tbody', null,
        dataSource?.map((row: any, rowIndex: number) => 
          React.createElement('tr', { key: rowIndex },
            columns?.map((col: any, colIndex: number) => 
              React.createElement('td', { key: colIndex },
                col.render ? col.render(row[col.dataIndex], row) : row[col.dataIndex]
              )
            )
          )
        )
      )
    ),
  
  Tabs: ({ items, activeKey, onChange, ...props }: any) => 
    React.createElement('div', props,
      React.createElement('div', { role: 'tablist' },
        items?.map((item: any, index: number) => 
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
  
  Tag: ({ children, color, ...props }: any) => 
    React.createElement('span', { style: { backgroundColor: color }, ...props }, children),
  
  Space: ({ children, ...props }: any) => 
    React.createElement('div', { style: { display: 'flex', gap: '8px' }, ...props }, children),

  Popconfirm: ({ children, onConfirm, title, ...props }: any) => 
    React.createElement('div', props,
      children,
      React.createElement('div', { 'data-testid': 'popconfirm' },
        React.createElement('p', null, title),
        React.createElement('button', { onClick: onConfirm }, 'Confirm')
      )
    ),
}