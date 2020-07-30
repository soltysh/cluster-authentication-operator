package endpointaccessible

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/events"
)

func Test_endpointAccessibleController_sync(t *testing.T) {
	tests := []struct {
		name           string
		endpointListFn EndpointListFunc
		wantErr        bool
	}{
		{
			name: "all endpoints working",
			endpointListFn: func() ([]string, error) {
				return []string{"https://google.com"}, nil
			},
		},
		{
			name: "endpoints lister error",
			endpointListFn: func() ([]string, error) {
				return nil, fmt.Errorf("some error")
			},
			wantErr: true,
		},
		{
			name: "non working endpoints",
			endpointListFn: func() ([]string, error) {
				return []string{"https://google.com", "https://nonexistenturl.com"}, nil
			},
			wantErr: true,
		},
		{
			name: "invalid url",
			endpointListFn: func() ([]string, error) {
				return []string{"htt//bad`string"}, nil
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &endpointAccessibleController{
				endpointListFn: tt.endpointListFn,
			}
			if err := c.sync(context.Background(), factory.NewSyncContext(tt.name, events.NewInMemoryRecorder(tt.name))); (err != nil) != tt.wantErr {
				t.Errorf("sync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_toHealthzURL(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "test urls",
			args: []string{"a", "b"},
			want: []string{"https://a/healthz", "https://b/healthz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toHealthzURL(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toHealthzURL() = %v, want %v", got, tt.want)
			}
		})
	}
}