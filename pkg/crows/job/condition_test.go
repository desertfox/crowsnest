package job

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
			result    Result
			threshold int
			state     string
			alert     bool
			alertStr  string
		}{
			{Result{Count: 0}, 1, ">", false, check},
			{Result{Count: 1}, 0, ">", true, fire},
			{Result{Count: 0}, 0, ">", true, fire},
			{Result{Count: 1}, 0, "<", false, check},
			{Result{Count: 0}, 1, "<", true, fire},
			{Result{Count: 0}, 0, "<", true, fire},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(fmt.Sprintf("%d %s= %d", tt.result.Count, tt.state, tt.threshold), func(t *testing.T) {
				got := Condition{
					Threshold: tt.threshold,
					State:     tt.state,
				}

				if got.IsAlert(tt.result) != tt.alert {
					t.Fatalf("got: %v, want: %v", got.IsAlert(tt.result), tt.alert)
				}

				if !strings.Contains(got.IsAlertText(tt.result), tt.alertStr) {
					t.Fatalf("got: %v, want: %v", got.IsAlertText(tt.result), tt.alertStr)
				}
			})
		}
	})
}

func ExampleCondition() {
	fmt.Println(Condition{
		Threshold: 1,
		State:     "<",
	})
	//Output: {1 < 0}
}
