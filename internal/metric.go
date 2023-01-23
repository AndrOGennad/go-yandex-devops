package internal

import (
	"fmt"
	"strconv"
)

type ID string
type Type string
type Counter int64

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

type Gauge float64

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', 3, 64)
}

type Metric struct {
	ID      ID
	Type    Type
	Gauge   Gauge
	Counter Counter
}

func (m Metric) Value() string {
	switch m.Type {
	case "gauge":
		return m.Gauge.String()
	case "counter":
		return m.Counter.String()
	default:
		return "Unknown Metric"
	}
}

func (m Metric) String() string {
	return fmt.Sprintf("Metric. ID: %s Type: %s Value: %s", m.ID, m.Type, m.Value())
}

func NewGauge(id ID, value Gauge) Metric {
	return Metric{id, "gauge", value, 0}
}

func NewCounter(id ID, value Counter) Metric {
	return Metric{id, "counter", 0.0, value}
}
