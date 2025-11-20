# Handler: Debt

## Visão Geral

O `DebtHandler` gerencia operações relacionadas a dívidas, incluindo criação, listagem com paginação, exclusão e cálculo de totais por pessoa.

## Estrutura

```go
type DebtHandler struct {
    service *service.DebtService
}
```

## DTOs

### CreateDebtRequest
```go
type CreateDebtRequest struct {
    Person string  `json:"person"`
    Value  float64 `json:"value"`
}
```

### DebtResponse
```go
type DebtResponse struct {
    ID        uint    `json:"id"`
    Person    string  `json:"person"`
    Value     float64 `json:"value"`
    CreatedAt string  `json:"created_at"`
    UpdatedAt string  `json:"updated_at"`
}
```

### DebtListResponse
```go
type DebtListResponse struct {
    Debts []DebtResponse `json:"debts"`
    Total int64          `json:"total"`
    Page  int            `json:"page"`
    Limit int            `json:"limit"`
}
```

## Métodos HTTP

### 1. CreateDebt
**Endpoint**: `POST /debts`

**Descrição**: Cria uma nova dívida.

**Parâmetros**: Body com `CreateDebtRequest`

**Resposta**: Dívida criada (201 Created).

---

### 2. GetDebts
**Endpoint**: `GET /debts?page={page}&limit={limit}&year={year}&month={month}`

**Descrição**: Lista dívidas com paginação e filtros opcionais.

**Parâmetros**:
- `page` (opcional, padrão: 1)
- `limit` (opcional, padrão: 10)
- `year` (opcional)
- `month` (opcional, 1-12)

**Resposta**: Lista paginada de dívidas com total.

---

### 3. DeleteDebt
**Endpoint**: `DELETE /debts/{id}`

**Descrição**: Remove uma dívida.

**Parâmetros**: Path `id` (obrigatório)

**Resposta**: Confirmação de exclusão (200 OK).

---

### 4. GetTotalByPerson
**Endpoint**: `GET /debts/total-by-person?year={year}&month={month}`

**Descrição**: Calcula total de dívidas por pessoa em um mês específico.

**Parâmetros**:
- `year` (obrigatório)
- `month` (obrigatório, 1-12)

**Resposta**: Lista de totais por pessoa.

