# Handler: User

## Visão Geral

O `UserHandler` gerencia operações relacionadas a usuários, incluindo criação, busca e atualização de dados pessoais.

## Estrutura

```go
type UserHandler struct {
    service *service.UserService
}
```

**Dependências**:
- `UserService`: Service que contém a lógica de negócio para usuários

## DTOs

### CreateUserRequest
```go
type CreateUserRequest struct {
    User   models.User   `json:"user"`
    Person models.Person `json:"person"`
}
```

## Métodos HTTP

### 1. GetUser
**Endpoint**: `GET /api/v1/users?email={email}`

**Descrição**: Busca um usuário por email.

**Parâmetros**: Query `email` (obrigatório)

**Resposta**: Retorna dados do usuário ou 404 se não encontrado.

---

### 2. CreateUser
**Endpoint**: `POST /api/v1/users`

**Descrição**: Cria um novo usuário no sistema.

**Parâmetros**: Body com `CreateUserRequest`

**Resposta**: Retorna ID do usuário criado (201 Created).

---

### 3. GetUserWithPerson
**Endpoint**: `GET /api/v1/users?id={userID}`

**Descrição**: Busca usuário com dados da pessoa associada.

**Parâmetros**: Query `id` (obrigatório)

**Nota**: Atualmente usa ID fixo (1) - pode ser um bug.

---

### 4. UpdatePersonData
**Endpoint**: `PUT /api/v1/users`

**Descrição**: Atualiza dados pessoais de um usuário.

**Parâmetros**: Body com `models.Person`

**Nota**: Atualmente usa userID fixo (1) - pode ser um bug.

