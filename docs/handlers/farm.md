# Handler: Farm

## Visão Geral

O `FarmHandler` gerencia operações relacionadas a fazendas, incluindo busca e atualização de dados.

## Estrutura

```go
type FarmHandler struct {
    service *service.FarmService
}
```

## DTOs

### UpdateFarmRequest
```go
type UpdateFarmRequest struct {
    Logo string `json:"logo"`
}
```

## Métodos HTTP

### 1. GetFarm
**Endpoint**: `GET /api/v1/farm?id={farmID}`

**Descrição**: Busca uma fazenda por ID, incluindo dados da empresa.

**Parâmetros**: Query `id` (obrigatório)

**Fluxo**:
1. Extrai farmID da query
2. Busca fazenda via service
3. Carrega dados da empresa associada
4. Retorna fazenda completa

**Resposta**: Dados da fazenda com informações da empresa.

---

### 2. UpdateFarm
**Endpoint**: `PUT /api/v1/farm?id={farmID}`

**Descrição**: Atualiza dados de uma fazenda (principalmente logo).

**Parâmetros**: 
- Query `id` (obrigatório)
- Body com `UpdateFarmRequest`

**Resposta**: Fazenda atualizada (200 OK).

