# Integração New Relic

## Visão Geral

O FazendaPro Backend está integrado com o New Relic APM (Application Performance Monitoring) para monitoramento de performance, transações HTTP, erros e logs da aplicação.

## Configuração

### Variáveis de Ambiente

Adicione as seguintes variáveis ao seu arquivo `.env`:

```bash
# New Relic Configuration
NEW_RELIC_ENABLED=true
NEW_RELIC_LICENSE_KEY=sua_chave_de_licenca_aqui
NEW_RELIC_APP_NAME=fazendapro-backend
```

### Parâmetros

- **NEW_RELIC_ENABLED**: Define se o New Relic está habilitado (true/false)
  - Padrão: `false`
  - Recomendado: `true` em produção, `false` em desenvolvimento

- **NEW_RELIC_LICENSE_KEY**: Chave de licença do New Relic
  - Obrigatório quando `NEW_RELIC_ENABLED=true`
  - Obtida no dashboard do New Relic

- **NEW_RELIC_APP_NAME**: Nome da aplicação no New Relic
  - Padrão: `fazendapro-backend`
  - Pode ser customizado por ambiente (ex: `fazendapro-backend-prod`)

## Funcionalidades

### 1. Monitoramento Automático de Rotas HTTP

Todas as rotas HTTP são automaticamente instrumentadas, exceto:
- `/health` - excluída para reduzir ruído no monitoramento

### 2. Informações Capturadas

Para cada requisição HTTP, o New Relic captura:
- Método HTTP (GET, POST, PUT, DELETE, etc.)
- Path da requisição
- Status code da resposta
- Tempo de resposta
- Headers da requisição
- Erros e exceções

### 3. Distributed Tracing

O distributed tracing está habilitado, permitindo rastrear requisições através de múltiplos serviços.

### 4. Log Forwarding

Os logs da aplicação são automaticamente encaminhados para o New Relic quando configurado.

## Arquitetura

### Componentes Modificados

1. **config/config.go**
   - Adicionados campos: `NewRelicLicenseKey`, `NewRelicAppName`, `NewRelicEnabled`

2. **cmd/app/app.go**
   - Adicionado campo `NewRelicApp` na struct `Application`

3. **main.go**
   - Inicialização do New Relic Application
   - Configuração de timeout de conexão (5 segundos)
   - Shutdown graceful ao encerrar a aplicação

4. **internal/routes/routes.go**
   - Middleware customizado `newRelicMiddleware`
   - ResponseWriter wrapper para capturar status codes

### Middleware Customizado

O middleware `newRelicMiddleware` foi implementado para integrar o New Relic com o Chi Router:

```go
func newRelicMiddleware(nrApp *newrelic.Application) func(next http.Handler) http.Handler
```

**Características:**
- Cria transações automaticamente para cada requisição
- Ignora rotas de health check
- Adiciona transação ao contexto da requisição
- Captura status codes de resposta
- Nome da transação: `{METHOD} {PATH}` (ex: `GET /api/v1/animals`)

## Uso Avançado

### Transações Customizadas em Services

Para adicionar instrumentação customizada em services:

```go
import "github.com/newrelic/go-agent/v3/newrelic"

func (s *SomeService) SomeMethod(ctx context.Context) error {
    // Obter transação do contexto
    txn := newrelic.FromContext(ctx)
    
    // Criar um segmento customizado
    segment := txn.StartSegment("CustomOperation")
    defer segment.End()
    
    // Sua lógica aqui
    
    return nil
}
```

### Reportar Erros Manualmente

```go
import "github.com/newrelic/go-agent/v3/newrelic"

func handler(w http.ResponseWriter, r *http.Request) {
    txn := newrelic.FromContext(r.Context())
    
    err := someOperation()
    if err != nil {
        txn.NoticeError(err)
        // Tratar erro
    }
}
```

### Adicionar Atributos Customizados

```go
txn := newrelic.FromContext(r.Context())
txn.AddAttribute("user_id", userID)
txn.AddAttribute("farm_id", farmID)
```

## Monitoramento no Dashboard

### Métricas Disponíveis

1. **Web Transactions**
   - Tempo de resposta por endpoint
   - Throughput (requisições por minuto)
   - Taxa de erro

2. **Database**
   - Queries mais lentas
   - Tempo de resposta do banco
   - Número de queries por transação

3. **External Services**
   - Chamadas a APIs externas
   - Tempo de resposta

4. **Errors**
   - Taxa de erro
   - Stack traces
   - Erros por endpoint

### Alertas Recomendados

Configure alertas para:
- Taxa de erro > 5%
- Tempo de resposta > 2 segundos
- Apdex score < 0.7
- Throughput anormal

## Compatibilidade

### Sentry

O New Relic coexiste com o Sentry sem conflitos:
- **Sentry**: Focado em error tracking e crash reporting
- **New Relic**: Focado em APM e performance monitoring

Ambos podem estar habilitados simultaneamente.

### Performance

- Overhead do New Relic: < 1% na maioria dos casos
- Impacto mínimo na latência das requisições
- Coleta de dados assíncrona

## Troubleshooting

### New Relic não está enviando dados

1. Verificar se `NEW_RELIC_ENABLED=true`
2. Verificar se a license key está correta
3. Verificar logs de inicialização:
   ```
   New Relic inicializado com sucesso
   ```
4. Aguardar até 5 minutos para os primeiros dados aparecerem

### Timeout de conexão

Se aparecer no log:
```
New Relic connection timeout
```

Isso é normal e não impede o funcionamento. O agente continuará tentando conectar em background.

### Desabilitar em Desenvolvimento

Para desabilitar o New Relic em desenvolvimento:

```bash
NEW_RELIC_ENABLED=false
```

Ou simplesmente não defina a variável (padrão é false).

## Referências

- [New Relic Go Agent Documentation](https://docs.newrelic.com/docs/apm/agents/go-agent/)
- [Go Agent API Guide](https://pkg.go.dev/github.com/newrelic/go-agent/v3/newrelic)
- [Best Practices](https://docs.newrelic.com/docs/apm/agents/go-agent/get-started/go-agent-best-practices/)

## Changelog

### v1.0.0 (2024)
- Integração inicial do New Relic
- Middleware customizado para Chi Router
- Monitoramento automático de rotas HTTP
- Configuração via variáveis de ambiente
- Distributed tracing habilitado
- Log forwarding habilitado

