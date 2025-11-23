package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type CacheInterface interface {
	Get(key string, dest interface{}) error
	Set(key string, value interface{}, expiration int32) error
	Delete(key string) error
	Increment(key string, delta uint64) (uint64, error)
}

type MemcacheClient struct {
	client *memcache.Client
}

func NewMemcacheClient(servers ...string) *MemcacheClient {
	if len(servers) == 0 {
		servers = []string{"localhost:11211"}
	}
	client := memcache.New(servers...)
	return &MemcacheClient{client: client}
}

func (m *MemcacheClient) Get(key string, dest interface{}) error {
	if m.client == nil {
		return ErrCacheMiss
	}

	item, err := m.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return ErrCacheMiss
		}
		log.Printf("Erro ao buscar do cache (key: %s): %v - fazendo fallback", key, err)
		return ErrCacheMiss
	}

	if err := json.Unmarshal(item.Value, dest); err != nil {
		log.Printf("Erro ao deserializar cache (key: %s): %v", key, err)
		return ErrCacheMiss
	}

	return nil
}

func (m *MemcacheClient) Set(key string, value interface{}, expiration int32) error {
	if m.client == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("Erro ao serializar para cache (key: %s): %v", key, err)
		return nil
	}

	item := &memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: expiration,
	}

	if err := m.client.Set(item); err != nil {
		log.Printf("Erro ao salvar no cache (key: %s): %v - continuando sem cache", key, err)
		return nil
	}

	return nil
}

func (m *MemcacheClient) Delete(key string) error {
	if m.client == nil {
		return nil
	}

	if err := m.client.Delete(key); err != nil {
		if err == memcache.ErrCacheMiss {
			return nil
		}
		log.Printf("Erro ao deletar do cache (key: %s): %v - continuando", key, err)
		return nil
	}
	return nil
}

func (m *MemcacheClient) Increment(key string, delta uint64) (uint64, error) {
	if m.client == nil {
		return delta, nil
	}

	newValue, err := m.client.Increment(key, delta)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			initialValue := int64(delta)
			item := &memcache.Item{
				Key:        key,
				Value:      []byte(fmt.Sprintf("%d", initialValue)),
				Expiration: int32((24 * time.Hour).Seconds()),
			}
			if err := m.client.Set(item); err != nil {
				log.Printf("Erro ao criar chave para incremento (key: %s): %v", key, err)
				return delta, nil
			}
			return delta, nil
		}
		log.Printf("Erro ao incrementar no cache (key: %s): %v", key, err)
		return delta, nil
	}
	return newValue, nil
}

var (
	ErrCacheMiss = fmt.Errorf("cache miss")
)
