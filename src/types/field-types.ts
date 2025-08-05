export type FieldType = {
  name: string;
  label: string;
  type: 'text' | 'password' | 'checkbox' | 'link' | 'number';
  placeholder?: string;
  colSpan?: number;
  isRequired?: boolean;
};