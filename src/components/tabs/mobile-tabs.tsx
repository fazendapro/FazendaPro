import React, { useState } from 'react'
import { Button } from 'antd'
import { 
  UnorderedListOutlined, 
  ExperimentOutlined, 
  HeartOutlined 
} from '@ant-design/icons'
import { Tab } from './types'

interface MobileTabsProps {
  tabs: Tab[]
  defaultTabIndex?: number
  onChange?: (index: number) => void
  onTabSelect?: (name: string) => void
}

const getTabIcon = (tabName: string) => {
  switch (tabName) {
    case 'animals-list':
      return <UnorderedListOutlined />
    case 'milk-production':
      return <ExperimentOutlined />
    case 'reproduction':
      return <HeartOutlined />
    default:
      return <UnorderedListOutlined />
  }
}

export const MobileTabs: React.FC<MobileTabsProps> = ({
  tabs,
  defaultTabIndex = 0,
  onChange,
  onTabSelect
}) => {
  const [selectedIndex, setSelectedIndex] = useState(defaultTabIndex)

  const handleTabSelect = (index: number) => {
    setSelectedIndex(index)
    onChange?.(index)
    onTabSelect?.(tabs[index]?.name)
  }

  return (
    <div>
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-around', 
        marginBottom: 16,
        padding: '8px 0',
        backgroundColor: '#fafafa',
        borderRadius: '8px',
        border: '1px solid #f0f0f0'
      }}>
        {tabs.map((tab, index) => (
            <Button
              key={`tab-${tab.name}-${index}`}
              type={selectedIndex === index ? 'primary' : 'text'}
              disabled={tab.isDisabled}
              onClick={() => handleTabSelect(index)}
              icon={getTabIcon(tab.name)}
              size="large"
              style={{
                minWidth: '60px',
                height: '50px',
                borderRadius: '8px',
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                gap: '4px'
              }}
            >
              <span style={{ 
                fontSize: '10px', 
                fontWeight: selectedIndex === index ? 'bold' : 'normal',
                lineHeight: 1
              }}>
                {tab.title.split(' ')[0]}
              </span>
            </Button>
        ))}
      </div>

      <div>
        {tabs[selectedIndex]?.component}
      </div>
    </div>
  )
}
