package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AndrOGennad/go-yandex-devops/internal"
	"github.com/go-chi/chi/v5"
)

const address = "127.0.0.1:8080"

type MetricHandler struct {
	store Storage
}

func NewMetricHandler(store Storage) *MetricHandler {
	return &MetricHandler{store}
}

func (mh *MetricHandler) PutMetric(w http.ResponseWriter, r *http.Request) {

	metricTypeParam := chi.URLParam(r, "type")
	metricIDParam := chi.URLParam(r, "id")
	metricValueParam := chi.URLParam(r, "value")

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
		// код 422 не проходит авто тесты
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Println("получена метрика неизвестного типа")
		return
	}

	_, _ = mh.store.Put(metric.ID, metric)
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	return
}

func (mh *MetricHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricTypeParam := chi.URLParam(r, "type")
	metricIDParam := chi.URLParam(r, "id")

	if metricTypeParam != "counter" && metricTypeParam != "gauge" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metric, err := mh.store.Get(internal.ID(metricIDParam))
	if err == ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(metric.Value()))
	return
}

func Run(ctx context.Context) error {
	store := NewMemStorage()
	metricHandler := NewMetricHandler(store)

	mux := chi.NewRouter()
	mux.Get("/value/{type}/{id}", metricHandler.GetMetric)
	mux.Post("/update/{type}/{id}/{value}", metricHandler.PutMetric)

	server := &http.Server{Handler: mux, Addr: address}

	errCh := make(chan error)
	defer close(errCh)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			return err
		} else {
			return nil
		}
	case err := <-errCh:
		fmt.Println("ошибка сервера", err)
		return err
	}
}
