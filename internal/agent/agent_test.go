package agent

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AndrOGennad/go-yandex-devops/internal"
	"github.com/stretchr/testify/assert"
)

func TestSendMetric(t *testing.T) {
	type args struct {
		metric internal.Metric
	}

	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			"send_counter_ok",
			args{metric: internal.Metric{
				ID:      "id",
				Type:    "counter",
				Counter: 1,
			}},
			func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(200)
				return
			},
			false,
		},
		{
			"send_gauge_ok",
			args{metric: internal.Metric{
				ID:    "id",
				Type:  "gauge",
				Gauge: 1.1,
			}},
			func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(200)
				return
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.handler)
			defer srv.Close()

			sender := NewHTTPSender(strings.TrimPrefix(srv.URL, "http://"), srv.Client())
			err := sender.Send(tt.args.metric)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
