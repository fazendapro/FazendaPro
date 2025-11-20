export interface Tab {
  title: string
  name: string
  component: React.ReactNode
  isDisabled?: boolean
  tooltip?: string
}
