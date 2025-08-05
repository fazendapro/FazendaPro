import * as yup from 'yup';
import { AnimalSex } from '../types/type';

export const animalSchema = yup.object().shape({
  animal_name: yup
    .string()
    .trim()
    .min(2, 'O nome deve ter pelo menos 2 caracteres')
    .max(50, 'O nome deve ter no máximo 50 caracteres')
    .required('O nome do animal é obrigatório'),
  ear_tag_number_local: yup
    .number()
    .min(1, 'O número do brinco é obrigatório')
    .required('O número do brinco é obrigatório'),
  type: yup
    .string()
    .oneOf(['vaca', 'bezerro', 'touro', 'novilho'], 'Selecione um tipo válido')
    .required('O tipo do animal é obrigatório'),
  sex: yup
    .number()
    .oneOf([AnimalSex.MALE, AnimalSex.FEMALE], 'Selecione um sexo válido')
    .required('O sexo do animal é obrigatório'),
  breed: yup
    .string()
    .trim()
    .min(2, 'A raça deve ter pelo menos 2 caracteres')
    .max(50, 'A raça deve ter no máximo 50 caracteres')
    .required('A raça do animal é obrigatória'),
  birth_date: yup
    .string()
    .matches(/^\d{4}-\d{2}-\d{2}$/, 'A data deve estar no formato YYYY-MM-DD')
    .required('A data de nascimento é obrigatória'),
  ear_tag_number_register: yup
    .number()
    .min(1, 'O número do brinco de registro é obrigatório')
    .required('O número do brinco de registro é obrigatório'),
}); 