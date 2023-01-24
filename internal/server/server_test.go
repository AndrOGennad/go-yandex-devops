package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndrOGennad/go-yandex-devops/internal"
	"github.com/stretchr/testify/assert"
)

type MemStorageMock struct {
	metric internal.Metric
}

func (m MemStorageMock) Get(key internal.ID) internal.Metric {
	return m.metric
}

func (m MemStorageMock) Put(key internal.ID, value internal.Metric) (newValue internal.Metric) {
	return m.metric
}

func TestMetricHandler_PutMetric(t *testing.T) {
	storeMock := MemStorageMock{}
	metricHandler := NewMetricHandler(storeMock)

	type request struct {
		method  string
		path    string
		handler http.HandlerFunc
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
				method:  http.MethodPost,
				path:    "/update/counter/testCounter/100",
				handler: metricHandler.PutMetric,
			},
			response{code: 200},
		},
		{
			"valid_gauge_200",
			request{
				method:  http.MethodPost,
				path:    "/update/gauge/testGauge/10.0",
				handler: metricHandler.PutMetric,
			},
			response{code: 200},
		},
		{
			"no_value_counter_404",
			request{
				method:  http.MethodPost,
				path:    "/update/counter/testCounter",
				handler: metricHandler.PutMetric,
			},
			response{code: 404},
		},
		{
			"no_value_gauge_404",
			request{
				method:  http.MethodPost,
				path:    "/update/gauge/testGauge",
				handler: metricHandler.PutMetric,
			},
			response{code: 404},
		},
		{
			"invalid_counter_400",
			request{
				method:  http.MethodPost,
				path:    "/update/counter/testCounter/none",
				handler: metricHandler.PutMetric,
			},
			response{code: 400},
		},
		{
			"invalid_gauge_400",
			request{
				method:  http.MethodPost,
				path:    "/update/gauge/testGauge/none",
				handler: metricHandler.PutMetric,
			},
			response{code: 400},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			request := httptest.NewRequest(tt.request.method, tt.request.path, nil)
			tt.request.handler(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
