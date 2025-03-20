export type FieldType = {
  name: string;
  label: string;
  type: 'text' | 'password' | 'checkbox' | 'link';
  placeholder?: string;
  colSpan?: number;
  isRequired?: boolean;
};