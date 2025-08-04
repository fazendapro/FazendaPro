import * as yup from 'yup';

export const animalSchema = yup.object().shape({
  animalName: yup
    .string()
    .min(2, 'O nome deve ter pelo menos 2 caracteres')
    .max(50, 'O nome deve ter no máximo 50 caracteres')
    .required('O nome do animal é obrigatório'),
  earringNumber: yup
    .string()
    .min(1, 'O número do brinco é obrigatório')
    .max(20, 'O número do brinco deve ter no máximo 20 caracteres')
    .required('O número do brinco é obrigatório'),
  type: yup
    .string()
    .oneOf(['vaca', 'bezerro', 'touro', 'novilho'], 'Selecione um tipo válido')
    .required('O tipo do animal é obrigatório'),
  sex: yup
    .string()
    .oneOf(['macho', 'fêmea'], 'Selecione um sexo válido')
    .required('O sexo do animal é obrigatório'),
  breed: yup
    .string()
    .min(2, 'A raça deve ter pelo menos 2 caracteres')
    .max(50, 'A raça deve ter no máximo 50 caracteres')
    .required('A raça do animal é obrigatória'),
  birthDate: yup
    .string()
    .matches(/^\d{4}-\d{2}-\d{2}$/, 'A data deve estar no formato YYYY-MM-DD')
    .required('A data de nascimento é obrigatória'),
  earringNumberGlobal: yup
    .string()
    .min(1, 'O número do brinco de registro é obrigatório')
    .max(20, 'O número do brinco de registro deve ter no máximo 20 caracteres')
    .required('O número do brinco de registro é obrigatório'),
}); 