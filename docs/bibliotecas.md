# Bibliotecas e Dependências do Projeto FazendaPro

Este documento descreve todas as bibliotecas utilizadas no projeto FazendaPro, explicando o propósito e uso de cada uma delas.

## Índice

- [Backend (Go)](#backend-go)
  - [Dependências Principais](#dependências-principais)
  - [Dependências Indiretas](#dependências-indiretas)
- [Frontend (React/TypeScript)](#frontend-reacttypescript)
  - [Dependências Principais](#dependências-principais-1)
  - [Dependências de Desenvolvimento](#dependências-de-desenvolvimento)

---

## Backend (Go)

### Dependências Principais

#### 1. `github.com/go-chi/chi/v5` v5.2.2
**Propósito**: Router HTTP leve e rápido para Go.

**Uso no Projeto**:
- Define todas as rotas da API REST
- Gerencia middlewares (autenticação, CORS, Sentry)
- Organiza rotas em grupos (ex: `/api/v1/animals`, `/api/v1/sales`)

**Exemplo de uso**:
```go
r := chi.NewRouter()
r.Route("/api/v1/animals", func(r chi.Router) {
    r.Use(middleware.Auth(cfg.JWTSecret))
    r.Post("/", animalHandler.CreateAnimal)
})
```

**Por que foi escolhido**: Chi é um router minimalista, performático e fácil de usar, ideal para APIs REST.

---

#### 2. `github.com/golang-jwt/jwt/v5` v5.2.2
**Propósito**: Biblioteca para criação e validação de tokens JWT (JSON Web Tokens).

**Uso no Projeto**:
- Geração de tokens de autenticação após login
- Validação de tokens em requisições protegidas
- Criação de refresh tokens para renovação de sessão

**Onde é usado**:
- `internal/api/handlers/auth.go` - Criação de tokens
- `internal/api/middleware/auth.go` - Validação de tokens

**Por que foi escolhido**: Biblioteca oficial e amplamente utilizada para JWT em Go.

---

#### 3. `github.com/google/uuid` v1.6.0
**Propósito**: Geração de UUIDs (Identificadores Únicos Universais).

**Uso no Projeto**:
- Geração de IDs únicos para refresh tokens
- Identificadores únicos em geral quando necessário

**Por que foi escolhido**: Biblioteca padrão do Google, confiável e amplamente adotada.

---

#### 4. `github.com/joho/godotenv` v1.5.1
**Propósito**: Carregamento de variáveis de ambiente a partir de arquivos `.env`.

**Uso no Projeto**:
- Carrega configurações de `.env` e `.env.production`
- Permite diferentes configurações por ambiente (dev, production)

**Onde é usado**:
- `config/config.go` - Função `Load()` carrega variáveis de ambiente

**Exemplo**:
```go
godotenv.Load()           // Carrega .env
godotenv.Load(".env.production")  // Carrega .env.production
```

**Por que foi escolhido**: Solução simples e popular para gerenciamento de configurações.

---

#### 5. `github.com/stretchr/testify` v1.11.1
**Propósito**: Biblioteca de testes com assertions, mocks e suíte de testes.

**Uso no Projeto**:
- Testes unitários e de integração
- Assertions para validação de resultados
- Mocks para isolar dependências em testes

**Onde é usado**:
- Todos os arquivos em `tests/` (handlers, services, repositories)

**Exemplo**:
```go
assert.Equal(t, expected, actual)
assert.NoError(t, err)
```

**Por que foi escolhido**: Facilita escrita de testes com assertions legíveis e mocks.

---

#### 6. `golang.org/x/crypto` v0.42.0
**Propósito**: Pacote de criptografia da equipe Go, incluindo funções de hash.

**Uso no Projeto**:
- Hash de senhas usando bcrypt
- Comparação de senhas para autenticação

**Onde é usado**:
- `internal/utils/password.go` - Funções de hash e verificação de senhas

**Exemplo**:
```go
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

**Por que foi escolhido**: Biblioteca oficial do Go para operações criptográficas.

---

#### 7. `gorm.io/gorm` v1.30.0
**Propósito**: ORM (Object-Relational Mapping) para Go, facilita interação com bancos de dados.

**Uso no Projeto**:
- Abstração de acesso ao banco de dados
- Migrations automáticas
- Queries type-safe
- Relacionamentos entre modelos

**Onde é usado**:
- Todos os repositories (`internal/repository/`)
- Migrations (`internal/migrations/migrations.go`)
- Models (`internal/models/`)

**Exemplo**:
```go
db.Create(&animal)
db.Where("farm_id = ?", farmID).Find(&animals)
```

**Por que foi escolhido**: ORM mais popular em Go, com excelente suporte a migrations e relacionamentos.

---

#### 8. `gorm.io/driver/postgres` v1.6.0
**Propósito**: Driver PostgreSQL para GORM.

**Uso no Projeto**:
- Conexão com banco de dados PostgreSQL em produção
- Driver principal usado no projeto

**Onde é usado**:
- `internal/repository/database.go` - Configuração de conexão

**Por que foi escolhido**: Driver oficial do GORM para PostgreSQL.

---

#### 9. `gorm.io/driver/sqlite` v1.6.0
**Propósito**: Driver SQLite para GORM.

**Uso no Projeto**:
- Usado principalmente em testes
- Permite testes sem necessidade de banco PostgreSQL

**Por que foi escolhido**: Facilita testes locais sem dependência de banco externo.

---

### Dependências Indiretas (Transitivas)

#### `github.com/getsentry/sentry-go` v0.37.0
**Propósito**: Integração com Sentry para monitoramento de erros e performance.

**Uso no Projeto**:
- Captura e reporta erros em produção
- Middleware HTTP para capturar panics
- Monitoramento de performance

**Onde é usado**:
- `main.go` - Inicialização do Sentry
- `internal/routes/routes.go` - Middleware Sentry nas rotas

**Exemplo**:
```go
sentry.Init(sentry.ClientOptions{Dsn: cfg.SentryDSN})
sentryHandler := sentryhttp.New(sentryhttp.Options{Repanic: true})
```

**Por que foi escolhido**: Ferramenta profissional de monitoramento de erros.

---

#### `github.com/jackc/pgx/v5` v5.6.0
**Propósito**: Driver PostgreSQL de alto desempenho (usado pelo GORM internamente).

**Uso no Projeto**: Usado indiretamente pelo GORM para comunicação com PostgreSQL.

---

#### Outras Dependências Indiretas

- **`github.com/jinzhu/inflection`** - Pluralização de nomes para GORM
- **`github.com/jinzhu/now`** - Utilitários de data/hora
- **`golang.org/x/sys`** - Interfaces de sistema
- **`golang.org/x/net`** - Funcionalidades de rede
- **`golang.org/x/text`** - Processamento de texto
- **`golang.org/x/tools`** - Ferramentas de desenvolvimento

---

## Frontend (React/TypeScript)

### Dependências Principais

#### 1. `react` ^19.0.0 e `react-dom` ^19.0.0
**Propósito**: Biblioteca principal para construção de interfaces de usuário.

**Uso no Projeto**:
- Base de toda a aplicação frontend
- Componentes React para todas as páginas e componentes

**Por que foi escolhido**: Framework padrão para desenvolvimento web moderno.

---

#### 2. `react-router-dom` ^7.4.0
**Propósito**: Roteamento client-side para aplicações React.

**Uso no Projeto**:
- Navegação entre páginas
- Rotas protegidas (requerem autenticação)
- Gerenciamento de histórico de navegação

**Exemplo de uso**:
```tsx
<Routes>
  <Route path="/login" element={<Login />} />
  <Route path="/animals" element={<Animals />} />
</Routes>
```

**Por que foi escolhido**: Solução padrão e robusta para roteamento em React.

---

#### 3. `antd` ^5.24.2
**Propósito**: Biblioteca de componentes UI (Ant Design) para React.

**Uso no Projeto**:
- Componentes de interface (tabelas, formulários, modais, etc.)
- Design system consistente
- Componentes prontos: Table, Form, Modal, DatePicker, etc.

**Exemplo de uso**:
```tsx
import { Table, Button, Modal } from 'antd';
```

**Por que foi escolhido**: Biblioteca completa de componentes com design profissional.

---

#### 4. `axios` ^1.8.4
**Propósito**: Cliente HTTP para fazer requisições à API.

**Uso no Projeto**:
- Todas as chamadas à API backend
- Interceptors para adicionar tokens de autenticação
- Tratamento de erros HTTP

**Onde é usado**:
- `src/components/services/` - Serviços de API
- `src/config/api.ts` - Configuração do cliente Axios

**Exemplo**:
```tsx
axios.get('/api/v1/animals', { headers: { Authorization: `Bearer ${token}` } })
```

**Por que foi escolhido**: Cliente HTTP popular, com suporte a interceptors e cancelamento de requisições.

---

#### 5. `react-hook-form` ^7.54.2
**Propósito**: Biblioteca para gerenciamento de formulários com hooks.

**Uso no Projeto**:
- Formulários de criação/edição (animais, vendas, etc.)
- Validação de campos
- Performance otimizada (re-renderização mínima)

**Exemplo de uso**:
```tsx
const { register, handleSubmit, formState: { errors } } = useForm();
```

**Por que foi escolhido**: Solução moderna e performática para formulários em React.

---

#### 6. `yup` ^1.6.1
**Propósito**: Biblioteca de validação de esquemas.

**Uso no Projeto**:
- Validação de formulários
- Integração com react-hook-form via `@hookform/resolvers`

**Exemplo**:
```tsx
const schema = yup.object({
  name: yup.string().required('Nome é obrigatório'),
  email: yup.string().email('Email inválido')
});
```

**Por que foi escolhido**: Validação declarativa e poderosa, integrada com react-hook-form.

---

#### 7. `@hookform/resolvers` ^4.1.3
**Propósito**: Resolvers para integrar react-hook-form com bibliotecas de validação (Yup, Zod, etc.).

**Uso no Projeto**:
- Conecta Yup com react-hook-form

**Exemplo**:
```tsx
const { register, handleSubmit } = useForm({
  resolver: yupResolver(schema)
});
```

---

#### 8. `jwt-decode` ^4.0.0
**Propósito**: Decodificação de tokens JWT no cliente.

**Uso no Projeto**:
- Extração de informações do token (user ID, expiração, etc.)
- Verificação de expiração do token
- Obtenção de dados do usuário autenticado

**Onde é usado**:
- Contextos de autenticação
- Middleware de rotas protegidas

**Por que foi escolhido**: Biblioteca simples e leve para decodificar JWT.

---

#### 9. `react-toastify` ^11.0.5
**Propósito**: Sistema de notificações toast (mensagens temporárias).

**Uso no Projeto**:
- Mensagens de sucesso/erro
- Feedback visual para ações do usuário

**Exemplo**:
```tsx
toast.success('Animal criado com sucesso!');
toast.error('Erro ao criar animal');
```

**Por que foi escolhido**: Solução simples e elegante para notificações.

---

#### 10. `i18next` ^25.4.1 e `react-i18next` ^15.4.1
**Propósito**: Internacionalização (i18n) - suporte a múltiplos idiomas.

**Uso no Projeto**:
- Tradução de textos da interface
- Suporte a português, inglês e espanhol
- Arquivos de tradução em `src/locale/`

**Exemplo**:
```tsx
const { t } = useTranslation();
<h1>{t('welcome')}</h1>
```

**Por que foi escolhido**: Solução padrão e robusta para i18n em React.

---

#### 11. `chart.js` ^4.5.0 e `react-chartjs-2` ^5.3.0
**Propósito**: Biblioteca de gráficos e visualização de dados.

**Uso no Projeto**:
- Gráficos de produção de leite
- Estatísticas de vendas
- Dashboards com visualizações

**Exemplo**:
```tsx
import { Line } from 'react-chartjs-2';
<Line data={chartData} />
```

**Por que foi escolhido**: Biblioteca popular e flexível para gráficos.

---

#### 12. `@types/chart.js` ^2.9.41
**Propósito**: Tipos TypeScript para Chart.js.

**Uso no Projeto**: Suporte a tipos para Chart.js em TypeScript.

---

#### 13. `jspdf` ^3.0.3 e `jspdf-autotable` ^5.0.2
**Propósito**: Geração de PDFs no navegador.

**Uso no Projeto**:
- Exportação de relatórios em PDF
- Geração de tabelas em PDF
- Histórico de animais em PDF

**Onde é usado**:
- Componentes de exportação (ex: `AnimalHistoryExport`)

**Por que foi escolhido**: Solução completa para geração de PDFs client-side.

---

#### 14. `serve` ^14.2.4
**Propósito**: Servidor HTTP simples para servir arquivos estáticos.

**Uso no Projeto**:
- Servir build de produção localmente
- Script `start` no package.json

**Por que foi escolhido**: Ferramenta simples para servir builds estáticos.

---

### Dependências de Desenvolvimento

#### 1. `vite` ^6.2.0
**Propósito**: Build tool e dev server extremamente rápido.

**Uso no Projeto**:
- Servidor de desenvolvimento
- Build de produção
- Hot Module Replacement (HMR)

**Por que foi escolhido**: Alternativa moderna e muito mais rápida que Webpack.

---

#### 2. `@vitejs/plugin-react` ^4.3.4
**Propósito**: Plugin do Vite para suporte a React.

**Uso no Projeto**: Habilita JSX, Fast Refresh, etc. no Vite.

---

#### 3. `typescript` ~5.7.2
**Propósito**: Superset do JavaScript com tipagem estática.

**Uso no Projeto**:
- Todo o código frontend é escrito em TypeScript
- Tipagem de props, estados, funções, etc.

**Por que foi escolhido**: Adiciona segurança de tipos e melhor DX (Developer Experience).

---

#### 4. `vitest` ^1.3.1
**Propósito**: Framework de testes rápido e moderno.

**Uso no Projeto**:
- Testes unitários e de integração
- Compatível com Jest, mas mais rápido

**Por que foi escolhido**: Integração nativa com Vite, muito rápido.

---

#### 5. `@vitest/coverage-v8` ^1.6.1
**Propósito**: Plugin para cobertura de código nos testes.

**Uso no Projeto**: Gera relatórios de cobertura de testes.

---

#### 6. `@vitest/ui` ^1.3.1
**Propósito**: Interface web para visualizar e executar testes.

**Uso no Projeto**: UI interativa para testes (comando `test:ui`).

---

#### 7. `@testing-library/react` ^16.0.0
**Propósito**: Utilitários para testar componentes React.

**Uso no Projeto**:
- Renderização de componentes em testes
- Queries para encontrar elementos
- Simulação de interações do usuário

**Exemplo**:
```tsx
import { render, screen } from '@testing-library/react';
render(<Component />);
expect(screen.getByText('Hello')).toBeInTheDocument();
```

**Por que foi escolhido**: Biblioteca padrão para testes de componentes React.

---

#### 8. `@testing-library/jest-dom` ^6.4.2
**Propósito**: Matchers customizados do Jest para DOM.

**Uso no Projeto**: Assertions como `toBeInTheDocument()`, `toHaveClass()`, etc.

---

#### 9. `@testing-library/user-event` ^14.5.2
**Propósito**: Simulação de eventos de usuário de forma realista.

**Uso no Projeto**: Simular cliques, digitação, etc. em testes.

---

#### 10. `eslint` ^9.21.0
**Propósito**: Linter para JavaScript/TypeScript.

**Uso no Projeto**:
- Verificação de qualidade de código
- Padrões de código consistentes
- Detecção de erros comuns

---

#### 11. `typescript-eslint` ^8.24.1
**Propósito**: Regras ESLint específicas para TypeScript.

**Uso no Projeto**: Linting de código TypeScript.

---

#### 12. `eslint-plugin-react-hooks` ^5.1.0
**Propósito**: Regras ESLint para hooks do React.

**Uso no Projeto**: Validação de uso correto de hooks (ex: dependências do useEffect).

---

#### 13. `eslint-plugin-react-refresh` ^0.4.19
**Propósito**: Regras para Fast Refresh do React.

**Uso no Projeto**: Garante compatibilidade com Fast Refresh.

---

#### 14. `@types/react` ^19.1.8 e `@types/react-dom` ^19.1.6
**Propósito**: Tipos TypeScript para React e React DOM.

**Uso no Projeto**: Tipagem para bibliotecas React.

---

#### 15. `@types/node` ^24.2.0
**Propósito**: Tipos TypeScript para Node.js.

**Uso no Projeto**: Tipagem para APIs do Node.js (usado em configurações do Vite).

---

#### 16. `@types/jest` ^29.5.12
**Propósito**: Tipos TypeScript para Jest (usado com Vitest).

---

#### 17. `jest` ^29.7.0 e `jest-environment-jsdom` ^29.7.0
**Propósito**: Framework de testes (usado como fallback ou em alguns testes).

**Nota**: O projeto usa principalmente Vitest, mas mantém Jest para compatibilidade.

---

#### 18. `ts-jest` ^29.1.2
**Propósito**: Processador TypeScript para Jest.

---

#### 19. `globals` ^15.15.0
**Propósito**: Configuração de globals para ESLint.

---

## Resumo por Categoria

### Backend - Funcionalidades Principais
- **Router**: Chi
- **ORM**: GORM
- **Autenticação**: JWT
- **Banco de Dados**: PostgreSQL (driver pgx)
- **Testes**: Testify
- **Monitoramento**: Sentry
- **Configuração**: godotenv

### Frontend - Funcionalidades Principais
- **Framework**: React 19
- **Roteamento**: React Router
- **UI Components**: Ant Design
- **HTTP Client**: Axios
- **Formulários**: React Hook Form + Yup
- **Gráficos**: Chart.js
- **PDF**: jsPDF
- **i18n**: i18next
- **Build Tool**: Vite
- **Testes**: Vitest + Testing Library
- **Linting**: ESLint

## Gerenciamento de Versões

### Backend (Go)
- Versões fixas no `go.mod` garantem builds reproduzíveis
- Comando `go mod tidy` mantém dependências atualizadas

### Frontend (Node.js)
- Versões com `^` permitem atualizações de patch e minor
- Comando `npm install` instala dependências
- `npm update` atualiza para versões compatíveis

## Atualização de Dependências

### Backend
```bash
go get -u ./...          # Atualiza todas as dependências
go mod tidy              # Limpa dependências não usadas
```

### Frontend
```bash
npm update               # Atualiza dentro do range permitido
npm install package@latest  # Atualiza pacote específico
```

## Considerações de Segurança

- **Backend**: Dependências são verificadas pelo Go modules
- **Frontend**: Use `npm audit` para verificar vulnerabilidades
- **Atualizações**: Mantenha dependências atualizadas para correções de segurança

## Conclusão

O projeto FazendaPro utiliza uma stack moderna e bem estabelecida:

- **Backend**: Go com GORM, Chi, JWT - stack performática e robusta
- **Frontend**: React 19 com TypeScript, Vite, Ant Design - stack moderna e produtiva

Todas as bibliotecas foram escolhidas considerando:
- Popularidade e comunidade ativa
- Performance
- Facilidade de uso
- Manutenibilidade
- Suporte a longo prazo

