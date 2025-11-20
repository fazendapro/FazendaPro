# Handler: Error Response

## Visão Geral

O arquivo `error_response.go` contém funções utilitárias para padronizar as respostas HTTP da API. Estas funções são usadas por todos os handlers para garantir consistência nas respostas de sucesso e erro.

## Estrutura

Este arquivo não possui uma struct handler, apenas funções auxiliares que podem ser chamadas por qualquer handler.

## Funções Utilitárias

### 1. `SendErrorResponse`

**Assinatura**:
```go
func SendErrorResponse(w http.ResponseWriter, message string, statusCode int)
```

**Propósito**: Envia uma resposta de erro padronizada em formato JSON.

**Parâmetros**:
- `w`: ResponseWriter HTTP
- `message`: Mensagem de erro descritiva
- `statusCode`: Código de status HTTP (ex: 400, 404, 500)

**Estrutura da Resposta**:
```json
{
  "success": false,
  "error": "Bad Request",
  "message": "Mensagem de erro descritiva",
  "code": 400
}
```

**Exemplo de Uso**:
```go
SendErrorResponse(w, "ID do animal é obrigatório", http.StatusBadRequest)
```

**Resposta HTTP**:
- Status Code: O valor passado em `statusCode`
- Content-Type: `application/json`
- Body: JSON com estrutura de erro

---

### 2. `SendSuccessResponse`

**Assinatura**:
```go
func SendSuccessResponse(w http.ResponseWriter, data interface{}, message string, statusCode int)
```

**Propósito**: Envia uma resposta de sucesso padronizada em formato JSON.

**Parâmetros**:
- `w`: ResponseWriter HTTP
- `data`: Dados a serem retornados (pode ser qualquer tipo: objeto, array, null)
- `message`: Mensagem de sucesso descritiva
- `statusCode`: Código de status HTTP (geralmente 200 ou 201)

**Estrutura da Resposta**:
```json
{
  "success": true,
  "message": "Operação realizada com sucesso",
  "data": { ... },
  "code": 200
}
```

**Exemplo de Uso**:
```go
data := map[string]interface{}{
    "id": animal.ID,
}
SendSuccessResponse(w, data, "Animal criado com sucesso", http.StatusCreated)
```

**Resposta HTTP**:
- Status Code: O valor passado em `statusCode` (geralmente 200 ou 201)
- Content-Type: `application/json`
- Body: JSON com estrutura de sucesso

---

## Estrutura de Resposta Padronizada

### Resposta de Erro

```json
{
  "success": false,
  "error": "Status Text do HTTP",
  "message": "Descrição detalhada do erro",
  "code": 400
}
```

**Campos**:
- `success`: Sempre `false` para erros
- `error`: Texto do status HTTP (ex: "Bad Request", "Not Found", "Internal Server Error")
- `message`: Mensagem específica do erro
- `code`: Código numérico do status HTTP

### Resposta de Sucesso

```json
{
  "success": true,
  "message": "Mensagem de sucesso",
  "data": { ... },
  "code": 200
}
```

**Campos**:
- `success`: Sempre `true` para sucessos
- `message`: Mensagem descritiva da operação
- `data`: Dados retornados (pode ser objeto, array ou null)
- `code`: Código numérico do status HTTP

---

## Exemplos de Uso nos Handlers

### Exemplo 1: Erro de Validação

```go
if animalID == "" {
    SendErrorResponse(w, "ID do animal é obrigatório", http.StatusBadRequest)
    return
}
```

**Resposta**:
```json
{
  "success": false,
  "error": "Bad Request",
  "message": "ID do animal é obrigatório",
  "code": 400
}
```

### Exemplo 2: Recurso Não Encontrado

```go
if animal == nil {
    SendErrorResponse(w, "Animal não encontrado", http.StatusNotFound)
    return
}
```

**Resposta**:
```json
{
  "success": false,
  "error": "Not Found",
  "message": "Animal não encontrado",
  "code": 404
}
```

### Exemplo 3: Sucesso com Dados

```go
response := modelToAnimalResponse(animal)
SendSuccessResponse(w, response, "Animal encontrado com sucesso", http.StatusOK)
```

**Resposta**:
```json
{
  "success": true,
  "message": "Animal encontrado com sucesso",
  "data": {
    "id": 1,
    "animal_name": "Branquinha",
    "farm_id": 1,
    ...
  },
  "code": 200
}
```

### Exemplo 4: Sucesso sem Dados (Criação)

```go
data := map[string]interface{}{
    "id": animal.ID,
}
SendSuccessResponse(w, data, "Animal criado com sucesso", http.StatusCreated)
```

**Resposta**:
```json
{
  "success": true,
  "message": "Animal criado com sucesso",
  "data": {
    "id": 1
  },
  "code": 201
}
```

---

## Vantagens da Padronização

1. **Consistência**: Todas as respostas seguem o mesmo formato
2. **Facilita Frontend**: O frontend sabe sempre onde encontrar os dados e mensagens
3. **Manutenibilidade**: Mudanças no formato de resposta são centralizadas
4. **Debugging**: Mais fácil identificar problemas com formato consistente

---

## Códigos de Status HTTP Comuns

### Sucesso
- `200 OK`: Operação bem-sucedida (GET, PUT, DELETE)
- `201 Created`: Recurso criado com sucesso (POST)

### Erro do Cliente
- `400 Bad Request`: Requisição inválida (validação, JSON malformado)
- `401 Unauthorized`: Não autenticado
- `403 Forbidden`: Não autorizado
- `404 Not Found`: Recurso não encontrado
- `405 Method Not Allowed`: Método HTTP não permitido

### Erro do Servidor
- `500 Internal Server Error`: Erro interno do servidor

---

## Notas de Implementação

- As funções sempre definem o header `Content-Type: application/json`
- O status code é definido antes de escrever o body
- A codificação JSON é feita automaticamente pelo `json.NewEncoder`
- Estas funções são chamadas por praticamente todos os handlers do projeto

