/// <reference types="vite/client" />

declare module '*.png' {
  const src: string;
  export default src;
}

declare module '*.jpg' {
  const src: string;
  export default src;
}

declare module '*.jpeg' {
  const src: string;
  export default src;
}

declare module '*.gif' {
  const src: string;
  export default src;
}

declare module '*.svg' {
  const src: string;
  export default src;
}

declare module 'antd' {
  export const Layout: any;
  export const Grid: any;
  export const Button: any;
  export const message: any;
  export const Card: any;
  export const Form: any;
  export const Checkbox: any;
  export const Row: any;
  export const Col: any;
  export const Typography: any;
  export const Input: any;
  export const InputNumber: any;
  export const Select: any;
  export const DatePicker: any;
  export const Tooltip: any;
  export const Space: any;
  export const Modal: any;
  export const Table: any;
  export const Pagination: any;
  export const Drawer: any;
  export const Collapse: any;
  export const Statistic: any;
  export const Spin: any;
  export const Dropdown: any;
  export const Menu: any;
  export const Badge: any;
  export const Tag: any;
  export const Switch: any;
  export const Radio: any;
  export const Upload: any;
  export const TimePicker: any;
  export const Cascader: any;
  export const TreeSelect: any;
  export const Transfer: any;
  export const Rate: any;
  export const Slider: any;
  export const Progress: any;
  export const Alert: any;
  export const Notification: any;
  export const Popconfirm: any;
  export const Popover: any;
  export const Tooltip: any;
  export const Tabs: any;
  export const Steps: any;
  export const Timeline: any;
  export const Calendar: any;
  export const Empty: any;
  export const Skeleton: any;
  export const Avatar: any;
  export const BackTop: any;
  export const Anchor: any;
  export const Breadcrumb: any;
  export const Affix: any;
  export const ConfigProvider: any;
  export const Divider: any;
  export const List: any;
  export const Descriptions: any;
  export const PageHeader: any;
  export const Result: any;
  export const Statistic: any;
  export const Tree: any;
  export const AutoComplete: any;
  export const Mentions: any;
  export const Segmented: any;
  export const QRCode: any;
  export const Watermark: any;
  export const FloatButton: any;
  export const App: any;
  export const Image: any;
  export const InputRef: any;
  export const theme: any;
  export const version: string;
  const antd: any;
  export default antd;
}

declare module 'antd/es/table' {
  export const Table: any;
  export const Column: any;
  export const ColumnGroup: any;
  export type ColumnType<T = any> = any;
  const table: any;
  export default table;
}
