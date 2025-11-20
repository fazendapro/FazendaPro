import React, { useState } from 'react'
import { Tabs } from 'antd'
import { Tab } from './types'

interface DesktopTabsProps {
  tabs: Tab[]
  defaultTabIndex?: number
  onChange?: (index: number) => void
  tabListMx?: string
  onTabSelect?: (name: string) => void
}

export const DesktopTabs: React.FC<DesktopTabsProps> = ({
  tabs,
  defaultTabIndex = 0,
  onChange,
  onTabSelect
}) => {
  const [selectedIndex, setSelectedIndex] = useState(defaultTabIndex)

  const handleTabChange = (key: string) => {
    const index = parseInt(key)
    setSelectedIndex(index)
    onChange?.(index)
    onTabSelect?.(tabs[index]?.name)
  }

  const tabItems = tabs.map((tab, index) => ({
    key: index.toString(),
    label: tab.title,
    disabled: tab.isDisabled,
    children: tab.component
  }))

  return (
    <Tabs
      activeKey={selectedIndex.toString()}
      onChange={handleTabChange}
      items={tabItems}
      type="card"
    />
  )
}
