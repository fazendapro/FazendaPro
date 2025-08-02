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