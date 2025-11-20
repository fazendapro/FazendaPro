# Sistema de Migrations do FazendaPro

## Visão Geral

O projeto FazendaPro utiliza um sistema de migrations customizado construído sobre o GORM. Este sistema permite controlar a evolução do esquema do banco de dados de forma versionada e rastreável, garantindo que todas as alterações sejam aplicadas na ordem correta e apenas uma vez.

## Como Funciona

### Tabela de Controle: `migrations`

O sistema mantém uma tabela especial chamada `migrations` que registra todas as migrations já executadas:

```go
type Migration struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"uniqueIndex"`
    CreatedAt time.Time
}
```

Esta tabela funciona como um log de execução, permitindo que o sistema saiba quais migrations já foram aplicadas ao banco de dados.

### Função Principal: `RunMigrations()`

A função `RunMigrations()` em `internal/migrations/migrations.go` é responsável por:

1. Criar a tabela `migrations` se ela não existir
2. Iterar sobre todas as migrations definidas
3. Verificar se cada migration já foi executada
4. Executar apenas as migrations pendentes
5. Registrar cada migration executada na tabela `migrations`

```go
func RunMigrations(db *gorm.DB) error {
    // 1. Cria tabela migrations se não existir
    if err := db.AutoMigrate(&Migration{}); err != nil {
        return fmt.Errorf("error creating migrations table: %w", err)
    }

    // 2. Define lista de migrations
    migrations := []struct {
        name string
        fn   func(*gorm.DB) error
    }{
        {"001_create_users_table", createUsersTable},
        {"002_create_companies_table", createCompaniesTable},
        // ... mais migrations
    }

    // 3. Executa cada migration pendente
    for _, migration := range migrations {
        var existingMigration Migration
        if err := db.Where("name = ?", migration.name).First(&existingMigration).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                // Migration ainda não executada
                log.Printf("Executando migração: %s", migration.name)
                
                // Executa a migration
                if err := migration.fn(db); err != nil {
                    return fmt.Errorf("error executing migration %s: %w", migration.name, err)
                }
                
                // Registra na tabela migrations
                if err := db.Create(&Migration{Name: migration.name}).Error; err != nil {
                    return fmt.Errorf("error registering migration %s: %w", migration.name, err)
                }
            }
        } else {
            log.Printf("Migration %s already executed", migration.name)
        }
    }

    return nil
}
```

## Estrutura de uma Migration

Cada migration possui:

1. **Nome único**: Identificador da migration (ex: `"001_create_users_table"`)
2. **Função de execução**: Função que realiza a alteração no banco
3. **Registro**: Após execução bem-sucedida, é registrada na tabela `migrations`

### Exemplo de Migration Simples

```go
func createUsersTable(db *gorm.DB) error {
    return db.AutoMigrate(&models.User{})
}
```

Esta migration usa `AutoMigrate` do GORM para criar a tabela `users` baseada no modelo `User`.

## Tipos de Migrations

O sistema suporta diferentes tipos de alterações no banco de dados:

### 1. Criação de Tabelas

Usa `AutoMigrate` para criar novas tabelas baseadas nos modelos:

```go
func createAnimalsTable(db *gorm.DB) error {
    return db.AutoMigrate(&models.Animal{})
}
```

**Migrations deste tipo**:
- `001_create_users_table`
- `002_create_companies_table`
- `003_create_farms_table`
- `004_create_animals_table`
- `005_create_milk_collections_table`
- `006_create_reproductions_table`
- `007_create_weights_table`
- `008_create_persons_table`
- `010_create_expenses_table`
- `013_create_refresh_tokens_table`
- `018_create_user_farms_table`
- `020_create_sales_table`
- `021_create_debts_table`

### 2. Atualização de Tabelas (Adicionar Colunas)

Usa `AutoMigrate` para adicionar novas colunas a tabelas existentes:

```go
func addCompanyName(db *gorm.DB) error {
    return db.AutoMigrate(&models.Company{})
}
```

**Migrations deste tipo**:
- `009_update_users_table` - Atualiza tabela users
- `011_update_users_with_person` - Adiciona relacionamento com Person
- `012_add_company_name` - Adiciona coluna `company_name` em Company
- `013_add_farm_logo` - Adiciona coluna `logo` em Farm
- `014_add_animal_photo` - Adiciona coluna `photo` em Animal

### 3. Modificação de Tabelas (Remover Colunas)

Usa `Migrator().DropColumn()` para remover colunas:

```go
func updateAnimalsTable(db *gorm.DB) error {
    // Verifica se coluna existe antes de remover
    if db.Migrator().HasColumn(&models.Animal{}, "ear_tag_number") {
        if err := db.Migrator().DropColumn(&models.Animal{}, "ear_tag_number"); err != nil {
            return fmt.Errorf("error dropping ear_tag_number column: %w", err)
        }
    }
    
    // Remove outras colunas...
    
    // Aplica mudanças do modelo atualizado
    return db.AutoMigrate(&models.Animal{})
}
```

**Migrations deste tipo**:
- `015_update_animals_table` - Remove colunas obsoletas de Animal
- `016_update_reproductions_table` - Remove colunas obsoletas de Reproduction

### 4. Migração de Dados

Executa transformações e inserções de dados:

```go
func seedInitialData(db *gorm.DB) error {
    // Verifica se dados já existem
    var companyCount int64
    db.Model(&models.Company{}).Count(&companyCount)
    
    if companyCount > 0 {
        log.Printf("Dados iniciais já existem, pulando seed")
        return nil
    }
    
    // Cria dados iniciais
    company := &models.Company{
        CompanyName: "FazendaPro Demo",
    }
    db.Create(company)
    
    farm := &models.Farm{
        CompanyID: company.ID,
        Logo:      "",
    }
    db.Create(farm)
    
    return nil
}
```

**Migrations deste tipo**:
- `017_seed_initial_data` - Cria dados iniciais (Company e Farm demo)
- `019_migrate_users_to_user_farms` - Migra dados de users para user_farms

#### Exemplo Detalhado: `migrateUsersToUserFarms`

Esta migration migra dados existentes de uma estrutura antiga para uma nova:

```go
func migrateUsersToUserFarms(db *gorm.DB) error {
    // 1. Busca todos os usuários
    var users []models.User
    if err := db.Find(&users).Error; err != nil {
        return fmt.Errorf("error finding users: %w", err)
    }
    
    // 2. Para cada usuário, cria registro em user_farms
    for _, user := range users {
        userFarm := &models.UserFarm{
            UserID:    user.ID,
            FarmID:    user.FarmID,
            IsPrimary: true,
        }
        
        // 3. Verifica se já existe para evitar duplicatas
        var existingUserFarm models.UserFarm
        if err := db.Where("user_id = ? AND farm_id = ?", user.ID, user.FarmID)
                    .First(&existingUserFarm).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                // 4. Cria apenas se não existir
                if err := db.Create(userFarm).Error; err != nil {
                    return fmt.Errorf("error creating user farm for user %d: %w", user.ID, err)
                }
            }
        }
    }
    
    return nil
}
```

## Execução de Migrations

### Execução Automática na Inicialização

As migrations são executadas automaticamente quando a aplicação inicia, no arquivo `main.go`:

```go
func main() {
    // ... configuração da aplicação ...
    
    db, err := repository.NewDatabase(cfg)
    if err != nil {
        // ... tratamento de erro ...
    }
    
    // Executa migrations automaticamente
    if err := migrations.RunMigrations(db.DB); err != nil {
        app.Logger.Printf("WARNING: Erro ao executar migrações: %v", err)
        app.Logger.Println("Continuando sem migrações...")
    }
    
    // ... resto da aplicação ...
}
```

### Execução Manual

Também é possível executar migrations manualmente através de um comando:

```go
func main() {
    if len(os.Args) > 1 && os.Args[1] == "migrate" {
        runMigrations()
        return
    }
    // ... resto do código ...
}

func runMigrations() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Erro ao carregar configuração:", err)
    }
    
    db, err := repository.NewDatabase(cfg)
    if err != nil {
        log.Fatal("Erro ao conectar ao banco:", err)
    }
    defer db.Close()
    
    log.Println("Executando migrações...")
    if err := migrations.RunMigrations(db.DB); err != nil {
        log.Fatal("Erro ao executar migrações:", err)
    }
    log.Println("Migrações executadas com sucesso!")
}
```

**Como executar**:
```bash
go run main.go migrate
```

## Rollback de Migrations

O sistema possui uma função `RollbackMigrations()` que permite reverter migrations:

```go
func RollbackMigrations(db *gorm.DB, steps int) error {
    // 1. Busca as últimas N migrations executadas
    var migrations []Migration
    if err := db.Order("id desc").Limit(steps).Find(&migrations).Error; err != nil {
        return fmt.Errorf("error searching migrations: %w", err)
    }
    
    // 2. Para cada migration, executa a reversão
    for _, migration := range migrations {
        log.Printf("Reverting migration: %s", migration.Name)
        
        // 3. Executa ação de reversão específica para cada migration
        switch migration.Name {
        case "001_create_users_table":
            if err := db.Migrator().DropTable(&models.User{}); err != nil {
                return fmt.Errorf("error reverting migration %s: %w", migration.Name, err)
            }
        // ... outros casos ...
        }
        
        // 4. Remove registro da tabela migrations
        if err := db.Delete(&migration).Error; err != nil {
            return fmt.Errorf("error removing migration record %s: %w", migration.Name, err)
        }
    }
    
    return nil
}
```

**Nota**: A função de rollback precisa ter a lógica de reversão implementada manualmente para cada migration. Nem todas as migrations têm rollback implementado.

## Lista Completa de Migrations

Aqui está a lista completa das 21 migrations atuais do projeto, na ordem de execução:

| # | Nome | Descrição |
|---|------|-----------|
| 001 | `create_users_table` | Cria tabela de usuários |
| 002 | `create_companies_table` | Cria tabela de empresas |
| 003 | `create_farms_table` | Cria tabela de fazendas |
| 004 | `create_animals_table` | Cria tabela de animais |
| 005 | `create_milk_collections_table` | Cria tabela de coletas de leite |
| 006 | `create_reproductions_table` | Cria tabela de reproduções |
| 007 | `create_weights_table` | Cria tabela de pesos |
| 008 | `create_persons_table` | Cria tabela de pessoas |
| 009 | `update_users_table` | Atualiza tabela de usuários |
| 010 | `create_expenses_table` | Cria tabela de despesas |
| 011 | `update_users_with_person` | Adiciona relacionamento User-Person |
| 012 | `add_company_name` | Adiciona coluna `company_name` em Company |
| 013 | `create_refresh_tokens_table` | Cria tabela de refresh tokens |
| 013 | `add_farm_logo` | Adiciona coluna `logo` em Farm |
| 014 | `add_animal_photo` | Adiciona coluna `photo` em Animal |
| 015 | `update_animals_table` | Remove colunas obsoletas de Animal |
| 016 | `update_reproductions_table` | Remove colunas obsoletas de Reproduction |
| 017 | `seed_initial_data` | Cria dados iniciais (Company e Farm demo) |
| 018 | `create_user_farms_table` | Cria tabela de relacionamento User-Farm |
| 019 | `migrate_users_to_user_farms` | Migra dados de users para user_farms |
| 020 | `create_sales_table` | Cria tabela de vendas |
| 021 | `create_debts_table` | Cria tabela de dívidas |

**Observação**: Note que há duas migrations com número 013 (`create_refresh_tokens_table` e `add_farm_logo`). Isso pode causar confusão, mas como o sistema usa o nome como identificador único, ambas funcionam corretamente.

## Boas Práticas

### 1. Nomes Descritivos

Use nomes que descrevam claramente o que a migration faz:
- ✅ `001_create_users_table`
- ❌ `001_migration`

### 2. Idempotência

As migrations devem ser idempotentes quando possível. O sistema já garante que cada migration execute apenas uma vez, mas é bom verificar dados existentes:

```go
func seedInitialData(db *gorm.DB) error {
    var companyCount int64
    db.Model(&models.Company{}).Count(&companyCount)
    
    if companyCount > 0 {
        return nil // Já existe, não faz nada
    }
    
    // Cria dados...
}
```

### 3. Verificação de Colunas

Ao remover colunas, sempre verifique se elas existem:

```go
if db.Migrator().HasColumn(&models.Animal{}, "ear_tag_number") {
    db.Migrator().DropColumn(&models.Animal{}, "ear_tag_number")
}
```

### 4. Ordem de Execução

As migrations são executadas na ordem em que aparecem no array. Certifique-se de que dependências sejam criadas antes de serem usadas.

### 5. Logs

Use logs para rastrear a execução:

```go
log.Printf("Executando migração: %s", migration.name)
// ... execução ...
log.Printf("Migration %s executed successfully", migration.name)
```

## Troubleshooting

### Migration já executada mas precisa reexecutar

Se precisar reexecutar uma migration (por exemplo, em desenvolvimento):

1. Remova o registro da tabela `migrations`:
```sql
DELETE FROM migrations WHERE name = 'nome_da_migration';
```

2. Execute a aplicação novamente ou rode `go run main.go migrate`

### Erro ao executar migration

Se uma migration falhar:

1. O sistema para a execução e retorna erro
2. Verifique os logs para identificar o problema
3. Corrija o problema na migration
4. Execute novamente

### Migration com número duplicado

O sistema usa o **nome** como identificador único, não o número. Portanto, duas migrations podem ter o mesmo número prefixo, mas devem ter nomes diferentes.

## Conclusão

O sistema de migrations do FazendaPro oferece controle versionado sobre o esquema do banco de dados, garantindo que todas as alterações sejam aplicadas de forma consistente e rastreável. É uma solução simples mas eficaz para gerenciar a evolução do banco de dados em projetos Go com GORM.

