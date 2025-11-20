import { vi } from 'vitest'

export const mockI18n = {
  t: (key: string, options?: Record<string, string | number>) => {
    const translations: Record<string, string> = {
      'navigation.dashboard': 'Dashboard',
      'navigation.animals': 'Vacas',
      'navigation.reports': 'Relatórios',
      'navigation.suppliers': 'Fornecedores',
      'navigation.sales': 'Vendas',
      'navigation.inventory': 'Estoque',
      'navigation.settings': 'Configurações',
      'navigation.logout': 'Sair',
      
      'animalTable.tabs.animalsList': 'Lista de Animais',
      'animalTable.tabs.milkProduction': 'Produção de Leite',
      'animalTable.tabs.reproduction': 'Reprodução',
      
      'animalTable.reproduction.title': 'Reprodução',
      'animalTable.reproduction.animalName': 'Nome do Animal',
      'animalTable.reproduction.earTag': 'Brinco',
      'animalTable.reproduction.currentPhase': 'Fase Atual',
      'animalTable.reproduction.inseminationDate': 'Data de Inseminação',
      'animalTable.reproduction.pregnancyDate': 'Data de Preñez',
      'animalTable.reproduction.expectedBirthDate': 'Data Esperada de Nascimento',
      'animalTable.reproduction.veterinaryConfirmation': 'Confirmação Veterinária',
      'animalTable.reproduction.actions': 'Ações',
      'animalTable.reproduction.addReproduction': 'Adicionar Reprodução',
      'animalTable.reproduction.updatePhase': 'Atualizar Fase',
      'animalTable.reproduction.delete': 'Excluir',
      'animalTable.reproduction.deleteConfirm': 'Tem certeza que deseja excluir?',
      'animalTable.reproduction.yes': 'Sim',
      'animalTable.reproduction.no': 'Não',
      'animalTable.reproduction.deletedSuccessfully': 'Reprodução excluída com sucesso',
      'animalTable.reproduction.deleteError': 'Erro ao excluir reprodução',
      
      'animalTable.reproduction.phases.insemination': 'Inseminação',
      'animalTable.reproduction.phases.pregnancy': 'Preñez',
      'animalTable.reproduction.phases.birth': 'Nascimento',
      'animalTable.reproduction.phases.lactation': 'Lactação',
      
      'dashboard.title': 'Dashboard',
      'dashboard.totalAnimals': 'Total de Animais',
      'dashboard.totalMilkProduction': 'Produção Total de Leite',
      'dashboard.activeReproductions': 'Reproduções Ativas',
      
      'login.title': 'Login',
      'login.email': 'Email',
      'login.password': 'Senha',
      'login.remember': 'Lembrar-me',
      'login.submit': 'Entrar',
      'login.error': 'Erro ao fazer login',
      
      'validation.required': 'Este campo é obrigatório',
      'validation.email': 'Email inválido',
      'validation.minLength': 'Mínimo de {min} caracteres',
      'validation.maxLength': 'Máximo de {max} caracteres',
      
      'common.loading': 'Carregando...',
      'common.error': 'Erro',
      'common.success': 'Sucesso',
      'common.cancel': 'Cancelar',
      'common.confirm': 'Confirmar',
      'common.save': 'Salvar',
      'common.edit': 'Editar',
      'common.delete': 'Excluir',
      'common.add': 'Adicionar',
      'common.search': 'Buscar',
      'common.filter': 'Filtrar',
      'common.export': 'Exportar',
    }
    
    let translation = translations[key] || key
    
    if (options) {
      Object.keys(options).forEach(optionKey => {
        translation = translation.replace(`{${optionKey}}`, options[optionKey])
      })
    }
    
    return translation
  },
  i18n: {
    changeLanguage: vi.fn(),
    language: 'pt',
    languages: ['pt', 'en', 'es'],
    isInitialized: true,
  },
  ready: true,
}

export const mockUseTranslation = () => mockI18n

export const createMockTranslations = (customTranslations: Record<string, string> = {}) => {
  return {
    t: (key: string, options?: Record<string, string | number>) => {
      const translation = customTranslations[key] || mockI18n.t(key, options)
      return translation
    },
    i18n: mockI18n.i18n,
    ready: true,
  }
}

export const createMockNamespace = (_namespace: string, translations: Record<string, string>) => {
  return {
    t: (key: string, options?: Record<string, string | number>) => {
      const translation = translations[key] || key
      
      if (options) {
        return translation.replace(/\{(\w+)\}/g, (match, key) => options[key] || match)
      }
      
      return translation
    },
    i18n: mockI18n.i18n,
    ready: true,
  }
}
