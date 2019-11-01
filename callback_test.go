// Copyright 2019 Profects Group B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
