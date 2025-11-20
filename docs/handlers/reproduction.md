# Handler: Reproduction

## Visão Geral

O `ReproductionHandler` gerencia operações relacionadas à reprodução de animais, incluindo criação de registros, atualização de fases, busca e estatísticas.

## Estrutura

```go
type ReproductionHandler struct {
    service *service.ReproductionService
}
```

## DTOs

### ReproductionData
```go
type ReproductionData struct {
    ID                     uint    `json:"id"`
    AnimalID               uint    `json:"animal_id"`
    CurrentPhase           int     `json:"current_phase"`
    InseminationDate       *string `json:"insemination_date,omitempty"`
    PregnancyDate          *string `json:"pregnancy_date,omitempty"`
    ExpectedBirthDate      *string `json:"expected_birth_date,omitempty"`
    ActualBirthDate        *string `json:"actual_birth_date,omitempty"`
    LactationStartDate     *string `json:"lactation_start_date,omitempty"`
    LactationEndDate       *string `json:"lactation_end_date,omitempty"`
    DryPeriodStartDate     *string `json:"dry_period_start_date,omitempty"`
    VeterinaryConfirmation bool    `json:"veterinary_confirmation"`
    Observations           string  `json:"observations,omitempty"`
}
```

### UpdateReproductionPhaseRequest
```go
type UpdateReproductionPhaseRequest struct {
    AnimalID       uint                   `json:"animal_id"`
    NewPhase       int                    `json:"new_phase"`
    AdditionalData map[string]interface{} `json:"additional_data,omitempty"`
}
```

## Métodos HTTP

### 1. CreateReproduction
**Endpoint**: `POST /api/v1/reproductions`

**Descrição**: Cria um novo registro de reprodução.

**Parâmetros**: Body com `ReproductionData`

**Resposta**: ID do registro criado (201 Created).

---

### 2. GetReproduction
**Endpoint**: `GET /api/v1/reproductions?id={id}`

**Descrição**: Busca um registro de reprodução por ID.

**Parâmetros**: Query `id` (obrigatório)

**Resposta**: Dados completos do registro.

---

### 3. GetReproductionByAnimal
**Endpoint**: `GET /api/v1/reproductions/animal?animalId={id}`

**Descrição**: Busca registro de reprodução de um animal específico.

**Parâmetros**: Query `animalId` (obrigatório)

**Resposta**: Registro de reprodução do animal.

---

### 4. GetReproductionsByFarm
**Endpoint**: `GET /api/v1/reproductions/farm?farmId={id}`

**Descrição**: Lista todos os registros de reprodução de uma fazenda.

**Parâmetros**: Query `farmId` (obrigatório)

**Resposta**: Lista de registros de reprodução.

---

### 5. GetReproductionsByPhase
**Endpoint**: `GET /api/v1/reproductions/phase?phase={phase}`

**Descrição**: Lista registros de reprodução por fase.

**Parâmetros**: Query `phase` (obrigatório, número da fase)

**Resposta**: Lista de registros na fase especificada.

---

### 6. UpdateReproduction
**Endpoint**: `PUT /api/v1/reproductions`

**Descrição**: Atualiza um registro de reprodução completo.

**Parâmetros**: Body com `ReproductionData` (deve incluir ID)

**Resposta**: Confirmação (200 OK).

---

### 7. UpdateReproductionPhase
**Endpoint**: `PUT /api/v1/reproductions/phase`

**Descrição**: Atualiza apenas a fase de reprodução de um animal.

**Parâmetros**: Body com `UpdateReproductionPhaseRequest`

**Resposta**: Confirmação (200 OK).

---

### 8. DeleteReproduction
**Endpoint**: `DELETE /api/v1/reproductions?id={id}`

**Descrição**: Remove um registro de reprodução.

**Parâmetros**: Query `id` (obrigatório)

**Resposta**: Confirmação (200 OK).

---

### 9. GetNextToCalve
**Endpoint**: `GET /api/v1/reproductions/next-to-calve?farmId={id}`

**Descrição**: Lista animais que estão próximos a parir, ordenados por data esperada.

**Parâmetros**: Query `farmId` (obrigatório)

**Funcionalidades**:
- Filtra apenas animais na fase "Prenhas"
- Calcula dias até o parto (283 dias após data de prenhez)
- Classifica por prioridade (Alto: ≤30 dias, Médio: ≤60 dias, Baixo: >60 dias)
- Ordena por dias até parto (crescente)

**Resposta**: Lista de animais com informações de parto esperado.

