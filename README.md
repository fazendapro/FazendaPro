# FazendaPro API

API backend para o projeto FazendaPro.

## 🚀 Como Usar

### Desenvolvimento (Docker)
```bash
make dev
```

### Produção
```bash
make prod
```

### Apenas Aplicação (banco já rodando)
```bash
make quick
```

## 📁 Estrutura

- `env.development` - Configurações para desenvolvimento
- `env.production` - Configurações para produção
- `scripts/dev.sh` - Script de desenvolvimento
- `docker-compose.yml` - Configuração do Docker

## 🔧 Configuração

### Desenvolvimento
- Usa Docker PostgreSQL
- Porta: 8080
- Banco: localhost:5432

### Produção
- Conecta em 127.0.0.1:5432
- Requer Cloud SQL Proxy ou conexão direta
- Configure as credenciais em `env.production`

## 📊 Monitoramento

### Sentry
Sistema de error tracking e crash reporting já integrado.

### New Relic
APM (Application Performance Monitoring) para monitoramento de performance, transações HTTP e logs.

**Configuração:**
```bash
NEW_RELIC_ENABLED=true
NEW_RELIC_LICENSE_KEY=sua_chave_aqui
NEW_RELIC_APP_NAME=fazendapro-backend
```

Veja a [documentação completa do New Relic](docs/newrelic.md) para mais detalhes.

## 📚 Documentação

- [Arquitetura](docs/arquitetura.md)
- [Bibliotecas](docs/bibliotecas.md)
- [Handlers](docs/handlers/README.md)
- [Middleware](docs/middleware/README.md)
- [Migrações](docs/migrations.md)
- [Rotas](docs/routes.md)
- [Testes](docs/testes.md)
- [New Relic](docs/newrelic.md) 