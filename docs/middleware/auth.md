# Middleware: Auth (Autenticação)

## Visão Geral

O middleware `Auth` é responsável por validar tokens JWT (JSON Web Tokens) em requisições HTTP. Ele verifica se o token é válido e extrai informações do usuário (como `farm_id`) para disponibilizar no contexto da requisição.

## Localização

`internal/api/middleware/auth.go`

## Assinatura

```go
func Auth(jwtSecret string) func(http.Handler) http.Handler
```

**Parâmetros**:
- `jwtSecret`: Chave secreta usada para validar e assinar tokens JWT

**Retorno**: Função middleware compatível com Chi Router

## Como Funciona

### Fluxo de Execução

```
Requisição HTTP
    ↓
Middleware Auth intercepta
    ↓
Extrai token do header Authorization
    ↓
Valida token JWT
    ↓
Se válido: extrai farm_id e adiciona ao contexto
    ↓
Se inválido: retorna 401 Unauthorized
    ↓
Passa requisição para próximo handler
```

### Passo a Passo

1. **Extração do Token**
   ```go
   tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
   ```
   - Remove o prefixo "Bearer " do header Authorization
   - Se não houver token, retorna erro 401

2. **Validação do Token**
   ```go
   token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
       return []byte(jwtSecret), nil
   })
   ```
   - Usa a biblioteca `golang-jwt/jwt/v5` para fazer parse do token
   - Valida assinatura usando `jwtSecret`
   - Verifica se o token é válido e não expirou

3. **Extração de Claims**
   ```go
   if claims, ok := token.Claims.(jwt.MapClaims); ok {
       if farmID, exists := claims["farm_id"]; exists {
           if farmIDFloat, ok := farmID.(float64); ok {
               ctx := r.Context()
               ctx = context.WithValue(ctx, "farm_id", uint(farmIDFloat))
               r = r.WithContext(ctx)
           }
       }
   }
   ```
   - Extrai claims do token
   - Busca o `farm_id` nos claims
   - Adiciona `farm_id` ao contexto da requisição
   - Converte de `float64` (formato JSON) para `uint`

4. **Continuação da Requisição**
   ```go
   next.ServeHTTP(w, r)
   ```
   - Se tudo estiver válido, passa a requisição para o próximo handler
   - O contexto agora contém o `farm_id` disponível

## Uso nas Rotas

### Aplicação em Grupo de Rotas

```go
r.Route("/animals", func(r chi.Router) {
    r.Use(middleware.Auth(cfg.JWTSecret))
    r.Post("/", animalHandler.CreateAnimal)
    r.Get("/", animalHandler.GetAnimal)
})
```

### Aplicação em Rota Individual

```go
r.Route("/api/v1", func(r chi.Router) {
    r.Use(middleware.Auth(cfg.JWTSecret))
    // Todas as rotas dentro deste grupo requerem autenticação
})
```

## Formato do Token JWT

O middleware espera tokens no formato:

```
Authorization: Bearer {token}
```

Onde `{token}` é um JWT válido assinado com `jwtSecret`.

### Claims Esperados

O token deve conter os seguintes claims:

- `sub`: ID do usuário (subject)
- `email`: Email do usuário
- `farm_id`: ID da fazenda (usado pelo middleware)
- `iat`: Timestamp de criação (issued at)
- `exp`: Timestamp de expiração (expiration)

## Acesso ao Farm ID no Handler

Após passar pelo middleware, os handlers podem acessar o `farm_id` do contexto:

```go
func (h *SaleHandler) CreateSale(w http.ResponseWriter, r *http.Request) {
    farmID, ok := r.Context().Value("farm_id").(uint)
    if !ok {
        http.Error(w, "Farm ID not found in context", http.StatusUnauthorized)
        return
    }
    // Usar farmID...
}
```

## Respostas de Erro

### Token Não Fornecido

**Status**: `401 Unauthorized`

**Resposta**:
```json
{
  "success": false,
  "error": "Unauthorized",
  "message": "Token não fornecido",
  "code": 401
}
```

**Quando ocorre**: Header `Authorization` não presente ou não contém "Bearer ".

---

### Token Inválido

**Status**: `401 Unauthorized`

**Resposta**:
```json
{
  "success": false,
  "error": "Unauthorized",
  "message": "Token inválido",
  "code": 401
}
```

**Quando ocorre**:
- Token malformado
- Token assinado com chave diferente
- Token expirado
- Token não pode ser parseado

## Exemplo de Requisição Autenticada

```http
GET /api/v1/animals/farm?farmId=1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjEsImVtYWlsIjoi...
```

## Segurança

### Boas Práticas Implementadas

1. **Validação de Assinatura**: Token deve ser assinado com a chave secreta correta
2. **Validação de Expiração**: Tokens expirados são rejeitados automaticamente
3. **Extração Segura**: Claims são extraídos de forma type-safe
4. **Contexto Isolado**: Cada requisição tem seu próprio contexto

### Considerações

1. **JWT Secret**: Deve ser uma string segura e aleatória, armazenada em variável de ambiente
2. **HTTPS**: Em produção, sempre use HTTPS para proteger tokens em trânsito
3. **Expiração**: Tokens têm tempo de vida limitado (24 horas no projeto)
4. **Refresh Tokens**: Use refresh tokens para renovar access tokens sem reautenticação

## Função Auxiliar: SendErrorResponse

O middleware também define uma função auxiliar para respostas de erro:

```go
func SendErrorResponse(w http.ResponseWriter, message string, statusCode int)
```

Esta função padroniza o formato de erro retornado pelo middleware.

## Integração com Outros Middlewares

O middleware Auth é geralmente aplicado após o CORS middleware:

```go
r.Use(middleware.CORSMiddleware(cfg))
r.Use(middleware.Auth(cfg.JWTSecret))
```

**Ordem Importante**:
1. CORS primeiro (permite requisições cross-origin)
2. Auth depois (valida autenticação)

## Rotas Públicas

Rotas que não devem usar este middleware:

- `/api/v1/auth/login`
- `/api/v1/auth/register`
- `/api/v1/auth/refresh`
- `/api/v1/auth/logout`
- `/health`
- `/`

Estas rotas são definidas antes da aplicação do middleware Auth.

## Debugging

### Problemas Comuns

1. **"Token não fornecido"**
   - Verificar se header `Authorization` está presente
   - Verificar se está no formato `Bearer {token}`

2. **"Token inválido"**
   - Verificar se token não expirou
   - Verificar se `jwtSecret` está correto
   - Verificar se token foi assinado corretamente

3. **"Farm ID not found in context"**
   - Verificar se token contém claim `farm_id`
   - Verificar se middleware Auth foi aplicado na rota

### Logs

O middleware não gera logs próprios. Para debug, verifique:
- Logs do servidor HTTP
- Respostas de erro do middleware
- Validação de tokens no handler de autenticação

## Código Completo

```go
func Auth(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 1. Extrai token
            tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
            if tokenString == "" {
                SendErrorResponse(w, "Token não fornecido", http.StatusUnauthorized)
                return
            }

            // 2. Valida token
            token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                return []byte(jwtSecret), nil
            })

            if err != nil || !token.Valid {
                SendErrorResponse(w, "Token inválido", http.StatusUnauthorized)
                return
            }

            // 3. Extrai farm_id e adiciona ao contexto
            if claims, ok := token.Claims.(jwt.MapClaims); ok {
                if farmID, exists := claims["farm_id"]; exists {
                    if farmIDFloat, ok := farmID.(float64); ok {
                        ctx := r.Context()
                        ctx = context.WithValue(ctx, "farm_id", uint(farmIDFloat))
                        r = r.WithContext(ctx)
                    }
                }
            }

            // 4. Continua requisição
            next.ServeHTTP(w, r)
        })
    }
}
```

## Referências

- Para entender como tokens são gerados, consulte `docs/handlers/auth.md`
- Para entender o sistema de rotas, consulte `docs/routes.md`
- Para entender a arquitetura, consulte `docs/arquitetura.md`

