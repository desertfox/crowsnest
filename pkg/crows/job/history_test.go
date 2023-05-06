package job

import (
	"testing"
)

func Test_history(t *testing.T) {
	t.Run("Push", func(t *testing.T) {
		var results []*Result
		for i := 0; i < maxHistory; i++ {
			results = append(results, &Result{Count: 0})
		}

		tests := []struct {
			name string
			push *Result
			got  *History
			want int
		}{
			{"One", &Result{Count: 1}, &History{AlertCount: 0}, 1},
			{"Five", &Result{Count: 1}, &History{Results: results[:5], AlertCount: 0}, 6},
			{"Max", &Result{Count: 1}, &History{Results: results, AlertCount: 0}, maxHistory},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				tt.got.Add(tt.push)

				if len(tt.got.Results) != tt.want {
					t.Fatalf("Incorrect length got: %v, want: %v, results: %v", len(tt.got.Results), tt.want, tt.got.Results)
				}

				if tt.push != tt.got.Results[0] {
					t.Fatalf("New results go on bottom got: %v, want: %v", tt.push, tt.got.Results[0])
				}
			})
		}
	})
}
