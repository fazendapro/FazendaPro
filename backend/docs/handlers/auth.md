# Handler: Auth (Autenticação)

## Visão Geral

O `AuthHandler` gerencia toda a autenticação do sistema, incluindo login, registro, renovação de tokens e logout. Utiliza JWT (JSON Web Tokens) para autenticação e refresh tokens para renovação de sessão.

## Estrutura

```go
type AuthHandler struct {
    service          *service.UserService
    refreshTokenRepo repository.RefreshTokenRepositoryInterface
    jwtSecret        string
}
```

**Dependências**:
- `UserService`: Service para operações de usuário
- `RefreshTokenRepository`: Repository para gerenciar refresh tokens
- `jwtSecret`: Chave secreta para assinar tokens JWT

**Construtor**:
```go
func NewAuthHandler(service *service.UserService, refreshTokenRepo repository.RefreshTokenRepositoryInterface, jwtSecret string) *AuthHandler
```

## DTOs (Data Transfer Objects)

### LoginRequest

Request para login:

```go
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

### RegisterRequest

Request para registro de novo usuário:

```go
type RegisterRequest struct {
    User   models.User   `json:"user"`
    Person models.Person `json:"person"`
}
```

### LoginResponse / RegisterResponse

Response após login ou registro bem-sucedido:

```go
type LoginResponse struct {
    Success      bool   `json:"success"`
    Message      string `json:"message"`
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    User         struct {
        ID    uint   `json:"id"`
        Email string `json:"email"`
        Name  string `json:"name"`
    } `json:"user"`
}
```

### RefreshTokenRequest

Request para renovar token:

```go
type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token"`
}
```

### RefreshTokenResponse

Response após renovação de token:

```go
type RefreshTokenResponse struct {
    Success     bool   `json:"success"`
    Message     string `json:"message"`
    AccessToken string `json:"access_token"`
}
```

## Métodos HTTP

### 1. Login

**Endpoint**: `POST /api/v1/auth/login`

**Método HTTP**: POST

**Autenticação**: Não requerida (endpoint público)

**Descrição**: Autentica um usuário e retorna tokens de acesso.

**Parâmetros**:
- Body (JSON): `LoginRequest`

**Validações**:
- Método HTTP deve ser POST
- JSON deve ser válido
- Email e senha devem ser fornecidos

**Fluxo**:
1. Valida método HTTP
2. Decodifica JSON do body
3. Busca usuário por email via service
4. Verifica se usuário existe
5. Valida senha via service
6. Gera JWT (access token)
7. Cria refresh token (válido por 7 dias)
8. Retorna tokens e dados do usuário

**Claims do JWT**:
```go
claims := jwt.MapClaims{
    "sub":     user.ID,           // Subject (ID do usuário)
    "email":   user.Person.Email,  // Email do usuário
    "farm_id": user.FarmID,        // ID da fazenda
    "iat":     time.Now().Unix(),  // Issued At
    "exp":     time.Now().Add(time.Hour * 24).Unix(), // Expiration (24 horas)
}
```

**Resposta de Sucesso** (200 OK):
```json
{
  "success": true,
  "message": "Login realizado com sucesso",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "550e8400-e29b-41d4-a716-446655440000",
  "user": {
    "id": 1,
    "email": "usuario@example.com",
    "name": "João Silva"
  }
}
```

**Resposta de Erro**:
- `400 Bad Request`: JSON inválido
- `401 Unauthorized`: Credenciais inválidas
- `405 Method Not Allowed`: Método HTTP incorreto
- `500 Internal Server Error`: Erro ao gerar tokens

**Exemplo de Requisição**:
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "usuario@example.com",
  "password": "senha123"
}
```

---

### 2. Register

**Endpoint**: `POST /api/v1/auth/register`

**Método HTTP**: POST

**Autenticação**: Não requerida (endpoint público)

**Descrição**: Registra um novo usuário no sistema e automaticamente faz login.

**Parâmetros**:
- Body (JSON): `RegisterRequest`

**Validações**:
- Método HTTP deve ser POST
- JSON deve ser válido
- Validações de negócio são feitas no service

**Fluxo**:
1. Valida método HTTP
2. Decodifica JSON do body
3. Cria usuário e pessoa via service
4. Busca usuário criado por email
5. Gera JWT (access token)
6. Cria refresh token (válido por 7 dias)
7. Retorna tokens e dados do usuário

**Resposta de Sucesso** (201 Created):
```json
{
  "success": true,
  "message": "Usuário criado e logado com sucesso",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "550e8400-e29b-41d4-a716-446655440000",
  "user": {
    "id": 1,
    "email": "novo@example.com",
    "name": "Maria Santos"
  }
}
```

**Resposta de Erro**:
- `400 Bad Request`: JSON inválido ou erro de validação (ex: email já existe)
- `405 Method Not Allowed`: Método HTTP incorreto
- `500 Internal Server Error`: Erro ao criar usuário ou gerar tokens

**Exemplo de Requisição**:
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "user": {
    "farm_id": 1
  },
  "person": {
    "first_name": "Maria",
    "last_name": "Santos",
    "email": "maria@example.com",
    "phone": "11999999999"
  }
}
```

**Nota**: A senha deve ser fornecida no campo apropriado do modelo User.

---

### 3. RefreshToken

**Endpoint**: `POST /api/v1/auth/refresh`

**Método HTTP**: POST

**Autenticação**: Não requerida (usa refresh token)

**Descrição**: Renova o access token usando um refresh token válido.

**Parâmetros**:
- Body (JSON): `RefreshTokenRequest`

**Validações**:
- Método HTTP deve ser POST
- JSON deve ser válido
- Refresh token deve ser fornecido

**Fluxo**:
1. Valida método HTTP
2. Decodifica JSON do body
3. Busca refresh token no banco
4. Verifica se refresh token existe e é válido
5. Gera novo JWT (access token) para o usuário
6. Retorna novo access token

**Resposta de Sucesso** (200 OK):
```json
{
  "success": true,
  "message": "Token renovado com sucesso",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Resposta de Erro**:
- `400 Bad Request`: JSON inválido
- `401 Unauthorized`: Refresh token inválido ou expirado
- `405 Method Not Allowed`: Método HTTP incorreto
- `500 Internal Server Error`: Erro ao gerar token

**Exemplo de Requisição**:
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Nota**: O refresh token não é renovado, apenas o access token. O refresh token continua válido até expirar (7 dias).

---

### 4. Logout

**Endpoint**: `POST /api/v1/auth/logout`

**Método HTTP**: POST

**Autenticação**: Não requerida (usa refresh token)

**Descrição**: Invalida um refresh token, efetivamente fazendo logout do usuário.

**Parâmetros**:
- Body (JSON): `RefreshTokenRequest`

**Validações**:
- Método HTTP deve ser POST
- JSON deve ser válido
- Refresh token deve ser fornecido

**Fluxo**:
1. Valida método HTTP
2. Decodifica JSON do body
3. Remove refresh token do banco
4. Retorna confirmação

**Resposta de Sucesso** (200 OK):
```json
{
  "success": true,
  "message": "Logout realizado com sucesso"
}
```

**Resposta de Erro**:
- `400 Bad Request`: JSON inválido
- `405 Method Not Allowed`: Método HTTP incorreto
- `500 Internal Server Error`: Erro ao remover token

**Exemplo de Requisição**:
```http
POST /api/v1/auth/logout
Content-Type: application/json

{
  "refresh_token": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Nota**: O access token JWT continua válido até expirar (24 horas), mas sem o refresh token, o usuário não poderá renovar a sessão.

---

## Funções Auxiliares

### generateJWT

Gera um token JWT para um usuário:

```go
func (h *AuthHandler) generateJWT(user *models.User) (string, error)
```

**Funcionalidades**:
- Cria claims JWT com informações do usuário
- Assina token com HS256
- Token expira em 24 horas
- Inclui: user ID, email, farm_id, iat, exp

**Claims Incluídos**:
- `sub`: ID do usuário
- `email`: Email do usuário
- `farm_id`: ID da fazenda do usuário
- `iat`: Timestamp de criação
- `exp`: Timestamp de expiração (24 horas)

---

## Fluxo de Autenticação Completo

### 1. Registro/Login Inicial

```
Cliente → POST /auth/register ou /auth/login
         ↓
    Handler valida credenciais
         ↓
    Gera Access Token (JWT, 24h)
         ↓
    Gera Refresh Token (UUID, 7 dias)
         ↓
    Retorna ambos os tokens
```

### 2. Uso do Access Token

```
Cliente → Requisição com Header: Authorization: Bearer {access_token}
         ↓
    Middleware Auth valida token
         ↓
    Se válido: permite acesso
    Se inválido/expirado: retorna 401
```

### 3. Renovação de Token

```
Cliente → POST /auth/refresh com refresh_token
         ↓
    Handler valida refresh_token
         ↓
    Gera novo Access Token
         ↓
    Retorna novo Access Token
```

### 4. Logout

```
Cliente → POST /auth/logout com refresh_token
         ↓
    Handler remove refresh_token do banco
         ↓
    Cliente não pode mais renovar sessão
```

---

## Segurança

### Access Token (JWT)
- **Validade**: 24 horas
- **Algoritmo**: HS256
- **Conteúdo**: ID do usuário, email, farm_id
- **Armazenamento**: Cliente (localStorage, sessionStorage, etc.)
- **Uso**: Enviado em todas as requisições autenticadas

### Refresh Token
- **Validade**: 7 dias
- **Tipo**: UUID armazenado no banco
- **Armazenamento**: Banco de dados + Cliente
- **Uso**: Apenas para renovar access token
- **Revogação**: Pode ser removido no logout

### Boas Práticas Implementadas
1. **Senhas**: Nunca retornadas nas respostas
2. **Validação**: Credenciais validadas antes de gerar tokens
3. **Expiração**: Tokens têm tempo de vida limitado
4. **Revogação**: Refresh tokens podem ser invalidados
5. **Mensagens Genéricas**: Erros de autenticação não revelam se email existe

---

## Exemplos de Uso

### Fluxo Completo de Autenticação

```javascript
// 1. Login
const loginResponse = await fetch('/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'usuario@example.com',
    password: 'senha123'
  })
});

const { access_token, refresh_token } = await loginResponse.json();

// 2. Usar Access Token
const animalsResponse = await fetch('/api/v1/animals/farm?farmId=1', {
  headers: {
    'Authorization': `Bearer ${access_token}`
  }
});

// 3. Renovar Token quando expirar
const refreshResponse = await fetch('/api/v1/auth/refresh', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    refresh_token: refresh_token
  })
});

const { access_token: newAccessToken } = await refreshResponse.json();

// 4. Logout
await fetch('/api/v1/auth/logout', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    refresh_token: refresh_token
  })
});
```

---

## Tratamento de Erros

### Erros Comuns

1. **Credenciais Inválidas** (401)
   - Email não existe
   - Senha incorreta
   - Resposta genérica: "Credenciais inválidas"

2. **Token Expirado** (401)
   - Access token expirado → usar refresh token
   - Refresh token expirado → fazer login novamente

3. **Token Inválido** (401)
   - Token malformado
   - Token assinado com chave diferente

4. **Erro ao Gerar Token** (500)
   - Problema interno do servidor
   - Verificar logs do servidor

---

## Notas de Implementação

1. **Nome Completo**: O nome completo do usuário é construído concatenando `FirstName + " " + LastName` da Person
2. **Refresh Token**: Criado com expiração de 7 dias usando `time.Now().Add(time.Hour*24*7)`
3. **JWT Secret**: Deve ser uma string segura e aleatória, armazenada em variável de ambiente
4. **Validação de Senha**: Feita no service usando bcrypt
5. **Respostas HTTP**: Login e Register retornam JSON diretamente (não usam SendSuccessResponse)

