# FazendaPro Database Schema

Este diretório contém todos os arquivos relacionados ao esquema do banco de dados do sistema FazendaPro - uma solução para gestão de fazendas de gado leiteiro.

## Estrutura do Diretório

```
database/
├── migrations/          # Arquivos de migração individuais
├── schema/             # Schema completo consolidado
├── seeds/              # Dados de exemplo para testes
└── README.md          # Esta documentação
```

## Funcionalidades Implementadas

O schema foi desenvolvido baseado nos requisitos funcionais (RF) documentados no projeto:

### Tabelas Principais

| Tabela | Descrição | Requisitos |
|--------|-----------|------------|
| `users` | Usuários do sistema (autenticação) | RF01 |
| `animals` | Animais da fazenda (core do sistema) | RF02, RF03 |
| `lots` | Lotes de produção por performance | RF08 |
| `vaccines` | Catálogo de vacinas | RF11 |
| `animal_vaccines` | Registros de vacinação | RF06 |
| `weight_records` | Histórico de peso dos animais | RF07 |
| `pregnancy_records` | Controle de gestação | RF09 |
| `sales` | Vendas de animais | RF10 |
| `suppliers` | Fornecedores | RF12 |

## Relacionamentos

### Diagrama Conceitual

```
Users (1) -----> (N) Animals
Animals (N) --> (1) Lots
Animals (1) -----> (N) Weight_Records
Animals (1) -----> (N) Animal_Vaccines
Animals (1) -----> (N) Pregnancy_Records
Animals (1) -----> (N) Sales
Vaccines (1) ----> (N) Animal_Vaccines
Users (1) -------> (N) Vaccines
Users (1) -------> (N) Suppliers

-- Relacionamentos hierárquicos
Animals (mother_id) --> Animals (id)
Animals (father_id) --> Animals (id)
Pregnancy_Records (father_id) --> Animals (id)
```

## Como Usar

### 1. Schema Completo (Recomendado para desenvolvimento)

```sql
-- Execute o schema completo para criar todas as tabelas
source database/schema/complete_schema.sql;
```

### 2. Migrações Individuais

```sql
-- Execute as migrações em ordem numérica
source database/migrations/001_create_users_table.sql;
source database/migrations/002_create_lots_table.sql;
source database/migrations/003_create_animals_table.sql;
-- ... continue até 009
```

### 3. Dados de Exemplo

```sql
-- Após criar as tabelas, insira dados de exemplo para testes
source database/seeds/sample_data.sql;
```

## Características Técnicas

### Engine e Charset
- **Engine**: InnoDB (suporte a transações e foreign keys)
- **Charset**: utf8mb4_unicode_ci (suporte completo a caracteres especiais)

### Indexação
- Índices automáticos em todas as foreign keys
- Índices em campos de consulta frequente (email, ear_tag, dates)
- Índices compostos para consultas específicas

### Soft Delete
- Todas as tabelas implementam soft delete via campo `deleted_at`
- Permite manter histórico mesmo após "exclusão"

### Timestamps
- `created_at` e `updated_at` automáticos em todas as tabelas
- `updated_at` atualizado automaticamente via trigger MySQL

## Recursos Especiais

### 1. Controle Genealógico
```sql
-- Tabela animals permite rastreamento completo da genealogia
mother_id INT -- Referência para a mãe
father_id INT -- Referência para o pai
```

### 2. Controle de Lotes Automático (RF08)
```sql
-- Lotes definem faixas de produção de leite
min_milk_production DECIMAL(8,2)
max_milk_production DECIMAL(8,2)
```

### 3. Sistema de Notificações (RF09)
```sql
-- Controle de notificações WhatsApp para gestação
notification_sent BOOLEAN
notification_date TIMESTAMP
```

### 4. Rastreabilidade Completa
- Todas as ações são rastreáveis via timestamps
- Relacionamentos preservam histórico familiar
- Registros de peso mantêm evolução temporal
- Vacinações com controle de lotes e validade

## Exemplos de Consultas

### Animais em Gestação Próximos do Parto
```sql
SELECT 
    a.name,
    a.ear_tag,
    pr.expected_birth_date,
    DATEDIFF(pr.expected_birth_date, CURDATE()) as days_to_birth
FROM animals a
JOIN pregnancy_records pr ON a.id = pr.animal_id
WHERE pr.pregnancy_status = 'confirmed'
  AND pr.expected_birth_date BETWEEN CURDATE() AND DATE_ADD(CURDATE(), INTERVAL 20 DAY)
  AND pr.notification_sent = FALSE;
```

### Histórico Familiar de um Animal
```sql
WITH RECURSIVE family_tree AS (
    -- Animal base
    SELECT id, name, ear_tag, mother_id, father_id, 0 as generation
    FROM animals 
    WHERE ear_tag = 'V001'
    
    UNION ALL
    
    -- Ascendentes
    SELECT a.id, a.name, a.ear_tag, a.mother_id, a.father_id, ft.generation + 1
    FROM animals a
    JOIN family_tree ft ON a.id IN (ft.mother_id, ft.father_id)
    WHERE ft.generation < 3
)
SELECT * FROM family_tree ORDER BY generation, name;
```

### Relatório de Produção por Lote
```sql
SELECT 
    l.name as lote,
    COUNT(a.id) as total_animais,
    AVG(a.daily_milk_production) as producao_media,
    SUM(a.daily_milk_production) as producao_total
FROM lots l
LEFT JOIN animals a ON l.id = a.lot_id AND a.status = 'active'
GROUP BY l.id, l.name
ORDER BY producao_total DESC;
```

## Configuração MySQL Recomendada

```sql
-- Configurações recomendadas para melhor performance
SET innodb_buffer_pool_size = 256M;
SET max_connections = 100;
SET innodb_log_file_size = 64M;
```

## Backup e Manutenção

### Backup Completo
```bash
mysqldump -u usuario -p fazendapro_db > backup_$(date +%Y%m%d).sql
```

### Limpeza de Soft Deletes (executar periodicamente)
```sql
-- Remove registros marcados como deletados há mais de 1 ano
DELETE FROM animals WHERE deleted_at < DATE_SUB(NOW(), INTERVAL 1 YEAR);
-- Repita para outras tabelas conforme necessário
```

## Observações de Segurança

1. **Senhas**: Sempre use hash bcrypt para senhas (implementado no seed)
2. **Sanitização**: Todas as entradas devem ser sanitizadas no backend
3. **Backup**: Configure backup automático diário
4. **SSL**: Use conexões SSL em produção
5. **Usuários**: Crie usuários específicos com permissões limitadas

## Próximos Passos

- [ ] Implementar triggers para atualização automática de lotes
- [ ] Criar views para consultas frequentes
- [ ] Implementar stored procedures para operações complexas
- [ ] Adicionar tabela de logs para auditoria
- [ ] Configurar replicação para backup em tempo real

## Suporte

Para dúvidas sobre o schema ou problemas na implementação, consulte:
- Documentação completa em `docs/full_rfc.md`
- Diagramas ER em `docs/images/entityRelationshipDiagram.drawio.png`
- Issues no GitHub do projeto