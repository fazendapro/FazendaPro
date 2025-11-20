# Handler: Farm Selection

## Visão Geral

O `FarmSelectionHandler` gerencia a seleção de fazendas por usuários. Permite que usuários vejam suas fazendas e selecionem uma para trabalhar.

## Estrutura

```go
type FarmSelectionHandler struct {
    service   *service.UserService
    jwtSecret string
}
```

## DTOs

### GetUserFarmsResponse
```go
type GetUserFarmsResponse struct {
    Success        bool        `json:"success"`
    Message        string      `json:"message"`
    Farms          interface{} `json:"farms"`
    AutoSelect     bool        `json:"auto_select"`
    SelectedFarmID *uint       `json:"selected_farm_id,omitempty"`
}
```

### SelectFarmRequest
```go
type SelectFarmRequest struct {
    FarmID uint `json:"farm_id"`
}
```

## Métodos HTTP

### 1. GetUserFarms
**Endpoint**: `GET /api/v1/farms/user`

**Descrição**: Lista todas as fazendas de um usuário e indica se deve fazer auto-seleção.

**Fluxo**:
1. Extrai userID do token JWT
2. Busca fazendas do usuário via service
3. Verifica se deve fazer auto-seleção (usuário tem apenas 1 fazenda)
4. Se auto-select, define selected_farm_id
5. Retorna lista de fazendas

**Resposta**: Lista de fazendas com flag de auto-seleção.

---

### 2. SelectFarm
**Endpoint**: `POST /api/v1/farms/select`

**Descrição**: Seleciona uma fazenda para o usuário trabalhar.

**Parâmetros**: Body com `SelectFarmRequest`

**Fluxo**:
1. Extrai userID do token JWT
2. Valida que a fazenda pertence ao usuário
3. Retorna confirmação

**Resposta**: Confirmação de seleção (200 OK).

---

## Função Auxiliar

### extractUserIDFromToken
Extrai o ID do usuário do token JWT no header Authorization.

**Nota**: Implementa parsing manual do JWT para extrair o claim "sub".

