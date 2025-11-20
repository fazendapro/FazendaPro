# Middleware: CORS (Cross-Origin Resource Sharing)

## Visão Geral

O middleware `CORSMiddleware` gerencia as políticas de Cross-Origin Resource Sharing (CORS) da API. Ele permite que aplicações web em diferentes origens façam requisições à API, configurando os headers HTTP apropriados.

## Localização

`internal/api/middleware/cors.go`

## Assinatura

```go
func CORSMiddleware(cfg *config.Config) func(http.Handler) http.Handler
```

**Parâmetros**:
- `cfg`: Configuração da aplicação contendo configurações CORS

**Retorno**: Função middleware compatível com Chi Router

## Como Funciona

### Fluxo de Execução

```
Requisição HTTP
    ↓
Middleware CORS intercepta
    ↓
Verifica origem da requisição
    ↓
Configura headers CORS apropriados
    ↓
Se OPTIONS: retorna resposta preflight
    ↓
Passa requisição para próximo handler
```

## Configuração CORS

O middleware usa a estrutura `CORSConfig` do arquivo de configuração:

```go
type CORSConfig struct {
    AllowedOrigins   []string
    AllowedMethods   []string
    AllowedHeaders   []string
    ExposedHeaders   []string
    AllowCredentials bool
    MaxAge           int
}
```

### Variáveis de Ambiente

As configurações são carregadas de variáveis de ambiente:

- `CORS_ALLOWED_ORIGINS`: Origens permitidas (separadas por vírgula)
- `CORS_ALLOWED_METHODS`: Métodos HTTP permitidos (padrão: GET,POST,PUT,DELETE,OPTIONS)
- `CORS_ALLOWED_HEADERS`: Headers permitidos (padrão: Content-Type,Authorization)
- `CORS_EXPOSED_HEADERS`: Headers expostos ao cliente
- `CORS_ALLOW_CREDENTIALS`: Permitir credenciais (padrão: true)
- `CORS_MAX_AGE`: Tempo de cache do preflight (padrão: 86400 segundos)

## Headers CORS Configurados

### 1. Access-Control-Allow-Origin

**Função**: Define quais origens podem acessar a API.

**Lógica**:
- Se `*` estiver nas origens permitidas → permite todas as origens
- Se origem específica corresponder → permite apenas essa origem
- Se nenhuma corresponder → não define o header (bloqueia)

**Exemplo**:
```go
w.Header().Set("Access-Control-Allow-Origin", "https://app.fazendapro.com")
```

### 2. Access-Control-Allow-Methods

**Função**: Define quais métodos HTTP são permitidos.

**Aplicado**: Apenas em requisições OPTIONS (preflight).

**Exemplo**:
```go
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
```

### 3. Access-Control-Allow-Headers

**Função**: Define quais headers podem ser enviados pelo cliente.

**Aplicado**: Apenas em requisições OPTIONS (preflight).

**Exemplo**:
```go
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
```

### 4. Access-Control-Expose-Headers

**Função**: Define quais headers podem ser lidos pelo cliente na resposta.

**Aplicado**: Em todas as requisições.

**Exemplo**:
```go
w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
```

### 5. Access-Control-Allow-Credentials

**Função**: Permite envio de credenciais (cookies, auth headers) em requisições cross-origin.

**Aplicado**: Apenas se `AllowCredentials` for `true` e origem estiver presente.

**Exemplo**:
```go
w.Header().Set("Access-Control-Allow-Credentials", "true")
```

**Nota**: Não pode ser usado com `Access-Control-Allow-Origin: *`. Deve especificar origem exata.

### 6. Access-Control-Max-Age

**Função**: Define por quanto tempo o navegador pode cachear a resposta do preflight.

**Aplicado**: Apenas em requisições OPTIONS (preflight).

**Exemplo**:
```go
w.Header().Set("Access-Control-Max-Age", "86400") // 24 horas
```

## Requisições Preflight (OPTIONS)

### O que é Preflight?

Antes de fazer uma requisição cross-origin "complexa", o navegador envia uma requisição OPTIONS para verificar se a requisição é permitida.

### Quando Ocorre Preflight?

- Métodos HTTP diferentes de GET, HEAD ou POST
- POST com Content-Type diferente de `application/x-www-form-urlencoded`, `multipart/form-data` ou `text/plain`
- Headers customizados (como `Authorization`)

### Tratamento no Middleware

```go
if r.Method == "OPTIONS" {
    w.WriteHeader(http.StatusOK)
    return
}
```

O middleware responde imediatamente a requisições OPTIONS com status 200 OK e os headers CORS apropriados, sem passar para os handlers seguintes.

## Lógica de Origens Permitidas

### 1. Permitir Todas as Origens

Se `*` estiver na lista de origens permitidas:

```go
if allowAll {
    w.Header().Set("Access-Control-Allow-Origin", "*")
}
```

**Uso**: Desenvolvimento local

**Limitação**: Não pode usar `Access-Control-Allow-Credentials: true`

### 2. Origens Específicas

Se a origem da requisição corresponder a uma origem permitida:

```go
if allowedOrigin == origin {
    w.Header().Set("Access-Control-Allow-Origin", origin)
    originAllowed = true
}
```

**Uso**: Produção

**Vantagem**: Pode usar credenciais

### 3. Origem Não Permitida

Se a origem não corresponder a nenhuma permitida:

```go
if !originAllowed {
    log.Printf("CORS: Origin %s not allowed", origin)
}
```

O header não é definido, bloqueando a requisição.

## Uso nas Rotas

### Aplicação Global

O middleware CORS é aplicado globalmente em todas as rotas:

```go
r := chi.NewRouter()
r.Use(middleware.CORSMiddleware(cfg))
// Todas as rotas agora têm CORS configurado
```

### Ordem de Middlewares

O CORS deve ser aplicado antes de outros middlewares:

```go
r.Use(middleware.CORSMiddleware(cfg))  // 1. CORS primeiro
r.Use(middleware.Auth(cfg.JWTSecret))  // 2. Auth depois
```

## Exemplos de Configuração

### Desenvolvimento Local

```env
CORS_ALLOWED_ORIGINS=*
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization
CORS_ALLOW_CREDENTIALS=false
CORS_MAX_AGE=86400
```

### Produção

```env
CORS_ALLOWED_ORIGINS=https://app.fazendapro.com,https://www.fazendapro.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization
CORS_EXPOSED_HEADERS=X-Total-Count
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=86400
```

## Exemplo de Requisição Cross-Origin

### Requisição Preflight (OPTIONS)

```http
OPTIONS /api/v1/animals
Origin: https://app.fazendapro.com
Access-Control-Request-Method: POST
Access-Control-Request-Headers: Content-Type,Authorization
```

**Resposta**:
```http
HTTP/1.1 200 OK
Access-Control-Allow-Origin: https://app.fazendapro.com
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
Access-Control-Allow-Credentials: true
Access-Control-Max-Age: 86400
```

### Requisição Real (POST)

```http
POST /api/v1/animals
Origin: https://app.fazendapro.com
Authorization: Bearer {token}
Content-Type: application/json

{
  "animal_name": "Branquinha",
  ...
}
```

**Resposta**:
```http
HTTP/1.1 201 Created
Access-Control-Allow-Origin: https://app.fazendapro.com
Access-Control-Allow-Credentials: true
Content-Type: application/json

{
  "success": true,
  ...
}
```

## Segurança

### Boas Práticas

1. **Origens Específicas em Produção**: Nunca use `*` em produção
2. **Credenciais**: Use `AllowCredentials: true` apenas com origens específicas
3. **Headers Mínimos**: Exponha apenas headers necessários
4. **Max-Age Apropriado**: Balance entre segurança e performance

### Considerações

1. **Wildcard vs Específico**: `*` é mais permissivo, mas não permite credenciais
2. **Múltiplas Origens**: Suporta múltiplas origens separadas por vírgula
3. **Logs**: Registra tentativas de acesso de origens não permitidas

## Debugging

### Problemas Comuns

1. **"CORS policy blocked"**
   - Verificar se origem está em `CORS_ALLOWED_ORIGINS`
   - Verificar se não está usando `*` com credenciais

2. **Preflight falhando**
   - Verificar se métodos estão em `CORS_ALLOWED_METHODS`
   - Verificar se headers estão em `CORS_ALLOWED_HEADERS`

3. **Credenciais não funcionando**
   - Verificar se não está usando `*` como origem
   - Verificar se `CORS_ALLOW_CREDENTIALS=true`

### Logs

O middleware registra tentativas de acesso de origens não permitidas:

```
CORS: Origin https://malicious-site.com not allowed
```

## Código Completo

```go
func CORSMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")

            // Configura Access-Control-Allow-Origin
            if len(cfg.CORS.AllowedOrigins) > 0 {
                allowAll := false
                for _, allowedOrigin := range cfg.CORS.AllowedOrigins {
                    if allowedOrigin == "*" {
                        allowAll = true
                        break
                    }
                }

                if allowAll {
                    w.Header().Set("Access-Control-Allow-Origin", "*")
                } else if origin != "" {
                    originAllowed := false
                    for _, allowedOrigin := range cfg.CORS.AllowedOrigins {
                        if allowedOrigin == origin {
                            w.Header().Set("Access-Control-Allow-Origin", origin)
                            originAllowed = true
                            break
                        }
                    }
                    if !originAllowed {
                        log.Printf("CORS: Origin %s not allowed", origin)
                    }
                } else {
                    if len(cfg.CORS.AllowedOrigins) > 0 {
                        w.Header().Set("Access-Control-Allow-Origin", cfg.CORS.AllowedOrigins[0])
                    }
                }
            }

            // Configura outros headers CORS
            if r.Method == "OPTIONS" && len(cfg.CORS.AllowedMethods) > 0 {
                methods := strings.Join(cfg.CORS.AllowedMethods, ", ")
                w.Header().Set("Access-Control-Allow-Methods", methods)
            }

            if r.Method == "OPTIONS" && len(cfg.CORS.AllowedHeaders) > 0 {
                headers := strings.Join(cfg.CORS.AllowedHeaders, ", ")
                w.Header().Set("Access-Control-Allow-Headers", headers)
            }

            if len(cfg.CORS.ExposedHeaders) > 0 {
                exposedHeaders := strings.Join(cfg.CORS.ExposedHeaders, ", ")
                w.Header().Set("Access-Control-Expose-Headers", exposedHeaders)
            }

            if cfg.CORS.AllowCredentials && origin != "" {
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }

            if r.Method == "OPTIONS" && cfg.CORS.MaxAge > 0 {
                maxAge := fmt.Sprintf("%d", cfg.CORS.MaxAge)
                w.Header().Set("Access-Control-Max-Age", maxAge)
            }

            // Responde preflight imediatamente
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }

            // Continua requisição
            next.ServeHTTP(w, r)
        })
    }
}
```

## Referências

- Para entender como as rotas são configuradas, consulte `docs/routes.md`
- Para entender a configuração, consulte `docs/arquitetura.md`
- [MDN: CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)

