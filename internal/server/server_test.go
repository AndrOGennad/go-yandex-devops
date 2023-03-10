package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndrOGennad/go-yandex-devops/internal"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type MemStorageMock struct {
	metric internal.Metric
}

func (m MemStorageMock) Get(internal.ID) (internal.Metric, error) {
	return m.metric, nil
}

func (m MemStorageMock) Put(internal.ID, internal.Metric) (newValue internal.Metric, error error) {
	return m.metric, nil
}

func TestMetricHandler_PutMetric(t *testing.T) {
	storeMock := MemStorageMock{}
	metricHandler := NewMetricHandler(storeMock)

	type request struct {
		method      string
		path        string
		pathPattern string
		handler     http.HandlerFunc
	}

	type response struct {
		code int
	}

	tests := []struct {
		name    string
		request request
		want    response
	}{
		{
			"valid_counter_200",
			request{
				method:      http.MethodPost,
				path:        "/update/counter/testCounter/100",
				pathPattern: "/update/{type}/{id}/{value}",
				handler:     metricHandler.PutMetric,
			},
			response{code: 200},
		},
		{
			"valid_gauge_200",
			request{
				method:      http.MethodPost,
				path:        "/update/gauge/testGauge/10.0",
				pathPattern: "/update/{type}/{id}/{value}",
				handler:     metricHandler.PutMetric,
			},
			response{code: 200},
		},
		{
			"no_value_counter_404",
			request{
				method:      http.MethodPost,
				path:        "/update/counter/testCounter",
				pathPattern: "/update/{type}/{id}/{value}",
				handler:     metricHandler.PutMetric,
			},
			response{code: 404},
		},
		{
			"no_value_gauge_404",
			request{
				method:      http.MethodPost,
				path:        "/update/gauge/testGauge",
				pathPattern: "/update/{type}/{id}/{value}",
				handler:     metricHandler.PutMetric,
			},
			response{code: 404},
		},
		{
			"invalid_counter_400",
			request{
				method:      http.MethodPost,
				path:        "/update/counter/testCounter/none",
				pathPattern: "/update/{type}/{id}/{value}",
				handler:     metricHandler.PutMetric,
			},
			response{code: 400},
		},
		{
			"invalid_gauge_400",
			request{
				method:      http.MethodPost,
				path:        "/update/gauge/testGauge/none",
				pathPattern: "/update/{type}/{id}/{value}",
				handler:     metricHandler.PutMetric,
			},
			response{code: 400},
		},
		{
			"invalid_type_501",
			request{
				method:      http.MethodPost,
				path:        "/update/invalid_type/testCounter/100",
				pathPattern: "/update/{type}/{id}/{value}",
				handler:     metricHandler.PutMetric,
			},
			response{code: 501},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := chi.NewRouter()
			mux.Method(tt.request.method, tt.request.pathPattern, tt.request.handler)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(tt.request.method, tt.request.path, nil)
			mux.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
