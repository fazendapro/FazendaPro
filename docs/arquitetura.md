# Arquitetura do Projeto FazendaPro

## Visão Geral

O projeto FazendaPro utiliza uma arquitetura em camadas (3-tier) que separa as responsabilidades em três níveis principais:

1. **Handlers** - Camada de apresentação/API
2. **Services** - Camada de lógica de negócio
3. **Repositories** - Camada de acesso a dados

Esta arquitetura facilita a manutenção, testabilidade e escalabilidade do código, seguindo os princípios SOLID e separação de responsabilidades.

## Fluxo de Dados

O fluxo de uma requisição HTTP segue esta sequência:

```
Request HTTP → Handler → Service → Repository → Database
                ↓         ↓          ↓
              Response ← Response ← Response
```

### Exemplo Visual do Fluxo

```
Cliente (Frontend)
    ↓
HTTP Request (POST /api/v1/animals)
    ↓
Routes (routes.go) - Define rotas e middleware
    ↓
Handler (animal.go) - Recebe e valida requisição
    ↓
Service (animal_service.go) - Aplica regras de negócio
    ↓
Repository (animal_repository.go) - Acessa banco de dados
    ↓
Database (PostgreSQL via GORM)
    ↓
Response HTTP ← Handler formata resposta
```

## Camadas Detalhadas

### 1. Handlers (`internal/api/handlers/`)

**Responsabilidade**: Receber requisições HTTP, validar entrada, formatar saída.

**Características**:
- Recebem `http.Request` e `http.ResponseWriter`
- Fazem parsing de JSON/form data
- Validam parâmetros de entrada
- Chamam métodos dos Services
- Formatam respostas HTTP (JSON)
- Tratam erros HTTP

**Exemplo de Handler**:

```go
type AnimalHandler struct {
    service *service.AnimalService
}

func (h *AnimalHandler) CreateAnimal(w http.ResponseWriter, r *http.Request) {
    // 1. Valida método HTTP
    // 2. Decodifica JSON do request
    // 3. Converte DTO para Model
    // 4. Chama service.CreateAnimal()
    // 5. Formata e retorna resposta HTTP
}
```

**Arquivos principais**:
- `animal.go` - Handlers para animais
- `auth.go` - Handlers de autenticação
- `milk_collection.go` - Handlers para coleta de leite
- `reproduction.go` - Handlers para reprodução
- `sale.go` - Handlers para vendas
- `user.go` - Handlers para usuários
- `farm.go` - Handlers para fazendas
- `debt.go` - Handlers para dívidas

### 2. Services (`internal/service/`)

**Responsabilidade**: Contém a lógica de negócio, validações complexas e orquestra repositories.

**Características**:
- Implementam regras de negócio
- Fazem validações complexas
- Orquestram múltiplos repositories quando necessário
- Não conhecem detalhes de HTTP
- Retornam erros de domínio

**Exemplo de Service**:

```go
type AnimalService struct {
    repository repository.AnimalRepositoryInterface
}

func (s *AnimalService) CreateAnimal(animal *models.Animal) error {
    // 1. Validações de negócio (farm ID obrigatório, etc)
    // 2. Verifica se animal já existe (número de brinca)
    // 3. Define valores padrão (status, timestamps)
    // 4. Chama repository.Create()
    // 5. Retorna erro se algo falhar
}
```

**Arquivos principais**:
- `animal_service.go` - Lógica de negócio para animais
- `user_service.go` - Lógica de negócio para usuários
- `milk_collection_service.go` - Lógica de negócio para coleta de leite
- `reproduction_service.go` - Lógica de negócio para reprodução
- `sale_service.go` - Lógica de negócio para vendas
- `farm_service.go` - Lógica de negócio para fazendas
- `debt_service.go` - Lógica de negócio para dívidas
- `batch_service.go` - Lógica de negócio para lotes

### 3. Repositories (`internal/repository/`)

**Responsabilidade**: Abstrair acesso ao banco de dados, operações CRUD.

**Características**:
- Implementam interfaces definidas em `interfaces.go`
- Usam GORM para acesso ao banco
- Fazem queries SQL através do GORM
- Retornam models ou erros
- Não contêm lógica de negócio

**Exemplo de Repository**:

```go
type AnimalRepository struct {
    db *Database
}

func (r *AnimalRepository) Create(animal *models.Animal) error {
    if err := r.db.DB.Create(animal).Error; err != nil {
        return fmt.Errorf("erro ao criar animal: %w", err)
    }
    return nil
}
```

**Arquivos principais**:
- `animal_repository.go` - Operações de banco para animais
- `user_repository.go` - Operações de banco para usuários
- `milk_collection_repository.go` - Operações de banco para coleta de leite
- `reproduction_repository.go` - Operações de banco para reprodução
- `sale_repository.go` - Operações de banco para vendas
- `farm_repository.go` - Operações de banco para fazendas
- `debt_repository.go` - Operações de banco para dívidas
- `database.go` - Configuração e conexão com banco
- `interfaces.go` - Interfaces que os repositories implementam

## Factories (Padrão Factory)

O projeto utiliza o padrão Factory para criar instâncias de repositories e services de forma centralizada.

### RepositoryFactory (`internal/repository/factory.go`)

Cria instâncias de repositories, recebendo a conexão com o banco de dados:

```go
type RepositoryFactory struct {
    db *Database
}

func (f *RepositoryFactory) CreateAnimalRepository() AnimalRepositoryInterface {
    return NewAnimalRepository(f.db)
}
```

**Vantagens**:
- Centraliza criação de repositories
- Facilita injeção de dependências
- Permite mock em testes

### ServiceFactory (`internal/service/factory.go`)

Cria instâncias de services, recebendo o RepositoryFactory:

```go
type ServiceFactory struct {
    repoFactory *repository.RepositoryFactory
}

func (f *ServiceFactory) CreateAnimalService() *AnimalService {
    animalRepo := f.repoFactory.CreateAnimalRepository()
    return NewAnimalService(animalRepo)
}
```

**Vantagens**:
- Centraliza criação de services
- Gerencia dependências entre services e repositories
- Facilita manutenção

## Exemplo Prático: Fluxo Completo de Criação de Animal

Vamos acompanhar o fluxo completo desde uma requisição HTTP até o banco de dados:

### 1. Requisição HTTP

```http
POST /api/v1/animals
Content-Type: application/json

{
  "farm_id": 1,
  "ear_tag_number_local": 123,
  "animal_name": "Branquinha",
  "sex": 0,
  "breed": "Holandesa",
  "type": "Bovino",
  "animal_type": 1,
  "purpose": 1
}
```

### 2. Routes (`internal/routes/routes.go`)

A rota é definida e conecta o handler:

```go
animalService := serviceFactory.CreateAnimalService()
animalHandler := handlers.NewAnimalHandler(animalService)

r.Route("/animals", func(r chi.Router) {
    r.Use(middleware.Auth(cfg.JWTSecret))
    r.Post("/", animalHandler.CreateAnimal)
})
```

### 3. Handler (`internal/api/handlers/animal.go`)

O handler recebe a requisição:

```go
func (h *AnimalHandler) CreateAnimal(w http.ResponseWriter, r *http.Request) {
    // Decodifica JSON
    var req CreateAnimalRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // Converte DTO para Model
    animal := animalDataToModel(req.AnimalData)
    
    // Chama service
    if err := h.service.CreateAnimal(&animal); err != nil {
        SendErrorResponse(w, "Erro ao criar animal: "+err.Error(), http.StatusBadRequest)
        return
    }
    
    // Retorna resposta
    SendSuccessResponse(w, map[string]interface{}{"id": animal.ID}, 
                       "Animal criado com sucesso", http.StatusCreated)
}
```

### 4. Service (`internal/service/animal_service.go`)

O service aplica regras de negócio:

```go
func (s *AnimalService) CreateAnimal(animal *models.Animal) error {
    // Validações de negócio
    if animal.FarmID == 0 {
        return errors.New("farm ID é obrigatório")
    }
    if animal.EarTagNumberLocal == 0 {
        return errors.New("número da brinca local é obrigatório")
    }
    
    // Verifica se animal já existe
    existingAnimal, err := s.repository.FindByEarTagNumber(animal.FarmID, animal.EarTagNumberLocal)
    if existingAnimal != nil {
        return errors.New("já existe um animal com este número de brinca")
    }
    
    // Define timestamps
    now := time.Now()
    animal.CreatedAt = now
    animal.UpdatedAt = now
    
    // Chama repository
    return s.repository.Create(animal)
}
```

### 5. Repository (`internal/repository/animal_repository.go`)

O repository acessa o banco de dados:

```go
func (r *AnimalRepository) Create(animal *models.Animal) error {
    if err := r.db.DB.Create(animal).Error; err != nil {
        return fmt.Errorf("erro ao criar animal: %w", err)
    }
    return nil
}
```

### 6. Database

O GORM executa o INSERT no PostgreSQL e retorna o resultado.

## Estrutura de Pastas

```
backend/FazendaPro/
├── cmd/
│   └── app/
│       └── app.go              # Estrutura da aplicação
├── internal/
│   ├── api/
│   │   ├── handlers/           # Camada Handler (HTTP)
│   │   │   ├── animal.go
│   │   │   ├── auth.go
│   │   │   └── ...
│   │   └── middleware/         # Middlewares (Auth, CORS)
│   ├── service/                # Camada Service (Lógica de Negócio)
│   │   ├── animal_service.go
│   │   ├── factory.go          # ServiceFactory
│   │   └── ...
│   ├── repository/             # Camada Repository (Acesso a Dados)
│   │   ├── animal_repository.go
│   │   ├── database.go         # Conexão com banco
│   │   ├── factory.go          # RepositoryFactory
│   │   ├── interfaces.go       # Interfaces dos repositories
│   │   └── ...
│   ├── models/                 # Modelos de dados (GORM)
│   │   ├── animal.go
│   │   ├── user.go
│   │   └── ...
│   ├── routes/                 # Definição de rotas
│   │   └── routes.go
│   └── migrations/             # Migrations do banco
│       └── migrations.go
├── main.go                     # Ponto de entrada
└── config/
    └── config.go              # Configurações
```

## Vantagens desta Arquitetura

1. **Separação de Responsabilidades**: Cada camada tem uma responsabilidade clara
2. **Testabilidade**: Fácil criar mocks para testes unitários
3. **Manutenibilidade**: Mudanças em uma camada não afetam outras
4. **Reutilização**: Services podem ser reutilizados em diferentes contextos
5. **Flexibilidade**: Fácil trocar implementações (ex: trocar banco de dados)

## Interfaces e Inversão de Dependência

O projeto utiliza interfaces para desacoplar as camadas:

- **Repositories** implementam interfaces definidas em `interfaces.go`
- **Services** recebem interfaces, não implementações concretas
- Isso permite criar mocks facilmente para testes

**Exemplo**:

```go
// Interface
type AnimalRepositoryInterface interface {
    Create(animal *models.Animal) error
    FindByID(id uint) (*models.Animal, error)
    // ...
}

// Service recebe interface, não implementação concreta
type AnimalService struct {
    repository repository.AnimalRepositoryInterface
}
```

## Middleware

O projeto utiliza middlewares para funcionalidades transversais:

- **Auth Middleware** (`internal/api/middleware/auth.go`): Valida tokens JWT
- **CORS Middleware** (`internal/api/middleware/cors.go`): Gerencia CORS
- **Sentry Middleware**: Captura erros para monitoramento

Os middlewares são aplicados nas rotas através do Chi Router.

## Conclusão

Esta arquitetura em camadas proporciona um código organizado, testável e fácil de manter. Cada camada tem sua responsabilidade bem definida, facilitando o desenvolvimento e a evolução do projeto.

