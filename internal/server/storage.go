package server

import (
	"errors"

	"github.com/AndrOGennad/go-yandex-devops/internal"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	Get(key internal.ID) (metric internal.Metric, error error)
	Put(key internal.ID, value internal.Metric) (newValue internal.Metric, error error)
}

type MemStorage struct {
	data map[internal.ID]internal.Metric
}

func (m *MemStorage) Get(key internal.ID) (metric internal.Metric, error error) {
	if metric, ok := m.data[key]; !ok {
		return metric, ErrNotFound
	}
	return m.data[key], nil
}
func NewMemStorage() *MemStorage {
	data := make(map[internal.ID]internal.Metric)
	return &MemStorage{data}
}

func (m *MemStorage) Put(key internal.ID, value internal.Metric) (newValue internal.Metric, error error) {

	oldValue, ok := m.data[key]
	if !ok || value.Type == "gauge" {
		m.data[key] = value
		return m.data[key], nil

	}

	oldValue.Counter += value.Counter
	m.data[key] = oldValue

	return m.data[key], nil

	/*if oldValue, exists := m.data[key]; exists {
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
	return newValue*/
}
