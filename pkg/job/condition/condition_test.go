package condition

import (
	"fmt"
	"strings"
	"testing"
)

var (
	fire  string = "🔥"
	check string = "✔️"
)

func Test_condition(t *testing.T) {
	t.Run("condition", func(t *testing.T) {
		tests := []struct {
			count     int
			threshold int
			state     string
			alert     bool
			alertStr  string
		}{
			{0, 1, ">", false, check},
			{1, 0, ">", true, fire},
			{0, 0, ">", true, fire},
			{1, 0, "<", false, check},
			{0, 1, "<", true, fire},
			{0, 0, "<", true, fire},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(fmt.Sprintf("%d %s= %d", tt.count, tt.state, tt.threshold), func(t *testing.T) {
				got := Condition{
					Threshold: tt.threshold,
					State:     tt.state,
				}

				if got.IsAlert(tt.count) != tt.alert {
					t.Fatalf("got: %v, want: %v", got.IsAlert(tt.count), tt.alert)
				}

				if !strings.Contains(got.IsAlertText(tt.count), tt.alertStr) {
					t.Fatalf("got: %v, want: %v", got.IsAlertText(tt.count), tt.alertStr)
				}

			})
		}
	})
}

func testCondition() Condition {
	return Condition{
		Threshold: 1,
		State:     "<",
	}
}
