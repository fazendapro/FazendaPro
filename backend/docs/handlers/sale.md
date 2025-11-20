# Handler: Sale

## Visão Geral

O `SaleChiHandler` gerencia todas as operações relacionadas a vendas de animais, incluindo CRUD completo, estatísticas e relatórios.

## Estrutura

```go
type SaleChiHandler struct {
    service service.SaleService
}
```

**Nota**: Usa interface `SaleService` ao invés de ponteiro.

## DTOs

### CreateSaleRequest
```go
type CreateSaleRequest struct {
    AnimalID  uint    `json:"animal_id"`
    BuyerName string  `json:"buyer_name"`
    Price     float64 `json:"price"`
    SaleDate  string  `json:"sale_date"`
    Notes     string  `json:"notes"`
}
```

### UpdateSaleRequest
```go
type UpdateSaleRequest struct {
    BuyerName string  `json:"buyer_name"`
    Price     float64 `json:"price"`
    SaleDate  string  `json:"sale_date"`
    Notes     string  `json:"notes"`
}
```

### SaleResponse
```go
type SaleResponse struct {
    ID        uint           `json:"id"`
    AnimalID  uint           `json:"animal_id"`
    FarmID    uint           `json:"farm_id"`
    BuyerName string         `json:"buyer_name"`
    Price     float64        `json:"price"`
    SaleDate  time.Time      `json:"sale_date"`
    Notes     string         `json:"notes"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    Animal    *models.Animal `json:"animal,omitempty"`
}
```

## Métodos HTTP

### 1. CreateSale
**Endpoint**: `POST /api/v1/sales`

**Descrição**: Registra uma nova venda de animal.

**Parâmetros**: Body com `CreateSaleRequest`

**Características**:
- Obtém `farm_id` do contexto (setado pelo middleware)
- Valida formato de data ("2006-01-02")
- Atualiza status do animal para "Vendido"

**Resposta**: Venda criada (201 Created).

---

### 2. GetSalesByFarm
**Endpoint**: `GET /api/v1/sales`

**Descrição**: Lista todas as vendas de uma fazenda.

**Características**:
- Obtém `farm_id` do contexto
- Inclui dados completos do animal vendido

**Resposta**: Lista de vendas.

---

### 3. GetSalesHistory
**Endpoint**: `GET /api/v1/sales/history`

**Descrição**: Retorna histórico completo de vendas da fazenda.

**Resposta**: Lista de vendas ordenadas por data.

---

### 4. GetSalesByDateRange
**Endpoint**: `GET /api/v1/sales/date-range?start_date={date}&end_date={date}`

**Descrição**: Lista vendas em um período específico.

**Parâmetros**:
- Query `start_date` (obrigatório, formato "2006-01-02")
- Query `end_date` (obrigatório, formato "2006-01-02")

**Resposta**: Lista de vendas no período.

---

### 5. GetSaleByID
**Endpoint**: `GET /api/v1/sales/{id}`

**Descrição**: Busca uma venda específica por ID.

**Parâmetros**: Path `id` (obrigatório)

**Resposta**: Dados completos da venda.

---

### 6. UpdateSale
**Endpoint**: `PUT /api/v1/sales/{id}`

**Descrição**: Atualiza dados de uma venda.

**Parâmetros**:
- Path `id` (obrigatório)
- Body com `UpdateSaleRequest`

**Características**:
- Mantém `animal_id` e `farm_id` originais
- Atualiza apenas campos permitidos

**Resposta**: Venda atualizada (200 OK).

---

### 7. DeleteSale
**Endpoint**: `DELETE /api/v1/sales/{id}`

**Descrição**: Remove uma venda do sistema.

**Parâmetros**: Path `id` (obrigatório)

**Resposta**: Confirmação de exclusão (200 OK).

---

### 8. GetSalesByAnimal
**Endpoint**: `GET /api/v1/animals/{animal_id}/sales`

**Descrição**: Lista todas as vendas de um animal específico.

**Parâmetros**: Path `animal_id` (obrigatório)

**Resposta**: Lista de vendas do animal.

---

### 9. GetMonthlySalesStats
**Endpoint**: `GET /api/v1/sales/monthly-stats`

**Descrição**: Retorna estatísticas de vendas do mês atual.

**Resposta**: Contagem de vendas no mês.

---

### 10. GetMonthlySalesAndPurchases
**Endpoint**: `GET /api/v1/sales/monthly-data?months={n}`

**Descrição**: Retorna dados mensais de vendas e compras.

**Parâmetros**: Query `months` (opcional, padrão: 12, máximo: 24)

**Resposta**: Dados mensais de vendas e compras (compras sempre 0.0).

---

### 11. GetOverviewStats
**Endpoint**: `GET /api/v1/sales/overview`

**Descrição**: Retorna estatísticas gerais da fazenda.

**Resposta**:
```json
{
  "males_count": 10,
  "females_count": 50,
  "total_sold": 5,
  "total_revenue": 15000.00
}
```

---

## Características Especiais

### Contexto de Farm ID
Todos os métodos obtêm o `farm_id` do contexto HTTP, que é setado pelo middleware de autenticação. Isso garante que:
- Usuários só vejam vendas de suas fazendas
- Novas vendas sejam associadas à fazenda correta

### Atualização de Status do Animal
Ao criar uma venda, o status do animal é automaticamente atualizado para "Vendido" no service.

### Formato de Data
Todas as datas devem estar no formato "2006-01-02" (YYYY-MM-DD).

