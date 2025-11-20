# Documentação dos Middlewares

Esta pasta contém a documentação detalhada de todos os middlewares da API FazendaPro.

## Índice

### Middlewares de Segurança e Autenticação

1. **[Auth Middleware](auth.md)** - Autenticação JWT
   - Validação de tokens JWT
   - Extração de farm_id do contexto
   - Proteção de rotas

2. **[CORS Middleware](cors.md)** - Cross-Origin Resource Sharing
   - Configuração de políticas CORS
   - Tratamento de requisições preflight
   - Gerenciamento de origens permitidas

## O que são Middlewares?

Middlewares são funções que interceptam requisições HTTP antes que cheguem aos handlers finais. Eles podem:

- Validar autenticação
- Modificar requisições/respostas
- Adicionar headers
- Registrar logs
- Tratar erros

## Ordem de Execução

Os middlewares são executados na ordem em que são aplicados:

```
Requisição HTTP
    ↓
1. Sentry Middleware (se configurado)
    ↓
2. CORS Middleware
    ↓
3. Logging Middleware
    ↓
4. Auth Middleware (em rotas protegidas)
    ↓
Handler Final
```

## Aplicação de Middlewares

### Global

Middlewares aplicados a todas as rotas:

```go
r := chi.NewRouter()
r.Use(middleware.CORSMiddleware(cfg))
r.Use(loggingMiddleware)
```

### Por Grupo

Middlewares aplicados a grupos específicos de rotas:

```go
r.Route("/api/v1/animals", func(r chi.Router) {
    r.Use(middleware.Auth(cfg.JWTSecret))
    r.Post("/", handler.CreateAnimal)
})
```

### Por Rota

Middlewares aplicados a rotas individuais (menos comum):

```go
r.With(middleware.Auth(cfg.JWTSecret)).Get("/protected", handler.Protected)
```

## Middlewares do Projeto

### 1. Sentry Middleware

**Arquivo**: Configurado em `routes.go`

**Função**: Captura erros e panics para monitoramento.

**Aplicação**: Global (se SentryDSN configurado)

**Documentação**: Ver `docs/routes.md` para detalhes.

---

### 2. CORS Middleware

**Arquivo**: `internal/api/middleware/cors.go`

**Função**: Gerencia políticas CORS.

**Aplicação**: Global

**Documentação**: [cors.md](cors.md)

---

### 3. Logging Middleware

**Arquivo**: Configurado em `routes.go`

**Função**: Registra todas as requisições HTTP.

**Aplicação**: Global

**Documentação**: Ver `docs/routes.md` para detalhes.

---

### 4. Auth Middleware

**Arquivo**: `internal/api/middleware/auth.go`

**Função**: Valida tokens JWT e protege rotas.

**Aplicação**: Por grupo (rotas protegidas)

**Documentação**: [auth.md](auth.md)

## Padrão de Middleware no Chi

Todos os middlewares seguem o padrão do Chi Router:

```go
func MyMiddleware(param string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Lógica do middleware
            // ...
            
            // Passa para próximo handler
            next.ServeHTTP(w, r)
        })
    }
}
```

## Contexto HTTP

Middlewares podem adicionar valores ao contexto da requisição:

```go
ctx := r.Context()
ctx = context.WithValue(ctx, "key", value)
r = r.WithContext(ctx)
```

Handlers podem acessar esses valores:

```go
value := r.Context().Value("key")
```

## Tratamento de Erros

Middlewares podem retornar erros antes de chegar aos handlers:

```go
if error {
    SendErrorResponse(w, "Error message", http.StatusUnauthorized)
    return // Não chama next.ServeHTTP
}
```

## Boas Práticas

1. **Ordem Importante**: CORS antes de Auth
2. **Contexto**: Use contexto para passar dados entre middlewares e handlers
3. **Erros**: Retorne erros apropriados e não continue a requisição
4. **Performance**: Middlewares devem ser rápidos (não fazer I/O pesado)
5. **Reutilização**: Middlewares devem ser genéricos e reutilizáveis

## Exemplo Completo

```go
// 1. Cria router
r := chi.NewRouter()

// 2. Middlewares globais
r.Use(middleware.CORSMiddleware(cfg))
r.Use(loggingMiddleware)

// 3. Rotas públicas
r.Post("/api/v1/auth/login", authHandler.Login)

// 4. Rotas protegidas
r.Route("/api/v1/animals", func(r chi.Router) {
    r.Use(middleware.Auth(cfg.JWTSecret))
    r.Post("/", animalHandler.CreateAnimal)
    r.Get("/", animalHandler.GetAnimal)
})
```

## Referências

- Para entender como middlewares são aplicados nas rotas, consulte `docs/routes.md`
- Para entender a arquitetura geral, consulte `docs/arquitetura.md`
- [Chi Router Middleware](https://github.com/go-chi/chi#middlewares)

