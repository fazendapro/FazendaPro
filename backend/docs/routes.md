# Documentação das Rotas da API FazendaPro

## Visão Geral

Este documento descreve todas as rotas da API REST do FazendaPro, organizadas por funcionalidade. A API utiliza o router Chi e segue padrões RESTful.

## Estrutura Base

A função `SetupRoutes()` em `internal/routes/routes.go` configura todas as rotas da aplicação.

### Router Principal

```go
r := chi.NewRouter()
```

## Middlewares Globais

### 1. Sentry Middleware (Opcional)

**Condição**: Aplicado apenas se `SentryDSN` estiver configurado.

**Função**: Captura erros e panics para monitoramento.

```go
if cfg.SentryDSN != "" {
    sentryHandler := sentryhttp.New(sentryhttp.Options{
        Repanic: true,
    })
    r.Use(sentryHandler.Handle)
}
```

### 2. CORS Middleware

**Função**: Gerencia Cross-Origin Resource Sharing.

**Aplicado**: Em todas as rotas.

```go
r.Use(middleware.CORSMiddleware(cfg))
```

### 3. Logging Middleware

**Função**: Registra todas as requisições HTTP.

**Aplicado**: Em todas as rotas.

```go
r.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        app.Logger.Printf("Request: %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
})
```

## Rotas Públicas (Sem Autenticação)

### Health Check

**Endpoint**: `GET /health`

**Descrição**: Verifica se a API está funcionando.

**Resposta**: `200 OK` com mensagem de status.

---

### Root

**Endpoint**: `GET /`

**Descrição**: Endpoint raiz da API.

**Resposta**: `200 OK` com mensagem "FazendaPro API is running!".

---

### Inicialização de Dados

**Endpoint**: `POST /init-data`

**Descrição**: Cria dados iniciais (Company e Farm demo) se não existirem.

**Resposta**:
- `200 OK`: Dados criados ou já existentes
- `500 Internal Server Error`: Erro ao criar dados

**Nota**: Esta rota não requer autenticação e deve ser usada apenas em desenvolvimento/primeira instalação.

---

## Rotas de Autenticação (`/api/v1/auth`)

**Base Path**: `/api/v1/auth`

**Autenticação**: Não requerida (endpoint público)

### Login

**Endpoint**: `POST /api/v1/auth/login`

**Handler**: `AuthHandler.Login`

**Descrição**: Autentica um usuário e retorna tokens.

**Body**:
```json
{
  "email": "usuario@example.com",
  "password": "senha123"
}
```

**Resposta**: Access token, refresh token e dados do usuário.

---

### Register

**Endpoint**: `POST /api/v1/auth/register`

**Handler**: `AuthHandler.Register`

**Descrição**: Registra um novo usuário e faz login automático.

**Body**:
```json
{
  "user": {
    "farm_id": 1
  },
  "person": {
    "first_name": "João",
    "last_name": "Silva",
    "email": "joao@example.com",
    "phone": "11999999999"
  }
}
```

**Resposta**: Access token, refresh token e dados do usuário.

---

### Refresh Token

**Endpoint**: `POST /api/v1/auth/refresh`

**Handler**: `AuthHandler.RefreshToken`

**Descrição**: Renova o access token usando um refresh token.

**Body**:
```json
{
  "refresh_token": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Resposta**: Novo access token.

---

### Logout

**Endpoint**: `POST /api/v1/auth/logout`

**Handler**: `AuthHandler.Logout`

**Descrição**: Invalida um refresh token.

**Body**:
```json
{
  "refresh_token": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Resposta**: Confirmação de logout.

---

## Rotas de Usuários (`/api/v1/users`)

**Base Path**: `/api/v1/users`

**Autenticação**: Requerida (middleware Auth)

### Criar Usuário

**Endpoint**: `POST /api/v1/users`

**Handler**: `UserHandler.CreateUser`

**Descrição**: Cria um novo usuário no sistema.

---

### Buscar Usuário

**Endpoint**: `GET /api/v1/users?email={email}`

**Handler**: `UserHandler.GetUser`

**Descrição**: Busca um usuário por email.

**Query Parameters**:
- `email` (obrigatório): Email do usuário

---

## Rotas de Fazendas (`/api/v1/farms`)

**Base Path**: `/api/v1/farms`

**Autenticação**: Requerida

### Listar Fazendas do Usuário

**Endpoint**: `GET /api/v1/farms/user`

**Handler**: `FarmSelectionHandler.GetUserFarms`

**Descrição**: Lista todas as fazendas de um usuário.

**Resposta**: Lista de fazendas com flag de auto-seleção.

---

### Selecionar Fazenda

**Endpoint**: `POST /api/v1/farms/select`

**Handler**: `FarmSelectionHandler.SelectFarm`

**Descrição**: Seleciona uma fazenda para o usuário trabalhar.

**Body**:
```json
{
  "farm_id": 1
}
```

---

## Rotas de Animais (`/api/v1/animals`)

**Base Path**: `/api/v1/animals`

**Autenticação**: Requerida

### Criar Animal

**Endpoint**: `POST /api/v1/animals`

**Handler**: `AnimalHandler.CreateAnimal`

**Descrição**: Cria um novo animal na fazenda.

---

### Buscar Animal

**Endpoint**: `GET /api/v1/animals?id={id}`

**Handler**: `AnimalHandler.GetAnimal`

**Descrição**: Busca um animal por ID.

**Query Parameters**:
- `id` (obrigatório): ID do animal

---

### Listar Animais da Fazenda

**Endpoint**: `GET /api/v1/animals/farm?farmId={farmId}`

**Handler**: `AnimalHandler.GetAnimalsByFarm`

**Descrição**: Lista todos os animais de uma fazenda.

**Query Parameters**:
- `farmId` (obrigatório): ID da fazenda

---

### Listar Animais por Sexo

**Endpoint**: `GET /api/v1/animals/sex?farmId={farmId}&sex={sex}`

**Handler**: `AnimalHandler.GetAnimalsBySex`

**Descrição**: Lista animais filtrados por sexo.

**Query Parameters**:
- `farmId` (obrigatório): ID da fazenda
- `sex` (obrigatório): 0 = Fêmea, 1 = Macho

---

### Atualizar Animal

**Endpoint**: `PUT /api/v1/animals`

**Handler**: `AnimalHandler.UpdateAnimal`

**Descrição**: Atualiza dados de um animal.

---

### Deletar Animal

**Endpoint**: `DELETE /api/v1/animals?id={id}`

**Handler**: `AnimalHandler.DeleteAnimal`

**Descrição**: Remove um animal do sistema.

**Query Parameters**:
- `id` (obrigatório): ID do animal

---

### Upload de Foto do Animal

**Endpoint**: `POST /api/v1/animals/photo`

**Handler**: `AnimalHandler.UploadAnimalPhoto`

**Descrição**: Faz upload de uma foto para um animal.

**Form Data**:
- `animal_id` (obrigatório): ID do animal
- `photo` (obrigatório): Arquivo de imagem (máx. 10MB)

---

## Rotas de Coleta de Leite (`/api/v1/milk-collections`)

**Base Path**: `/api/v1/milk-collections`

**Autenticação**: Requerida

### Criar Coleta de Leite

**Endpoint**: `POST /api/v1/milk-collections`

**Handler**: `MilkCollectionHandler.CreateMilkCollection`

**Descrição**: Registra uma nova coleta de leite.

---

### Atualizar Coleta de Leite

**Endpoint**: `PUT /api/v1/milk-collections/{id}`

**Handler**: `MilkCollectionHandler.UpdateMilkCollection`

**Descrição**: Atualiza uma coleta de leite existente.

**Path Parameters**:
- `id` (obrigatório): ID da coleta

---

### Listar Coletas por Fazenda

**Endpoint**: `GET /api/v1/milk-collections/farm/{farmId}?start_date={date}&end_date={date}`

**Handler**: `MilkCollectionHandler.GetMilkCollectionsByFarmID`

**Descrição**: Lista coletas de leite de uma fazenda, opcionalmente filtradas por período.

**Path Parameters**:
- `farmId` (obrigatório): ID da fazenda

**Query Parameters**:
- `start_date` (opcional): Data inicial (formato: YYYY-MM-DD)
- `end_date` (opcional): Data final (formato: YYYY-MM-DD)

---

### Listar Coletas por Animal

**Endpoint**: `GET /api/v1/milk-collections/animal/{animalId}`

**Handler**: `MilkCollectionHandler.GetMilkCollectionsByAnimalID`

**Descrição**: Lista todas as coletas de leite de um animal.

**Path Parameters**:
- `animalId` (obrigatório): ID do animal

---

### Top Produtoras de Leite

**Endpoint**: `GET /api/v1/milk-collections/top-producers?farmId={id}&limit={limit}&periodDays={days}`

**Handler**: `MilkCollectionHandler.GetTopMilkProducers`

**Descrição**: Retorna as maiores produtoras de leite.

**Query Parameters**:
- `farmId` (obrigatório): ID da fazenda
- `limit` (opcional, padrão: 10): Número máximo de resultados
- `periodDays` (opcional, padrão: 30): Período em dias para análise

---

## Rotas de Reprodução (`/api/v1/reproductions`)

**Base Path**: `/api/v1/reproductions`

**Autenticação**: Requerida

### Criar Registro de Reprodução

**Endpoint**: `POST /api/v1/reproductions`

**Handler**: `ReproductionHandler.CreateReproduction`

**Descrição**: Cria um novo registro de reprodução.

---

### Buscar Registro de Reprodução

**Endpoint**: `GET /api/v1/reproductions?id={id}`

**Handler**: `ReproductionHandler.GetReproduction`

**Descrição**: Busca um registro de reprodução por ID.

**Query Parameters**:
- `id` (obrigatório): ID do registro

---

### Buscar Reprodução por Animal

**Endpoint**: `GET /api/v1/reproductions/animal?animalId={id}`

**Handler**: `ReproductionHandler.GetReproductionByAnimal`

**Descrição**: Busca registro de reprodução de um animal.

**Query Parameters**:
- `animalId` (obrigatório): ID do animal

---

### Listar Reproduções por Fazenda

**Endpoint**: `GET /api/v1/reproductions/farm?farmId={id}`

**Handler**: `ReproductionHandler.GetReproductionsByFarm`

**Descrição**: Lista todos os registros de reprodução de uma fazenda.

**Query Parameters**:
- `farmId` (obrigatório): ID da fazenda

---

### Listar Reproduções por Fase

**Endpoint**: `GET /api/v1/reproductions/phase?phase={phase}`

**Handler**: `ReproductionHandler.GetReproductionsByPhase`

**Descrição**: Lista registros de reprodução por fase.

**Query Parameters**:
- `phase` (obrigatório): Número da fase

---

### Próximas Vacas a Parir

**Endpoint**: `GET /api/v1/reproductions/next-to-calve?farmId={id}`

**Handler**: `ReproductionHandler.GetNextToCalve`

**Descrição**: Lista animais próximos a parir, ordenados por data esperada.

**Query Parameters**:
- `farmId` (obrigatório): ID da fazenda

---

### Atualizar Registro de Reprodução

**Endpoint**: `PUT /api/v1/reproductions`

**Handler**: `ReproductionHandler.UpdateReproduction`

**Descrição**: Atualiza um registro de reprodução completo.

---

### Atualizar Fase de Reprodução

**Endpoint**: `PUT /api/v1/reproductions/phase`

**Handler**: `ReproductionHandler.UpdateReproductionPhase`

**Descrição**: Atualiza apenas a fase de reprodução de um animal.

---

### Deletar Registro de Reprodução

**Endpoint**: `DELETE /api/v1/reproductions?id={id}`

**Handler**: `ReproductionHandler.DeleteReproduction`

**Descrição**: Remove um registro de reprodução.

**Query Parameters**:
- `id` (obrigatório): ID do registro

---

## Rotas de Fazenda (`/api/v1/farm`)

**Base Path**: `/api/v1/farm`

**Autenticação**: Requerida

### Buscar Fazenda

**Endpoint**: `GET /api/v1/farm?id={id}`

**Handler**: `FarmHandler.GetFarm`

**Descrição**: Busca uma fazenda por ID, incluindo dados da empresa.

**Query Parameters**:
- `id` (obrigatório): ID da fazenda

---

### Atualizar Fazenda

**Endpoint**: `PUT /api/v1/farm?id={id}`

**Handler**: `FarmHandler.UpdateFarm`

**Descrição**: Atualiza dados de uma fazenda (principalmente logo).

**Query Parameters**:
- `id` (obrigatório): ID da fazenda

---

## Rotas de Vendas (`/api/v1/sales`)

**Base Path**: `/api/v1/sales`

**Autenticação**: Requerida

### Criar Venda

**Endpoint**: `POST /api/v1/sales`

**Handler**: `SaleChiHandler.CreateSale`

**Descrição**: Registra uma nova venda de animal.

---

### Listar Vendas da Fazenda

**Endpoint**: `GET /api/v1/sales`

**Handler**: `SaleChiHandler.GetSalesByFarm`

**Descrição**: Lista todas as vendas de uma fazenda.

---

### Histórico de Vendas

**Endpoint**: `GET /api/v1/sales/history`

**Handler**: `SaleChiHandler.GetSalesHistory`

**Descrição**: Retorna histórico completo de vendas.

---

### Estatísticas Mensais

**Endpoint**: `GET /api/v1/sales/monthly-stats`

**Handler**: `SaleChiHandler.GetMonthlySalesStats`

**Descrição**: Retorna estatísticas de vendas do mês atual.

---

### Dados Mensais de Vendas e Compras

**Endpoint**: `GET /api/v1/sales/monthly-data?months={n}`

**Handler**: `SaleChiHandler.GetMonthlySalesAndPurchases`

**Descrição**: Retorna dados mensais de vendas e compras.

**Query Parameters**:
- `months` (opcional, padrão: 12, máximo: 24): Número de meses

---

### Estatísticas Gerais

**Endpoint**: `GET /api/v1/sales/overview`

**Handler**: `SaleChiHandler.GetOverviewStats`

**Descrição**: Retorna estatísticas gerais da fazenda.

---

### Vendas por Período

**Endpoint**: `GET /api/v1/sales/date-range?start_date={date}&end_date={date}`

**Handler**: `SaleChiHandler.GetSalesByDateRange`

**Descrição**: Lista vendas em um período específico.

**Query Parameters**:
- `start_date` (obrigatório): Data inicial (formato: YYYY-MM-DD)
- `end_date` (obrigatório): Data final (formato: YYYY-MM-DD)

---

### Buscar Venda por ID

**Endpoint**: `GET /api/v1/sales/{id}`

**Handler**: `SaleChiHandler.GetSaleByID`

**Descrição**: Busca uma venda específica por ID.

**Path Parameters**:
- `id` (obrigatório): ID da venda

---

### Atualizar Venda

**Endpoint**: `PUT /api/v1/sales/{id}`

**Handler**: `SaleChiHandler.UpdateSale`

**Descrição**: Atualiza dados de uma venda.

**Path Parameters**:
- `id` (obrigatório): ID da venda

---

### Deletar Venda

**Endpoint**: `DELETE /api/v1/sales/{id}`

**Handler**: `SaleChiHandler.DeleteSale`

**Descrição**: Remove uma venda do sistema.

**Path Parameters**:
- `id` (obrigatório): ID da venda

---

## Rotas de Vendas por Animal (`/api/v1/animals/{animal_id}/sales`)

**Base Path**: `/api/v1/animals/{animal_id}/sales`

**Autenticação**: Requerida

### Listar Vendas de um Animal

**Endpoint**: `GET /api/v1/animals/{animal_id}/sales`

**Handler**: `SaleChiHandler.GetSalesByAnimal`

**Descrição**: Lista todas as vendas de um animal específico.

**Path Parameters**:
- `animal_id` (obrigatório): ID do animal

---

## Rotas de Dívidas (`/debts`)

**Base Path**: `/debts`

**Autenticação**: Não requerida (fora do grupo `/api/v1`)

**Nota**: Estas rotas estão fora do grupo `/api/v1` e não requerem autenticação.

### Criar Dívida

**Endpoint**: `POST /debts`

**Handler**: `DebtHandler.CreateDebt`

**Descrição**: Cria uma nova dívida.

---

### Listar Dívidas

**Endpoint**: `GET /debts?page={page}&limit={limit}&year={year}&month={month}`

**Handler**: `DebtHandler.GetDebts`

**Descrição**: Lista dívidas com paginação e filtros.

**Query Parameters**:
- `page` (opcional, padrão: 1): Número da página
- `limit` (opcional, padrão: 10): Itens por página
- `year` (opcional): Filtrar por ano
- `month` (opcional, 1-12): Filtrar por mês

---

### Deletar Dívida

**Endpoint**: `DELETE /debts/{id}`

**Handler**: `DebtHandler.DeleteDebt`

**Descrição**: Remove uma dívida.

**Path Parameters**:
- `id` (obrigatório): ID da dívida

---

### Total por Pessoa

**Endpoint**: `GET /debts/total-by-person?year={year}&month={month}`

**Handler**: `DebtHandler.GetTotalByPerson`

**Descrição**: Calcula total de dívidas por pessoa em um mês.

**Query Parameters**:
- `year` (obrigatório): Ano
- `month` (obrigatório, 1-12): Mês

---

## Autenticação

### Middleware de Autenticação

A maioria das rotas utiliza o middleware `Auth`:

```go
r.Use(middleware.Auth(cfg.JWTSecret))
```

### Como Autenticar

1. **Fazer Login**: `POST /api/v1/auth/login`
2. **Obter Tokens**: Receber `access_token` e `refresh_token`
3. **Usar Access Token**: Incluir no header de requisições:

```http
Authorization: Bearer {access_token}
```

### Renovação de Token

Quando o access token expirar (24 horas), use o refresh token:

```http
POST /api/v1/auth/refresh
{
  "refresh_token": "{refresh_token}"
}
```

---

## Códigos de Status HTTP

### Sucesso
- `200 OK`: Operação bem-sucedida
- `201 Created`: Recurso criado com sucesso

### Erro do Cliente
- `400 Bad Request`: Requisição inválida
- `401 Unauthorized`: Não autenticado ou token inválido
- `403 Forbidden`: Não autorizado
- `404 Not Found`: Recurso não encontrado
- `405 Method Not Allowed`: Método HTTP não permitido

### Erro do Servidor
- `500 Internal Server Error`: Erro interno do servidor

---

## Formato de Datas

Todas as datas devem estar no formato: **YYYY-MM-DD** (ex: `2024-01-15`)

Exemplo:
```json
{
  "sale_date": "2024-01-15"
}
```

---

## Paginação

Rotas que suportam paginação usam os parâmetros:
- `page`: Número da página (começa em 1)
- `limit`: Itens por página

Exemplo:
```
GET /debts?page=1&limit=10
```

---

## Resumo de Rotas por Grupo

| Grupo | Base Path | Autenticação | Métodos |
|-------|-----------|--------------|---------|
| Públicas | `/`, `/health`, `/init-data` | Não | 3 |
| Autenticação | `/api/v1/auth` | Não | 4 |
| Usuários | `/api/v1/users` | Sim | 2 |
| Fazendas | `/api/v1/farms` | Sim | 2 |
| Animais | `/api/v1/animals` | Sim | 7 |
| Coleta de Leite | `/api/v1/milk-collections` | Sim | 5 |
| Reprodução | `/api/v1/reproductions` | Sim | 9 |
| Fazenda (singular) | `/api/v1/farm` | Sim | 2 |
| Vendas | `/api/v1/sales` | Sim | 12 |
| Vendas por Animal | `/api/v1/animals/{id}/sales` | Sim | 1 |
| Dívidas | `/debts` | Não | 4 |

**Total**: ~51 endpoints

---

## Notas Importantes

1. **Versão da API**: A maioria das rotas está sob `/api/v1`
2. **Dívidas**: As rotas de dívidas estão fora do grupo `/api/v1` e não requerem autenticação
3. **Farm ID**: Muitas rotas obtêm o `farm_id` do contexto (setado pelo middleware de autenticação)
4. **Path vs Query**: Algumas rotas usam path parameters (`/{id}`), outras query parameters (`?id={id}`)
5. **Consistência**: Nem todas as rotas seguem exatamente o mesmo padrão (algumas melhorias podem ser feitas)

---

## Exemplo de Fluxo Completo

### 1. Registrar Usuário
```http
POST /api/v1/auth/register
```

### 2. Fazer Login
```http
POST /api/v1/auth/login
```

### 3. Listar Fazendas
```http
GET /api/v1/farms/user
Authorization: Bearer {token}
```

### 4. Selecionar Fazenda
```http
POST /api/v1/farms/select
Authorization: Bearer {token}
```

### 5. Criar Animal
```http
POST /api/v1/animals
Authorization: Bearer {token}
```

### 6. Listar Animais
```http
GET /api/v1/animals/farm?farmId=1
Authorization: Bearer {token}
```

---

## Referências

- Para detalhes de cada handler, consulte a documentação em `docs/handlers/`
- Para entender a arquitetura, consulte `docs/arquitetura.md`
- Para informações sobre autenticação, consulte `docs/handlers/auth.md`

