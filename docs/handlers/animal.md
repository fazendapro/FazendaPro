# Handler: Animal

## Visão Geral

O `AnimalHandler` gerencia todas as operações relacionadas a animais na fazenda. Este handler implementa o CRUD completo (Create, Read, Update, Delete) e funcionalidades adicionais como busca por sexo e upload de fotos.

## Estrutura

```go
type AnimalHandler struct {
    service *service.AnimalService
}
```

**Dependências**:
- `AnimalService`: Service que contém a lógica de negócio para animais

**Construtor**:
```go
func NewAnimalHandler(service *service.AnimalService) *AnimalHandler
```

## DTOs (Data Transfer Objects)

### AnimalData

Estrutura base que representa os dados de um animal:

```go
type AnimalData struct {
    ID                   uint   `json:"id"`
    FarmID               uint   `json:"farm_id"`
    EarTagNumberLocal     int    `json:"ear_tag_number_local"`
    EarTagNumberRegister  int    `json:"ear_tag_number_register"`
    AnimalName            string `json:"animal_name"`
    Sex                   int    `json:"sex"`  // 0 = Fêmea, 1 = Macho
    Breed                 string `json:"breed"`
    Type                  string `json:"type"`
    BirthDate             string `json:"birth_date,omitempty"`
    Photo                 string `json:"photo,omitempty"`
    FatherID              *uint  `json:"father_id,omitempty"`
    MotherID              *uint  `json:"mother_id,omitempty"`
    Confinement           bool   `json:"confinement"`
    AnimalType            int    `json:"animal_type"`
    Status                int    `json:"status"`
    Fertilization         bool   `json:"fertilization"`
    Castrated             bool   `json:"castrated"`
    Purpose               int    `json:"purpose"`  // 0 = Carne, 1 = Leite, 2 = Reprodução
    CurrentBatch          int    `json:"current_batch"`
}
```

### CreateAnimalRequest

Request para criação de animal:

```go
type CreateAnimalRequest struct {
    AnimalData
}
```

### AnimalResponse

Response com dados completos do animal, incluindo informações dos pais:

```go
type AnimalResponse struct {
    AnimalData
    Father    *AnimalParent `json:"father,omitempty"`
    Mother    *AnimalParent `json:"mother,omitempty"`
    CreatedAt string        `json:"createdAt"`
    UpdatedAt string        `json:"updatedAt"`
}
```

### AnimalParent

Informações resumidas do pai/mãe do animal:

```go
type AnimalParent struct {
    ID                uint   `json:"id"`
    AnimalName        string `json:"animal_name"`
    EarTagNumberLocal int    `json:"ear_tag_number_local"`
}
```

## Métodos HTTP

### 1. CreateAnimal

**Endpoint**: `POST /api/v1/animals`

**Método HTTP**: POST

**Autenticação**: Requerida (middleware Auth)

**Descrição**: Cria um novo animal na fazenda.

**Parâmetros**:
- Body (JSON): `CreateAnimalRequest`

**Validações**:
- Método HTTP deve ser POST
- JSON deve ser válido
- Validações de negócio são feitas no service

**Fluxo**:
1. Valida método HTTP
2. Decodifica JSON do body
3. Converte DTO para Model
4. Chama `service.CreateAnimal()`
5. Retorna resposta com ID do animal criado

**Resposta de Sucesso** (201 Created):
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

**Resposta de Erro**:
- `400 Bad Request`: JSON inválido ou erro de validação
- `405 Method Not Allowed`: Método HTTP incorreto

---

### 2. GetAnimal

**Endpoint**: `GET /api/v1/animals?id={animalID}`

**Método HTTP**: GET

**Autenticação**: Requerida

**Descrição**: Busca um animal específico por ID.

**Parâmetros**:
- Query: `id` (uint, obrigatório)

**Validações**:
- ID deve ser fornecido
- ID deve ser um número válido

**Fluxo**:
1. Extrai ID da query string
2. Converte ID para uint
3. Chama `service.GetAnimalByID()`
4. Verifica se animal existe
5. Converte model para response
6. Retorna dados do animal

**Resposta de Sucesso** (200 OK):
```json
{
  "success": true,
  "message": "Animal encontrado com sucesso",
  "data": {
    "id": 1,
    "farm_id": 1,
    "animal_name": "Branquinha",
    "ear_tag_number_local": 123,
    "sex": 0,
    "breed": "Holandesa",
    "type": "Bovino",
    "father": {
      "id": 2,
      "animal_name": "Touro",
      "ear_tag_number_local": 50
    },
    "mother": {
      "id": 3,
      "animal_name": "Vaca",
      "ear_tag_number_local": 100
    },
    "createdAt": "2024-01-15 10:30:00",
    "updatedAt": "2024-01-15 10:30:00"
  },
  "code": 200
}
```

**Resposta de Erro**:
- `400 Bad Request`: ID não fornecido ou inválido
- `404 Not Found`: Animal não encontrado
- `500 Internal Server Error`: Erro interno

---

### 3. GetAnimalsByFarm

**Endpoint**: `GET /api/v1/animals/farm?farmId={farmID}`

**Método HTTP**: GET

**Autenticação**: Requerida

**Descrição**: Lista todos os animais de uma fazenda.

**Parâmetros**:
- Query: `farmId` (uint, obrigatório)

**Validações**:
- farmId deve ser fornecido
- farmId deve ser um número válido

**Fluxo**:
1. Extrai farmId da query string
2. Converte farmId para uint
3. Chama `service.GetAnimalsByFarmID()`
4. Converte cada animal para response
5. Retorna lista de animais

**Resposta de Sucesso** (200 OK):
```json
{
  "success": true,
  "message": "Animais encontrados com sucesso (5 animais)",
  "data": [
    {
      "id": 1,
      "animal_name": "Branquinha",
      ...
    },
    {
      "id": 2,
      "animal_name": "Preta",
      ...
    }
  ],
  "code": 200
}
```

**Resposta de Erro**:
- `400 Bad Request`: farmId não fornecido ou inválido
- `500 Internal Server Error`: Erro interno

---

### 4. UpdateAnimal

**Endpoint**: `PUT /api/v1/animals`

**Método HTTP**: PUT

**Autenticação**: Requerida

**Descrição**: Atualiza os dados de um animal existente.

**Parâmetros**:
- Body (JSON): `CreateAnimalRequest` (deve incluir `id`)

**Validações**:
- Método HTTP deve ser PUT
- JSON deve ser válido
- ID deve estar presente no body
- Validações de negócio são feitas no service

**Fluxo**:
1. Valida método HTTP
2. Decodifica JSON do body
3. Converte DTO para Model
4. Chama `service.UpdateAnimal()`
5. Busca animal atualizado
6. Retorna dados atualizados

**Resposta de Sucesso** (200 OK):
```json
{
  "success": true,
  "message": "Animal atualizado com sucesso",
  "data": {
    "id": 1,
    "animal_name": "Branquinha Atualizada",
    ...
  },
  "code": 200
}
```

**Resposta de Erro**:
- `400 Bad Request`: JSON inválido ou erro de validação
- `405 Method Not Allowed`: Método HTTP incorreto

---

### 5. DeleteAnimal

**Endpoint**: `DELETE /api/v1/animals?id={animalID}`

**Método HTTP**: DELETE

**Autenticação**: Requerida

**Descrição**: Remove um animal do sistema.

**Parâmetros**:
- Query: `id` (uint, obrigatório)

**Validações**:
- Método HTTP deve ser DELETE
- ID deve ser fornecido
- ID deve ser um número válido

**Fluxo**:
1. Valida método HTTP
2. Extrai ID da query string
3. Converte ID para uint
4. Chama `service.DeleteAnimal()`
5. Retorna confirmação

**Resposta de Sucesso** (200 OK):
```json
{
  "success": true,
  "message": "Animal deletado com sucesso",
  "data": null,
  "code": 200
}
```

**Resposta de Erro**:
- `400 Bad Request`: ID não fornecido, inválido ou erro de validação
- `405 Method Not Allowed`: Método HTTP incorreto

---

### 6. GetAnimalsBySex

**Endpoint**: `GET /api/v1/animals/sex?farmId={farmID}&sex={sex}`

**Método HTTP**: GET

**Autenticação**: Requerida

**Descrição**: Lista animais de uma fazenda filtrados por sexo.

**Parâmetros**:
- Query: `farmId` (uint, obrigatório)
- Query: `sex` (int, obrigatório) - 0 = Fêmea, 1 = Macho

**Validações**:
- farmId deve ser fornecido
- sex deve ser fornecido
- Ambos devem ser números válidos

**Fluxo**:
1. Extrai farmId e sex da query string
2. Converte ambos para tipos apropriados
3. Chama `service.GetAnimalsByFarmIDAndSex()`
4. Converte cada animal para response
5. Retorna lista filtrada

**Resposta de Sucesso** (200 OK):
```json
{
  "success": true,
  "message": "Animais encontrados com sucesso (3 animais)",
  "data": [
    {
      "id": 1,
      "animal_name": "Branquinha",
      "sex": 0,
      ...
    }
  ],
  "code": 200
}
```

**Resposta de Erro**:
- `400 Bad Request`: Parâmetros não fornecidos ou inválidos
- `500 Internal Server Error`: Erro interno

---

### 7. UploadAnimalPhoto

**Endpoint**: `POST /api/v1/animals/photo`

**Método HTTP**: POST

**Autenticação**: Requerida

**Descrição**: Faz upload de uma foto para um animal, convertendo para base64.

**Parâmetros**:
- Form Data:
  - `animal_id` (string, obrigatório)
  - `photo` (file, obrigatório) - máximo 10MB

**Validações**:
- Método HTTP deve ser POST
- animal_id deve ser fornecido
- Arquivo photo deve ser fornecido
- Tamanho máximo: 10MB

**Fluxo**:
1. Valida método HTTP
2. Faz parse do multipart form
3. Extrai animal_id e arquivo
4. Lê arquivo e converte para base64
5. Busca animal existente
6. Atualiza campo photo com base64
7. Chama `service.UpdateAnimal()`
8. Busca animal atualizado
9. Retorna animal com foto atualizada

**Formato da Foto**:
A foto é armazenada como string base64 no formato:
```
data:image/jpeg;base64,{base64_string}
```

**Resposta de Sucesso** (200 OK):
```json
{
  "success": true,
  "message": "Foto do animal atualizada com sucesso",
  "data": {
    "id": 1,
    "animal_name": "Branquinha",
    "photo": "data:image/jpeg;base64,/9j/4AAQSkZJRg...",
    ...
  },
  "code": 200
}
```

**Resposta de Erro**:
- `400 Bad Request`: Parâmetros inválidos ou arquivo não fornecido
- `404 Not Found`: Animal não encontrado
- `405 Method Not Allowed`: Método HTTP incorreto
- `500 Internal Server Error`: Erro ao processar arquivo

---

## Funções Auxiliares

### animalDataToModel

Converte `AnimalData` (DTO) para `models.Animal` (Model):

```go
func animalDataToModel(data AnimalData) models.Animal
```

**Funcionalidades**:
- Converte string de data para `*time.Time`
- Mapeia todos os campos do DTO para o model

### modelToAnimalResponse

Converte `models.Animal` (Model) para `AnimalResponse` (DTO):

```go
func modelToAnimalResponse(animal *models.Animal) AnimalResponse
```

**Funcionalidades**:
- Formata datas para strings
- Inclui informações dos pais (Father e Mother) se disponíveis
- Formata timestamps (CreatedAt, UpdatedAt)

---

## Exemplos de Requisições

### Criar Animal

```http
POST /api/v1/animals
Authorization: Bearer {token}
Content-Type: application/json

{
  "farm_id": 1,
  "ear_tag_number_local": 123,
  "animal_name": "Branquinha",
  "sex": 0,
  "breed": "Holandesa",
  "type": "Bovino",
  "animal_type": 1,
  "purpose": 1,
  "birth_date": "2020-05-15",
  "father_id": 2,
  "mother_id": 3
}
```

### Buscar Animal por ID

```http
GET /api/v1/animals?id=1
Authorization: Bearer {token}
```

### Listar Animais da Fazenda

```http
GET /api/v1/animals/farm?farmId=1
Authorization: Bearer {token}
```

### Atualizar Animal

```http
PUT /api/v1/animals
Authorization: Bearer {token}
Content-Type: application/json

{
  "id": 1,
  "animal_name": "Branquinha Atualizada",
  "breed": "Holandesa Preta e Branca",
  ...
}
```

### Deletar Animal

```http
DELETE /api/v1/animals?id=1
Authorization: Bearer {token}
```

### Buscar por Sexo

```http
GET /api/v1/animals/sex?farmId=1&sex=0
Authorization: Bearer {token}
```

### Upload de Foto

```http
POST /api/v1/animals/photo
Authorization: Bearer {token}
Content-Type: multipart/form-data

animal_id: 1
photo: [arquivo de imagem]
```

---

## Notas de Implementação

1. **Validação de Método HTTP**: Alguns métodos validam explicitamente o método HTTP antes de processar
2. **Conversão de Tipos**: Conversões de string para uint/int são feitas manualmente com validação
3. **Formatação de Datas**: Datas são formatadas no padrão "2006-01-02" para JSON
4. **Base64 para Fotos**: Fotos são convertidas para base64 e armazenadas como string no banco
5. **Relacionamentos**: O handler carrega automaticamente informações dos pais (Father/Mother) quando disponíveis
6. **Logs de Debug**: Alguns métodos incluem logs de debug (fmt.Printf) para facilitar troubleshooting

