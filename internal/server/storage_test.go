package server

import (
	"testing"

	"github.com/AndrOGennad/go-yandex-devops/internal"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_Get(t *testing.T) {
	type fields struct {
		data map[internal.ID]internal.Metric
	}
	type args struct {
		key internal.ID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    internal.Metric
		wantErr bool
	}{
		{
			name: "counter metric exists for key",
			fields: fields{data: map[internal.ID]internal.Metric{
				internal.ID("id"): {
					ID:      internal.ID("id"),
					Type:    "type",
					Counter: 1,
				},
			}},
			args: args{key: internal.ID("id")},
			want: internal.Metric{
				ID:      internal.ID("id"),
				Type:    "type",
				Counter: 1,
			},
			wantErr: false,
		},
		{
			name: "gauge metric exists for key",
			fields: fields{data: map[internal.ID]internal.Metric{
				internal.ID("id"): {
					ID:    internal.ID("id"),
					Type:  "type",
					Gauge: 1.1,
				},
			}},
			args: args{key: internal.ID("id")},
			want: internal.Metric{
				ID:    internal.ID("id"),
				Type:  "type",
				Gauge: 1.1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{data: tt.fields.data}
			got, err := storage.Get(tt.args.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMemStorage_Put(t *testing.T) {
	type fields struct {
		data map[internal.ID]internal.Metric
	}
	type args struct {
		key   internal.ID
		value internal.Metric
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantNewValue internal.Metric
		wantErr      bool
	}{
		{
			"counter metric not found",
			fields{data: map[internal.ID]internal.Metric{}},
			args{
				key: "id",
				value: internal.Metric{
					ID:      "id",
					Type:    "counter",
					Counter: 1,
				},
			},
			internal.Metric{
				ID:      "id",
				Type:    "counter",
				Counter: 1,
			},
			false,
		},
		{
			"gauge metric not found",
			fields{data: map[internal.ID]internal.Metric{}},
			args{
				key: "id",
				value: internal.Metric{
					ID:    "id",
					Type:  "gauge",
					Gauge: 1.1,
				},
			},
			internal.Metric{
				ID:    "id",
				Type:  "gauge",
				Gauge: 1.1,
			},
			false,
		},
		{
			"counter metric found",
			fields{data: map[internal.ID]internal.Metric{
				internal.ID("id"): {
					ID:      "id",
					Type:    "counter",
					Counter: 1,
				},
			}},
			args{
				key: "id",
				value: internal.Metric{
					ID:      "id",
					Type:    "counter",
					Counter: 2,
				},
			},
			internal.Metric{
				ID:      "id",
				Type:    "counter",
				Counter: 3,
			},
			false,
		},
		{
			"gauge metric found",
			fields{data: map[internal.ID]internal.Metric{
				internal.ID("id"): {
					ID:    "id",
					Type:  "gauge",
					Gauge: 1.1,
				},
			}},
			args{
				key: "id",
				value: internal.Metric{
					ID:    "id",
					Type:  "gauge",
					Gauge: 2.2,
				},
			},
			internal.Metric{
				ID:    "id",
				Type:  "gauge",
				Gauge: 2.2,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &MemStorage{data: tt.fields.data}
			got, err := store.Put(tt.args.key, tt.args.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, tt.wantNewValue)
			written, exists := tt.fields.data[tt.args.key]
			assert.True(t, exists)
			assert.Equal(t, got, written)
		})
	}
}
