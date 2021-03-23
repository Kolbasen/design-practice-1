package binary

import "testing"

func TestMain(t *testing.T) {
	res := Hello()

	if res != 1 {
		t.Errorf("Result is incorrect")
	}
}
