package agent

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/AndrOGennad/go-yandex-devops/internal"
)

const (
	address        = "127.0.0.1:8080"
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func GetGaugeMetrics() []internal.Metric {
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
	data := make([]internal.Metric, 0)
	data = append(data, internal.NewGauge("Alloc", internal.Gauge(stats.Alloc)))
	data = append(data, internal.NewGauge("BuckHashSys", internal.Gauge(stats.BuckHashSys)))
	data = append(data, internal.NewGauge("Frees", internal.Gauge(stats.Frees)))
	data = append(data, internal.NewGauge("GCCPUFraction", internal.Gauge(stats.GCCPUFraction)))
	data = append(data, internal.NewGauge("GCSys", internal.Gauge(stats.GCSys)))
	data = append(data, internal.NewGauge("HeapAlloc", internal.Gauge(stats.HeapAlloc)))
	data = append(data, internal.NewGauge("HeapIdle", internal.Gauge(stats.HeapIdle)))
	data = append(data, internal.NewGauge("HeapInuse", internal.Gauge(stats.HeapInuse)))
	data = append(data, internal.NewGauge("HeapObjects", internal.Gauge(stats.HeapObjects)))
	data = append(data, internal.NewGauge("HeapReleased", internal.Gauge(stats.HeapReleased)))
	data = append(data, internal.NewGauge("HeapSys", internal.Gauge(stats.HeapSys)))
	data = append(data, internal.NewGauge("LastGC", internal.Gauge(stats.LastGC)))
	data = append(data, internal.NewGauge("Lookups", internal.Gauge(stats.Lookups)))
	data = append(data, internal.NewGauge("MCacheInuse", internal.Gauge(stats.MCacheInuse)))
	data = append(data, internal.NewGauge("MCacheSys", internal.Gauge(stats.MCacheSys)))
	data = append(data, internal.NewGauge("MSpanInuse", internal.Gauge(stats.MSpanInuse)))
	data = append(data, internal.NewGauge("MSpanSys", internal.Gauge(stats.MSpanSys)))
	data = append(data, internal.NewGauge("Mallocs", internal.Gauge(stats.Mallocs)))
	data = append(data, internal.NewGauge("NextGC", internal.Gauge(stats.NextGC)))
	data = append(data, internal.NewGauge("NumForcedGC", internal.Gauge(stats.NumForcedGC)))
	data = append(data, internal.NewGauge("NumGC", internal.Gauge(stats.NumGC)))
	data = append(data, internal.NewGauge("OtherSys", internal.Gauge(stats.OtherSys)))
	data = append(data, internal.NewGauge("PauseTotalNs", internal.Gauge(stats.PauseTotalNs)))
	data = append(data, internal.NewGauge("StackInuse", internal.Gauge(stats.StackInuse)))
	data = append(data, internal.NewGauge("StackSys", internal.Gauge(stats.StackSys)))
	data = append(data, internal.NewGauge("Sys", internal.Gauge(stats.Sys)))
	data = append(data, internal.NewGauge("TotalAlloc", internal.Gauge(stats.TotalAlloc)))
	data = append(data, internal.NewGauge("RandomValue", internal.Gauge(rand.Float64())))
	return data
}

type Sender interface {
	Send(metric internal.Metric) error
}

type HTTPSender struct {
	address string
	client  *http.Client
}

func NewHTTPSender(address string, client *http.Client) *HTTPSender {
	return &HTTPSender{address, client}
}

func (hs *HTTPSender) Send(metric internal.Metric) error {
	_url := url.URL{
		Scheme: "http",
		Host:   hs.address,
		Path:   fmt.Sprintf("update/%s/%s/%s", metric.Type, metric.ID, metric.Value()),
	}
	response, err := hs.client.Post(_url.String(), "text/plain", nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		return nil
	}
	return errors.New(fmt.Sprintln("сервер вернул плохой код: ", response.Status))
}

func Run(ctx context.Context) error {
	sender := NewHTTPSender(address, http.DefaultClient)
	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	var pollCounter internal.Counter
	var metrics []internal.Metric
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-pollTicker.C:
			metrics = GetGaugeMetrics()
			pollCounter += 1
			counterMetric := internal.NewCounter("PollCount", pollCounter)
			metrics = append(metrics, counterMetric)
			fmt.Println("обновили метрики, текущий PollCount: ", pollCounter)
		case <-reportTicker.C:
			for _, metric := range metrics {
				err := sender.Send(metric)
				if err != nil {
					fmt.Println("ошибка при отправке метрики: ", err)
				}
			}
		}
	}
}
