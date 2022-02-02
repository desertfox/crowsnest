package job

import (
	"testing"
)

func Test_history(t *testing.T) {
	t.Run("Push", func(t *testing.T) {
		var results []Result
		for i := 0; i < 10; i++ {
			results = append(results, Result{Count: 0})
		}

		tests := []struct {
			name string
			push Result
			got  History
			want int
		}{
			{"One", Result{Count: 1}, History{}, 1},
			{"Five", Result{Count: 1}, History{results[:5]}, 6},
			{"Max", Result{Count: 1}, History{results: results}, 10},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				tt.got.Push(tt.push)

				if len(tt.got.results) != tt.want {
					t.Fatalf("Incorrect length got: %v, want: %v, results: %v", len(tt.got.results), tt.want, tt.got.results)
				}

				if tt.push != tt.got.results[0] {
					t.Fatalf("New results go on bottom got: %v, want: %v", tt.push, tt.got.results[0])
				}
			})
		}
	})
}
