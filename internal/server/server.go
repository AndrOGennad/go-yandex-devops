package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/AndrOGennad/go-yandex-devops/internal"
)

const address = "127.0.0.1:8080"

type MetricHandler struct {
	store Storage
}

func NewMetricHandler(store Storage) *MetricHandler {
	return &MetricHandler{store}
}

func (mh *MetricHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}
	/*if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}*/

	path := strings.Trim(r.URL.Path, " ")
	pathVars := strings.Split(path, "/")

	fmt.Printf("параметры запроса: %s", pathVars)

	if len(pathVars) < 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricTypeParam := pathVars[2]
	metricIDParam := pathVars[3]
	metricValueParam := pathVars[4]

	fmt.Printf("Path: (%s) Type: (%s) Name: (%s) Value: (%s)", pathVars[1], metricTypeParam, metricIDParam, metricValueParam)

	metric := internal.Metric{
		ID:   internal.ID(metricIDParam),
		Type: internal.Type(metricTypeParam),
	}
	switch metricTypeParam {
	case "counter":
		value, err := strconv.ParseInt(metricValueParam, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("получено плохое значение, конвертация в int64 провалилась")
			return
		}
		metric.Counter = internal.Counter(value)

	case "gauge":
		value, err := strconv.ParseFloat(metricValueParam, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("получено плохое значение, конвертация вo float64 провалилась")
			return
		}
		metric.Gauge = internal.Gauge(value)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Println("получена метрика неизвестного типа")
		return
	}

	_ = mh.store.Put(metric.ID, metric)
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	return
}

func Run(ctx context.Context) error {
	store := NewMemStorage()
	metricHandler := NewMetricHandler(store)
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", metricHandler.GetMetric)
	server := &http.Server{Handler: mux, Addr: address}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("ошибка сервера", err)
		}
	}()

	select {
	case <-ctx.Done():
		if err := server.Shutdown(ctx); err != nil {
			return err
		}
		return ctx.Err()
	}

}
