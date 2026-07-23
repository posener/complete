package complete

import "testing"

func TestNilCommandPredict(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked: %v", r)
		}
	}()
	var c *Command
	if c.Predict(Args{}) != nil {
		t.Fatal("want nil slice")
	}
}
