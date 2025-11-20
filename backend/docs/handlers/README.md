# Documentação dos Handlers

Esta pasta contém a documentação detalhada de todos os handlers da API FazendaPro.

## Índice

### Handlers de Entidades Principais

1. **[Animal Handler](animal.md)** - Gerencia operações de animais
   - 7 métodos HTTP
   - CRUD completo
   - Upload de fotos
   - Busca por sexo

2. **[Milk Collection Handler](milk_collection.md)** - Gerencia coletas de leite
   - 5 métodos HTTP
   - Criação e atualização de coletas
   - Estatísticas de produção
   - Top produtoras

3. **[Reproduction Handler](reproduction.md)** - Gerencia reprodução
   - 9 métodos HTTP
   - Gestão de fases de reprodução
   - Próximas vacas a parir
   - Histórico completo

4. **[Sale Handler](sale.md)** - Gerencia vendas
   - 12 métodos HTTP
   - CRUD completo
   - Estatísticas e relatórios
   - Histórico de vendas

5. **[Debt Handler](debt.md)** - Gerencia dívidas
   - 4 métodos HTTP
   - Paginação
   - Totais por pessoa

### Handlers de Autenticação e Usuários

6. **[Auth Handler](auth.md)** - Autenticação e autorização
   - Login e registro
   - Renovação de tokens (JWT)
   - Logout
   - Gerenciamento de sessão

7. **[User Handler](user.md)** - Gerenciamento de usuários
   - 4 métodos HTTP
   - Criação e busca de usuários
   - Atualização de dados pessoais

### Handlers de Configuração

8. **[Farm Handler](farm.md)** - Gerenciamento de fazendas
   - 2 métodos HTTP
   - Busca e atualização de fazendas
   - Dados da empresa

9. **[Farm Selection Handler](farm_selection.md)** - Seleção de fazendas
   - 2 métodos HTTP
   - Lista fazendas do usuário
   - Seleção de fazenda ativa

### Utilitários

10. **[Error Response](error_response.md)** - Funções utilitárias
    - Padronização de respostas
    - SendSuccessResponse
    - SendErrorResponse

## Estrutura Comum dos Handlers

Todos os handlers seguem um padrão similar:

1. **Struct do Handler**: Contém dependências (services, repositories)
2. **Construtor**: Função `New*Handler()` para inicialização
3. **DTOs**: Estruturas para Request e Response
4. **Métodos HTTP**: Implementam `http.HandlerFunc`

## Padrão de Resposta

Todos os handlers usam as funções utilitárias de `error_response.go`:

- **Sucesso**: `SendSuccessResponse(w, data, message, statusCode)`
- **Erro**: `SendErrorResponse(w, message, statusCode)`

## Autenticação

A maioria dos handlers requer autenticação via middleware:

```go
r.Use(middleware.Auth(cfg.JWTSecret))
```

Handlers públicos (sem autenticação):
- `AuthHandler`: Login, Register, RefreshToken, Logout

## Convenções

- **Endpoints**: Seguem padrão RESTful
- **Métodos HTTP**: GET (buscar), POST (criar), PUT (atualizar), DELETE (remover)
- **Formato de Data**: "2006-01-02" (YYYY-MM-DD)
- **IDs**: Sempre uint, extraídos de query params ou path params
- **Validação**: Feita no handler (parâmetros) e no service (regras de negócio)

## Como Usar Esta Documentação

1. **Para Desenvolvedores**: Use para entender como cada endpoint funciona
2. **Para Integração**: Use os exemplos de requisições e respostas
3. **Para Manutenção**: Entenda o fluxo de cada método antes de modificar

## Próximos Passos

- Adicionar exemplos de curl para cada endpoint
- Documentar códigos de erro específicos
- Adicionar diagramas de fluxo para métodos complexos

