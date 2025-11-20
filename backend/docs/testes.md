# Documentação de Testes - FazendaPro

## Visão Geral

O projeto FazendaPro utiliza testes unitários e de integração para garantir a qualidade e confiabilidade do código. Os testes são escritos em Go usando a biblioteca `testify` para assertions e mocks.

## Estrutura de Testes

A estrutura de testes segue a organização do código principal:

```
tests/
├── handlers/          # Testes dos handlers HTTP
├── services/          # Testes dos services (lógica de negócio)
├── repositories/      # Testes dos repositories (acesso a dados)
├── middleware/       # Testes dos middlewares
├── models/           # Testes dos models
├── utils/            # Testes de utilitários
├── integration/      # Testes de integração
├── mocks/            # Mocks de repositories
└── integration_test.go  # Testes de integração principais
```

## Tipos de Testes

### 1. Testes Unitários

Testam componentes isolados, usando mocks para dependências.

**Localização**: `tests/handlers/`, `tests/services/`, `tests/repositories/`

**Características**:
- Rápidos de executar
- Isolados (não dependem de banco de dados real)
- Usam mocks para simular dependências
- Focam em uma única unidade de código

### 2. Testes de Integração

Testam o sistema completo, incluindo banco de dados e rotas.

**Localização**: `tests/integration/`, `tests/integration_test.go`

**Características**:
- Mais lentos (usam banco de dados)
- Testam fluxos completos
- Usam banco SQLite em memória
- Testam interação entre componentes

## Bibliotecas Utilizadas

### testify

Biblioteca principal para testes em Go:

- **assert**: Assertions para validação de resultados
- **require**: Assertions que param o teste em caso de falha
- **mock**: Criação de mocks para dependências

**Exemplo**:
```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/mock"
)
```

## Mocks

### O que são Mocks?

Mocks são objetos simulados que imitam o comportamento de dependências reais, permitindo testar componentes isoladamente.

### Estrutura de Mocks

Os mocks implementam as interfaces dos repositories:

```go
type MockAnimalRepository struct {
    mock.Mock
}

func (m *MockAnimalRepository) Create(animal *models.Animal) error {
    args := m.Called(animal)
    return args.Error(0)
}

func (m *MockAnimalRepository) FindByID(id uint) (*models.Animal, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Animal), args.Error(1)
}
```

### Localização dos Mocks

- `tests/services/mocks.go`: Mocks de repositories para testes de services
- `tests/handlers/mocks.go`: Mocks para testes de handlers
- `tests/mocks/`: Mocks adicionais

### Como Usar Mocks

1. **Criar instância do mock**:
```go
mockRepo := new(MockAnimalRepository)
```

2. **Configurar comportamento esperado**:
```go
mockRepo.On("FindByID", uint(1)).Return(&models.Animal{ID: 1}, nil)
mockRepo.On("Create", mock.AnythingOfType("*models.Animal")).Return(nil)
```

3. **Usar o mock no teste**:
```go
service := service.NewAnimalService(mockRepo)
err := service.CreateAnimal(animal)
```

4. **Verificar se todas as expectativas foram atendidas**:
```go
mockRepo.AssertExpectations(t)
```

## Exemplos de Testes

### Teste Unitário de Service

```go
func TestAnimalService_CreateAnimal(t *testing.T) {
    // 1. Criar mock
    mockRepo := new(MockAnimalRepository)
    service := service.NewAnimalService(mockRepo)

    // 2. Preparar dados
    animal := &models.Animal{
        FarmID:            1,
        EarTagNumberLocal: 123,
        AnimalName:        "Vaca Teste",
        Sex:               0,
        Breed:             "Holandesa",
        Type:              "Bovino",
        Purpose:           1,
    }

    // 3. Configurar expectativas do mock
    mockRepo.On("FindByEarTagNumber", uint(1), 123).Return((*models.Animal)(nil), nil)
    mockRepo.On("Create", animal).Return(nil)

    // 4. Executar função testada
    err := service.CreateAnimal(animal)

    // 5. Verificar resultados
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### Teste Unitário de Handler

```go
func TestAnimalHandler_CreateAnimal_Success(t *testing.T) {
    // 1. Setup
    mockRepo := new(services.MockAnimalRepository)
    router, _ := setupAnimalRouter(mockRepo)

    // 2. Preparar requisição
    animalData := map[string]interface{}{
        "farm_id":              1,
        "ear_tag_number_local": 123,
        "animal_name":          "Boi Teste",
        "sex":                  1,
        "breed":                "Nelore",
        "type":                 "Bovino",
    }
    jsonData, _ := json.Marshal(animalData)
    req, _ := http.NewRequest("POST", "/animals", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    // 3. Configurar mocks
    mockRepo.On("FindByEarTagNumber", uint(1), 123).Return((*models.Animal)(nil), nil)
    mockRepo.On("Create", mock.AnythingOfType("*models.Animal")).Return(nil).Run(func(args mock.Arguments) {
        animal := args.Get(0).(*models.Animal)
        animal.ID = 1
    })

    // 4. Executar requisição
    router.ServeHTTP(w, req)

    // 5. Verificar resultados
    assert.Equal(t, http.StatusCreated, w.Code)
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.True(t, response["success"].(bool))
    mockRepo.AssertExpectations(t)
}
```

### Teste de Integração

```go
func TestAuthIntegration(t *testing.T) {
    // 1. Setup com banco de dados real (SQLite em memória)
    app, db, cfg := setupTestApp(t)
    defer db.Close()
    router := routes.SetupRoutes(app, db, cfg)

    // 2. Testar registro de usuário
    t.Run("Register_User_Success", func(t *testing.T) {
        userData := map[string]interface{}{
            "user": map[string]interface{}{
                "farm_id": 1,
            },
            "person": map[string]interface{}{
                "first_name": "João",
                "last_name":  "Silva",
                "email":      "joao@test.com",
                "password":   "123456",
            },
        }
        jsonData, _ := json.Marshal(userData)
        req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
        req.Header.Set("Content-Type", "application/json")
        rr := httptest.NewRecorder()

        router.ServeHTTP(rr, req)

        assert.Equal(t, http.StatusCreated, rr.Code)
        var response map[string]interface{}
        json.Unmarshal(rr.Body.Bytes(), &response)
        assert.True(t, response["success"].(bool))
        assert.Contains(t, response, "access_token")
    })
}
```

## Setup de Testes

### Setup de Banco de Dados para Integração

```go
func setupTestDB(t *testing.T) *repository.Database {
    // 1. Criar banco SQLite em memória
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)

    // 2. Executar migrations
    err = migrations.RunMigrations(db)
    require.NoError(t, err)

    // 3. Criar dados iniciais (seed)
    company := &models.Company{
        CompanyName: "Test Company",
    }
    db.Create(company)

    farm := &models.Farm{
        CompanyID: company.ID,
    }
    db.Create(farm)

    return &repository.Database{DB: db}
}
```

### Setup de Aplicação Completa

```go
func setupTestApp(t *testing.T) (*app.Application, *repository.Database, *config.Config) {
    app, err := app.NewApplication()
    require.NoError(t, err)

    db := setupTestDB(t)

    cfg := &config.Config{
        Port:      "8080",
        JWTSecret: "test-secret-key",
        CORS: config.CORSConfig{
            AllowedOrigins:   []string{"*"},
            AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
            AllowedHeaders:   []string{"Content-Type", "Authorization"},
            AllowCredentials: true,
        },
    }

    return app, db, cfg
}
```

## Executando Testes

### Comandos Básicos

#### Executar Todos os Testes

```bash
go test ./...
```

#### Executar Testes de um Pacote Específico

```bash
go test ./tests/handlers
go test ./tests/services
go test ./tests/repositories
```

#### Executar Testes com Verbose

```bash
go test -v ./...
```

#### Executar Testes com Coverage

```bash
go test -cover ./...
```

#### Executar Testes com Coverage Detalhado

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Usando o Makefile

O projeto inclui um Makefile com comandos úteis:

```bash
# Executar todos os testes
make test

# Executar testes com coverage
make test-coverage

# Executar apenas testes unitários
make test-unit

# Executar apenas testes de handlers
make test-handlers
```

## Padrões de Nomenclatura

### Arquivos de Teste

- Nome do arquivo: `{nome}_test.go`
- Exemplo: `animal_test.go`, `auth_test.go`

### Funções de Teste

- Nome: `Test{NomeDaFunção}_{Cenário}`
- Exemplo: `TestAnimalService_CreateAnimal_Success`
- Exemplo: `TestAnimalHandler_CreateAnimal_InvalidJSON`

### Subtests

Use `t.Run()` para criar subtests:

```go
func TestAnimalService(t *testing.T) {
    t.Run("CreateAnimal_Success", func(t *testing.T) {
        // teste
    })
    
    t.Run("CreateAnimal_InvalidSex", func(t *testing.T) {
        // teste
    })
}
```

## Assertions Comuns

### assert vs require

- **assert**: Continua execução mesmo se falhar
- **require**: Para execução imediatamente se falhar

### Exemplos de Assertions

```go
// Igualdade
assert.Equal(t, expected, actual)
assert.NotEqual(t, expected, actual)

// Nil/Not Nil
assert.Nil(t, err)
assert.NotNil(t, result)

// Erros
assert.NoError(t, err)
assert.Error(t, err)
assert.Contains(t, err.Error(), "mensagem")

// Boolean
assert.True(t, condition)
assert.False(t, condition)

// Comparações
assert.Greater(t, value, threshold)
assert.Less(t, value, threshold)

// Strings
assert.Contains(t, str, substring)
assert.Empty(t, str)

// Slices/Maps
assert.Len(t, slice, expectedLength)
```

## Testes de Middleware

### Exemplo: Teste do Auth Middleware

```go
func TestAuthMiddleware(t *testing.T) {
    t.Run("Valid_Token", func(t *testing.T) {
        // Criar token válido
        token := generateTestToken()
        
        req, _ := http.NewRequest("GET", "/protected", nil)
        req.Header.Set("Authorization", "Bearer "+token)
        
        w := httptest.NewRecorder()
        handler := middleware.Auth("test-secret")
        handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
        })).ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusOK, w.Code)
    })
    
    t.Run("Invalid_Token", func(t *testing.T) {
        req, _ := http.NewRequest("GET", "/protected", nil)
        req.Header.Set("Authorization", "Bearer invalid-token")
        
        w := httptest.NewRecorder()
        handler := middleware.Auth("test-secret")
        handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
        })).ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusUnauthorized, w.Code)
    })
}
```

## Testes de Performance

O projeto inclui testes de performance:

```go
func TestPerformance(t *testing.T) {
    t.Run("Health_Check_Performance", func(t *testing.T) {
        start := time.Now()
        
        req, _ := http.NewRequest("GET", "/health", nil)
        rr := httptest.NewRecorder()
        router.ServeHTTP(rr, req)
        
        duration := time.Since(start)
        
        assert.Equal(t, http.StatusOK, rr.Code)
        assert.Less(t, duration, 100*time.Millisecond)
    })
}
```

## Cobertura de Código

### Verificar Cobertura

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### Visualizar Cobertura em HTML

```bash
go tool cover -html=coverage.out
```

### Meta de Cobertura

O projeto busca manter alta cobertura de código, especialmente em:
- Services (lógica de negócio)
- Handlers (endpoints HTTP)
- Repositories (acesso a dados)

## Boas Práticas

### 1. Isolamento

- Cada teste deve ser independente
- Não compartilhar estado entre testes
- Usar mocks para dependências externas

### 2. Nomenclatura Clara

- Nomes de testes devem descrever o que está sendo testado
- Incluir cenário esperado no nome (Success, Error, Invalid, etc.)

### 3. Arrange-Act-Assert (AAA)

```go
func TestExample(t *testing.T) {
    // Arrange: Preparar dados e mocks
    mockRepo := new(MockRepository)
    service := NewService(mockRepo)
    
    // Act: Executar ação testada
    result, err := service.DoSomething()
    
    // Assert: Verificar resultados
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 4. Testar Casos de Sucesso e Erro

Sempre teste:
- Caso de sucesso (happy path)
- Casos de erro
- Casos extremos (edge cases)
- Validações

### 5. Evitar Testes Frágeis

- Não depender de ordem de execução
- Não usar valores hardcoded quando possível
- Usar factories para criar dados de teste

## Estrutura de um Teste Completo

```go
func TestAnimalService_CreateAnimal(t *testing.T) {
    // Setup
    mockRepo := new(MockAnimalRepository)
    service := service.NewAnimalService(mockRepo)
    
    // Test cases
    t.Run("Success", func(t *testing.T) {
        // Arrange
        animal := createTestAnimal()
        mockRepo.On("FindByEarTagNumber", uint(1), 123).Return(nil, nil)
        mockRepo.On("Create", animal).Return(nil)
        
        // Act
        err := service.CreateAnimal(animal)
        
        // Assert
        assert.NoError(t, err)
        mockRepo.AssertExpectations(t)
    })
    
    t.Run("Duplicate_EarTag", func(t *testing.T) {
        // Arrange
        existingAnimal := &models.Animal{ID: 1}
        animal := createTestAnimal()
        mockRepo.On("FindByEarTagNumber", uint(1), 123).Return(existingAnimal, nil)
        
        // Act
        err := service.CreateAnimal(animal)
        
        // Assert
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "já existe")
    })
}
```

## Troubleshooting

### Problemas Comuns

1. **Mock não está sendo chamado**
   - Verificar se está usando `AssertExpectations(t)`
   - Verificar se os parâmetros do mock correspondem aos chamados

2. **Teste de integração falhando**
   - Verificar se migrations foram executadas
   - Verificar se dados de seed foram criados
   - Verificar se banco está sendo fechado (`defer db.Close()`)

3. **Testes lentos**
   - Usar mocks ao invés de banco real quando possível
   - Executar testes em paralelo quando seguro (`t.Parallel()`)

## Referências

- [Go Testing Package](https://pkg.go.dev/testing)
- [testify Documentation](https://github.com/stretchr/testify)
- Para entender a arquitetura, consulte `docs/arquitetura.md`
- Para entender os handlers, consulte `docs/handlers/`

