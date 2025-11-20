# Handler: Milk Collection

## Visão Geral

O `MilkCollectionHandler` gerencia operações relacionadas à coleta de leite, incluindo criação, atualização, listagem e estatísticas de produção.

## Estrutura

```go
type MilkCollectionHandler struct {
    service *service.MilkCollectionService
}
```

## DTOs

### CreateMilkCollectionRequest
```go
type CreateMilkCollectionRequest struct {
    AnimalID uint    `json:"animal_id" validate:"required"`
    Liters   float64 `json:"liters" validate:"required,min=0"`
    Date     string  `json:"date" validate:"required"`
}
```

### MilkCollectionData
```go
type MilkCollectionData struct {
    ID        uint       `json:"id"`
    AnimalID  uint       `json:"animal_id"`
    Animal    AnimalData `json:"animal"`
    Liters    float64    `json:"liters"`
    Date      time.Time  `json:"date"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}
```

## Métodos HTTP

### 1. CreateMilkCollection
**Endpoint**: `POST /api/v1/milk-collections`

**Descrição**: Registra uma nova coleta de leite.

**Parâmetros**: Body com `CreateMilkCollectionRequest`

**Validações**: Data no formato "2006-01-02"

**Resposta**: Coleta criada (201 Created).

---

### 2. UpdateMilkCollection
**Endpoint**: `PUT /api/v1/milk-collections/{id}`

**Descrição**: Atualiza uma coleta de leite existente.

**Parâmetros**: 
- Path `id` (obrigatório)
- Body com `CreateMilkCollectionRequest`

**Resposta**: Coleta atualizada (200 OK).

---

### 3. GetMilkCollectionsByFarmID
**Endpoint**: `GET /api/v1/milk-collections/farm/{farmId}?start_date={date}&end_date={date}`

**Descrição**: Lista coletas de leite de uma fazenda, opcionalmente filtradas por período.

**Parâmetros**:
- Path `farmId` (obrigatório)
- Query `start_date` (opcional, formato "2006-01-02")
- Query `end_date` (opcional, formato "2006-01-02")

**Resposta**: Lista de coletas com dados do animal.

---

### 4. GetMilkCollectionsByAnimalID
**Endpoint**: `GET /api/v1/milk-collections/animal/{animalId}`

**Descrição**: Lista todas as coletas de leite de um animal específico.

**Parâmetros**: Path `animalId` (obrigatório)

**Resposta**: Lista de coletas do animal.

---

### 5. GetTopMilkProducers
**Endpoint**: `GET /api/v1/milk-collections/top-producers?farmId={id}&limit={limit}&periodDays={days}`

**Descrição**: Retorna as maiores produtoras de leite de uma fazenda.

**Parâmetros**:
- Query `farmId` (obrigatório)
- Query `limit` (opcional, padrão: 10)
- Query `periodDays` (opcional, padrão: 30)

**Funcionalidades**:
- Calcula produção total por animal
- Calcula média diária de produção
- Ordena por produção total (decrescente)
- Retorna top N produtoras

**Resposta**: Lista de animais com estatísticas de produção.

