package server

import "github.com/AndrOGennad/go-yandex-devops/internal"

type Storage interface {
	Get(key internal.ID) internal.Metric
	Put(key internal.ID, value internal.Metric) (newValue internal.Metric)
}

type MemStorage struct {
	data map[internal.ID]internal.Metric
}

func (m *MemStorage) Get(key internal.ID) internal.Metric {
	return m.data[key]
}
func NewMemStorage() *MemStorage {
	data := make(map[internal.ID]internal.Metric)
	return &MemStorage{data}
}

func (m *MemStorage) Put(key internal.ID, value internal.Metric) (newValue internal.Metric) {

	if oldValue, exists := m.data[key]; exists {
		if oldValue.Type == "counter" {
			newValue = internal.Metric{
				ID:      value.ID,
				Type:    value.Type,
				Gauge:   0,
				Counter: oldValue.Counter + value.Counter,
			}
			m.data[key] = newValue
		} else {
			newValue = value
		}
	} else {
		newValue = value
	}

	m.data[key] = newValue
	return newValue
}
