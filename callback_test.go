package gormetrics

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/prometheus/client_golang/prometheus"
)

func TestMergeLabels(t *testing.T) {
	tests := []struct {
		a, b, want prometheus.Labels
	}{
		{
			a: prometheus.Labels{"test": "test", "a": "b"},
			b: prometheus.Labels{"test": "test2", "b": "a"},
			want: prometheus.Labels{
				"test": "test",
				"a":    "b",
				"b":    "a",
			},
		},
	}

	for _, tc := range tests {
		got := mergeLabels(tc.a, tc.b)

		if diff := deep.Equal(tc.want, got); diff != nil {
			t.Fatal(diff)
		}
	}
}
