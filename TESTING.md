# 🧪 Guia de Testes - FazendaPro Frontend

Este documento fornece um guia completo para executar e escrever testes no frontend da aplicação FazendaPro.

## 📋 Índice

- [Configuração](#configuração)
- [Executando Testes](#executando-testes)
- [Estrutura de Testes](#estrutura-de-testes)
- [Escrevendo Testes](#escrevendo-testes)
- [Mocks e Utilitários](#mocks-e-utilitários)
- [Cobertura de Código](#cobertura-de-código)
- [Boas Práticas](#boas-práticas)
- [Troubleshooting](#troubleshooting)

## ⚙️ Configuração

### Dependências Instaladas

O projeto utiliza as seguintes ferramentas de teste:

- **Vitest**: Framework de testes principal
- **React Testing Library**: Para testes de componentes React
- **Jest DOM**: Matchers customizados para DOM
- **@testing-library/user-event**: Para simulação de eventos do usuário

### Arquivos de Configuração

- `vitest.config.ts`: Configuração principal do Vitest
- `coverage.config.js`: Configuração de cobertura de código
- `.eslintrc.test.js`: Configuração do ESLint para testes
- `src/test/setup.ts`: Setup global dos testes

## 🚀 Executando Testes

### Comandos Disponíveis

```bash
# Executar todos os testes
npm run test

# Executar testes em modo watch
npm run test -- --watch

# Executar testes com interface gráfica
npm run test:ui

# Executar testes com cobertura
npm run test:coverage

# Executar testes específicos
npm run test -- --run animals-container.test.tsx

# Executar testes de um diretório
npm run test -- --run src/components
```

### Script de Teste

Use o script personalizado para executar testes:

```bash
./scripts/test.sh
```

## 📁 Estrutura de Testes

```
src/test/
├── setup.ts                    # Configuração global
├── test-utils.tsx              # Utilitários e wrapper customizado
├── mocks/                      # Mocks para dependências
│   ├── antd-mocks.ts          # Mocks do Ant Design
│   ├── router-mocks.ts        # Mocks do React Router
│   ├── i18n-mocks.ts          # Mocks do react-i18next
│   ├── api-mocks.ts           # Mocks da API
│   └── index.ts               # Exportações
├── helpers/                    # Funções auxiliares
│   ├── test-helpers.ts        # Helpers para testes
│   └── index.ts               # Exportações
├── integration/                # Testes de integração
│   ├── animal-table-integration.test.tsx
│   └── index.ts
├── examples/                   # Exemplos de testes
│   └── example.test.tsx
└── README.md                   # Documentação detalhada
```

## ✍️ Escrevendo Testes

### Estrutura Básica

```typescript
import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { customRender } from '../test-utils'
import { MeuComponente } from './meu-componente'

describe('MeuComponente', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('deve renderizar corretamente', () => {
    render(<MeuComponente />)
    
    expect(screen.getByText('Texto esperado')).toBeInTheDocument()
  })

  it('deve responder a interações', () => {
    const handleClick = vi.fn()
    render(<MeuComponente onClick={handleClick} />)
    
    fireEvent.click(screen.getByRole('button'))
    
    expect(handleClick).toHaveBeenCalledTimes(1)
  })
})
```

### Testando Hooks

```typescript
import { describe, it, expect } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { useCounter } from './use-counter'

describe('useCounter', () => {
  it('deve inicializar com valor 0', () => {
    const { result } = renderHook(() => useCounter())
    
    expect(result.current.count).toBe(0)
  })

  it('deve incrementar o contador', () => {
    const { result } = renderHook(() => useCounter())
    
    act(() => {
      result.current.increment()
    })
    
    expect(result.current.count).toBe(1)
  })
})
```

### Testando com Mocks

```typescript
import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { UserProfile } from './user-profile'

// Mock do hook
vi.mock('../hooks/useUser', () => ({
  useUser: () => ({
    user: { name: 'João', email: 'joao@test.com' },
    loading: false,
  }),
}))

describe('UserProfile', () => {
  it('deve exibir informações do usuário', () => {
    render(<UserProfile />)
    
    expect(screen.getByText('João')).toBeInTheDocument()
    expect(screen.getByText('joao@test.com')).toBeInTheDocument()
  })
})
```

## 🎭 Mocks e Utilitários

### Mocks Disponíveis

#### Ant Design
```typescript
import { mockAntd } from '../test/mocks'

mockAntd.message.success('Sucesso!')
mockAntd.message.error('Erro!')
```

#### React Router
```typescript
import { mockRouter } from '../test/mocks'

const mockNavigate = mockRouter.useNavigate()
```

#### API
```typescript
import { mockApiHooks } from '../test/mocks'

const mockUseAnimals = mockApiHooks.useAnimals()
```

#### Tradução
```typescript
import { mockI18n } from '../test/mocks'

const translation = mockI18n.t('test.key')
```

### Utilitários de Teste

```typescript
import { 
  createMockFunction,
  createMockAsyncFunction,
  mockWindowSize,
  mockLocalStorage,
  waitForElement
} from '../test/helpers'

// Criar mock de função
const mockFn = createMockFunction(() => 'result')

// Simular tamanho de tela
mockWindowSize(1024, 768)

// Simular localStorage
const mockStorage = mockLocalStorage()

// Aguardar elemento
const element = await waitForElement('.my-selector')
```

## 📊 Cobertura de Código

### Metas de Cobertura

- **Statements**: 70%
- **Branches**: 70%
- **Functions**: 70%
- **Lines**: 70%

### Visualizando Cobertura

```bash
npm run test:coverage
```

O relatório será gerado em `coverage/index.html`.

### Configuração de Cobertura

O arquivo `coverage.config.js` define:

- Exclusões de arquivos
- Thresholds mínimos
- Formatos de relatório
- Diretório de saída

## ✅ Boas Práticas

### 1. Nomenclatura

- **Arquivos**: `[nome].test.tsx` ou `[nome].spec.tsx`
- **Descrições**: Em português, usando "deve"
- **Casos**: Descritivos e específicos

### 2. Estrutura de Testes

```typescript
describe('Componente', () => {
  beforeEach(() => {
    // Setup comum
  })

  describe('quando renderizado', () => {
    it('deve mostrar conteúdo inicial', () => {
      // Teste
    })
  })

  describe('quando interagido', () => {
    it('deve responder corretamente', () => {
      // Teste
    })
  })
})
```

### 3. Queries Acessíveis

```typescript
// ❌ Ruim
screen.getByClassName('button-primary')

// ✅ Bom
screen.getByRole('button', { name: 'Salvar' })
```

### 4. Teste de Estados

```typescript
it('deve mostrar loading', () => {
  render(<Component loading={true} />)
  expect(screen.getByText('Carregando...')).toBeInTheDocument()
})

it('deve mostrar erro', () => {
  render(<Component error="Erro ao carregar" />)
  expect(screen.getByText('Erro ao carregar')).toBeInTheDocument()
})
```

### 5. Limpeza de Mocks

```typescript
beforeEach(() => {
  vi.clearAllMocks()
})
```

## 🔧 Troubleshooting

### Problemas Comuns

#### 1. Erro de Importação
```
Cannot find module '@testing-library/jest-dom'
```
**Solução**: Verifique se as dependências estão instaladas:
```bash
npm install
```

#### 2. Mock Não Funciona
```
Mock function was not called
```
**Solução**: Certifique-se de que o mock está no escopo correto e use `vi.clearAllMocks()` no `beforeEach`.

#### 3. Elemento Não Encontrado
```
Unable to find an element
```
**Solução**: Use `waitFor` para elementos assíncronos:
```typescript
await waitFor(() => {
  expect(screen.getByText('Texto')).toBeInTheDocument()
})
```

#### 4. Teste Falha Intermitentemente
**Solução**: Use `waitFor` e evite dependências de timing.

### Debug de Testes

```typescript
// Imprimir HTML renderizado
screen.debug()

// Imprimir elemento específico
screen.debug(screen.getByRole('button'))

// Aguardar elemento aparecer
await waitFor(() => {
  expect(screen.getByText('Texto')).toBeInTheDocument()
})
```

## 📚 Recursos Adicionais

- [Documentação do Vitest](https://vitest.dev/)
- [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/)
- [Jest DOM Matchers](https://github.com/testing-library/jest-dom)
- [Testing Best Practices](https://kentcdodds.com/blog/common-mistakes-with-react-testing-library)

## 🤝 Contribuindo

Ao adicionar novos testes:

1. Siga as convenções de nomenclatura
2. Escreva testes descritivos
3. Cubra casos de sucesso e erro
4. Use mocks apropriados
5. Mantenha os testes simples e focados
6. Atualize esta documentação se necessário

---

**Última atualização**: Dezembro 2024
