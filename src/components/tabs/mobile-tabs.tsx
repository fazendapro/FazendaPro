import React, { useState } from 'react'
import { Button, Space } from 'antd'
import { Tab } from './types'

interface MobileTabsProps {
  tabs: Tab[]
  defaultTabIndex?: number
  onChange?: (index: number) => void
  onTabSelect?: (name: string) => void
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
      <Space direction="vertical" style={{ width: '100%', marginBottom: 16 }}>
        {tabs.map((tab, index) => (
          <Button
            key={tab.name}
            type={selectedIndex === index ? 'primary' : 'default'}
            disabled={tab.isDisabled}
            onClick={() => handleTabSelect(index)}
            title={tab.tooltip}
            block
          >
            {tab.title}
          </Button>
        ))}
      </Space>
      <div>
        {tabs[selectedIndex]?.component}
      </div>
    </div>
  )
}
