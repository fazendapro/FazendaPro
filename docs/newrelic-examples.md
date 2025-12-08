# Exemplos de Uso do New Relic

Este documento contém exemplos práticos de como usar o New Relic no FazendaPro Backend.

## 1. Instrumentação Automática (Já Implementada)

Todas as rotas HTTP são automaticamente instrumentadas pelo middleware customizado em `internal/routes/routes.go`.

**Exemplo de transação automática:**
```
GET /api/v1/animals/farm
POST /api/v1/sales
PUT /api/v1/reproductions
```

Cada uma dessas requisições gera automaticamente uma transação no New Relic com:
- Nome: `{METHOD} {PATH}`
- Tempo de resposta
- Status code
- Headers

## 2. Acessar Transação no Handler

Você pode acessar a transação atual em qualquer handler:

```go
import "github.com/newrelic/go-agent/v3/newrelic"

func (h *AnimalHandler) GetAnimalsByFarm(w http.ResponseWriter, r *http.Request) {
    // Obter transação do contexto
    txn := newrelic.FromContext(r.Context())
    
    // Adicionar atributos customizados
    farmID := r.URL.Query().Get("farmId")
    txn.AddAttribute("farm_id", farmID)
    
    // Seu código aqui...
}
```

## 3. Segmentos em Services

Para operações importantes em services, você pode criar segmentos customizados:

```go
import "github.com/newrelic/go-agent/v3/newrelic"

func (s *saleService) CreateSale(ctx context.Context, sale *models.Sale) error {
    txn := newrelic.FromContext(ctx)
    
    // Segmento para validação
    validationSegment := txn.StartSegment("ValidateSaleData")
    if err := s.validateSale(sale); err != nil {
        validationSegment.End()
        return err
    }
    validationSegment.End()
    
    // Segmento para operação de banco de dados
    dbSegment := txn.StartSegment("CreateSaleInDB")
    err := s.saleRepo.Create(ctx, sale)
    dbSegment.End()
    
    if err != nil {
        txn.NoticeError(err)
        return err
    }
    
    // Segmento para invalidação de cache
    cacheSegment := txn.StartSegment("InvalidateCache")
    s.cache.Delete(fmt.Sprintf("sales:farm:%d", sale.FarmID))
    cacheSegment.End()
    
    return nil
}
```

## 4. Reportar Erros Customizados

```go
func (s *animalService) UpdateAnimal(ctx context.Context, animal *models.Animal) error {
    txn := newrelic.FromContext(ctx)
    
    // Validação de negócio
    if animal.Status == models.AnimalStatusSold {
        err := errors.New("cannot update sold animal")
        txn.NoticeError(err)
        txn.AddAttribute("animal_id", animal.ID)
        txn.AddAttribute("error_type", "business_validation")
        return err
    }
    
    // Operação de atualização
    if err := s.animalRepo.Update(ctx, animal); err != nil {
        txn.NoticeError(err)
        txn.AddAttribute("operation", "update_animal")
        return err
    }
    
    return nil
}
```

## 5. Transações Assíncronas

Para operações que rodam em background:

```go
func (s *reportService) GenerateMonthlyReport(farmID uint) {
    // Criar nova transação para operação assíncrona
    txn := s.nrApp.StartTransaction("GenerateMonthlyReport")
    defer txn.End()
    
    txn.AddAttribute("farm_id", farmID)
    txn.AddAttribute("report_type", "monthly")
    
    // Criar contexto com a transação
    ctx := newrelic.NewContext(context.Background(), txn)
    
    // Sua lógica aqui
    if err := s.generateReport(ctx, farmID); err != nil {
        txn.NoticeError(err)
        log.Printf("Error generating report: %v", err)
    }
}
```

## 6. Segmentos de Database Queries

Para queries específicas que você quer monitorar:

```go
func (r *animalRepository) GetTopMilkProducers(ctx context.Context, farmID uint, limit int) ([]*models.Animal, error) {
    txn := newrelic.FromContext(ctx)
    
    // Criar segmento de datastore
    segment := &newrelic.DatastoreSegment{
        StartTime:  txn.StartSegmentNow(),
        Product:    newrelic.DatastorePostgres,
        Collection: "animals",
        Operation:  "SELECT",
        ParameterizedQuery: `
            SELECT a.* FROM animals a
            JOIN milk_collections mc ON a.id = mc.animal_id
            WHERE a.farm_id = ?
            GROUP BY a.id
            ORDER BY AVG(mc.quantity) DESC
            LIMIT ?
        `,
    }
    defer segment.End()
    
    var animals []*models.Animal
    err := r.db.Where("farm_id = ?", farmID).
        Joins("JOIN milk_collections ON animals.id = milk_collections.animal_id").
        Group("animals.id").
        Order("AVG(milk_collections.quantity) DESC").
        Limit(limit).
        Find(&animals).Error
    
    return animals, err
}
```

## 7. Segmentos de External Services

Para chamadas a APIs externas:

```go
func (s *weatherService) GetWeatherData(ctx context.Context, location string) (*WeatherData, error) {
    txn := newrelic.FromContext(ctx)
    
    // Criar segmento de external service
    segment := newrelic.StartExternalSegment(txn, &http.Request{
        Method: "GET",
        URL:    &url.URL{
            Scheme: "https",
            Host:   "api.weather.com",
            Path:   "/v1/current",
        },
    })
    defer segment.End()
    
    // Fazer requisição HTTP
    resp, err := http.Get(fmt.Sprintf("https://api.weather.com/v1/current?location=%s", location))
    if err != nil {
        txn.NoticeError(err)
        return nil, err
    }
    defer resp.Body.Close()
    
    // Processar resposta...
    var data WeatherData
    json.NewDecoder(resp.Body).Decode(&data)
    
    return &data, nil
}
```

## 8. Atributos Customizados para Análise

Adicione atributos que ajudem na análise:

```go
func (h *SaleHandler) CreateSale(w http.ResponseWriter, r *http.Request) {
    txn := newrelic.FromContext(r.Context())
    
    var sale models.Sale
    json.NewDecoder(r.Body).Decode(&sale)
    
    // Adicionar atributos de negócio
    txn.AddAttribute("sale_type", sale.Type)
    txn.AddAttribute("sale_value", sale.Value)
    txn.AddAttribute("farm_id", sale.FarmID)
    txn.AddAttribute("has_buyer_cpf", sale.BuyerCPF != "")
    
    // Categorizar por valor
    if sale.Value > 10000 {
        txn.AddAttribute("sale_category", "high_value")
    } else if sale.Value > 5000 {
        txn.AddAttribute("sale_category", "medium_value")
    } else {
        txn.AddAttribute("sale_category", "low_value")
    }
    
    // Processar venda...
}
```

## 9. Monitorar Performance de Cache

```go
func (s *animalService) GetAnimalsByFarm(ctx context.Context, farmID uint) ([]*models.Animal, error) {
    txn := newrelic.FromContext(ctx)
    
    cacheKey := fmt.Sprintf("animals:farm:%d", farmID)
    
    // Tentar buscar do cache
    cacheSegment := txn.StartSegment("CacheLookup")
    cachedData, err := s.cache.Get(cacheKey)
    cacheSegment.End()
    
    if err == nil && cachedData != nil {
        txn.AddAttribute("cache_hit", true)
        var animals []*models.Animal
        json.Unmarshal([]byte(cachedData), &animals)
        return animals, nil
    }
    
    txn.AddAttribute("cache_hit", false)
    
    // Buscar do banco
    dbSegment := txn.StartSegment("DatabaseQuery")
    animals, err := s.animalRepo.GetByFarmID(ctx, farmID)
    dbSegment.End()
    
    if err != nil {
        txn.NoticeError(err)
        return nil, err
    }
    
    // Salvar no cache
    cacheSetSegment := txn.StartSegment("CacheSet")
    data, _ := json.Marshal(animals)
    s.cache.Set(cacheKey, string(data), 300)
    cacheSetSegment.End()
    
    return animals, nil
}
```

## 10. Logs Contextualizados

Adicionar informações da transação aos logs:

```go
func (s *reproductionService) UpdateReproductionPhase(ctx context.Context, id uint, phase string) error {
    txn := newrelic.FromContext(ctx)
    
    log.Printf("[Transaction: %s] Updating reproduction %d to phase %s", 
        txn.GetTraceMetadata().TraceID, id, phase)
    
    // Sua lógica aqui...
    
    if err != nil {
        log.Printf("[Transaction: %s] Error updating reproduction: %v", 
            txn.GetTraceMetadata().TraceID, err)
        txn.NoticeError(err)
        return err
    }
    
    return nil
}
```

## Queries Úteis no New Relic

### 1. Taxa de Erro por Endpoint
```sql
SELECT count(*) 
FROM Transaction 
WHERE appName = 'fazendapro-backend' 
FACET request.uri 
WHERE error IS true
```

### 2. Endpoints Mais Lentos
```sql
SELECT average(duration) 
FROM Transaction 
WHERE appName = 'fazendapro-backend' 
FACET name 
ORDER BY average(duration) DESC
```

### 3. Análise de Vendas por Categoria
```sql
SELECT count(*) 
FROM Transaction 
WHERE appName = 'fazendapro-backend' 
AND name = 'POST /api/v1/sales'
FACET sale_category
```

### 4. Performance de Cache
```sql
SELECT percentage(count(*), WHERE cache_hit IS true) 
FROM Transaction 
WHERE appName = 'fazendapro-backend'
```

### 5. Erros de Validação
```sql
SELECT count(*) 
FROM TransactionError 
WHERE appName = 'fazendapro-backend' 
FACET error_type
```

## Melhores Práticas

1. **Sempre propague o contexto**: Passe o `context.Context` através das camadas (handler → service → repository)

2. **Use segmentos para operações importantes**: Identifique gargalos criando segmentos para operações críticas

3. **Adicione atributos relevantes**: Facilite a análise adicionando informações de negócio

4. **Reporte erros com contexto**: Use `NoticeError` e adicione atributos que ajudem no debugging

5. **Nomeie segmentos claramente**: Use nomes descritivos como "ValidateUserPermissions" ao invés de "Validate"

6. **Monitore operações assíncronas**: Crie transações explícitas para jobs em background

7. **Evite overhead desnecessário**: Não crie segmentos para operações muito rápidas (<1ms)

8. **Use cache de forma inteligente**: Monitore hit rate e ajuste TTL conforme necessário

## Troubleshooting

### Transação não aparece no New Relic

Verifique se:
1. O contexto está sendo propagado corretamente
2. A transação foi finalizada com `defer txn.End()`
3. O New Relic está habilitado (`NEW_RELIC_ENABLED=true`)

### Segmentos não aparecem

Certifique-se de:
1. Chamar `segment.End()` ou usar `defer segment.End()`
2. O segmento está dentro de uma transação ativa
3. O tempo de execução é significativo (>1ms)

### Erros não são reportados

Verifique:
1. Está usando `txn.NoticeError(err)` e não apenas retornando o erro
2. A transação está ativa quando o erro ocorre
3. O erro não está sendo silenciado em algum lugar

