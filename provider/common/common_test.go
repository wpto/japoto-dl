package common

import "testing"

func TestEncodeIdx(t *testing.T) {
	cases := []struct {
		args   []int
		result string
	}{
		{args: []int{1}, result: "1"},
		{args: []int{1, 2}, result: "1z2"},
		{args: []int{34, 34}, result: "yzy"},
	}

	for _, c := range cases {
		res := EncodeIdx(c.args[0], c.args[1:]...)
		if res != c.result {
			t.Logf("expected %s, got %s for %v", c.result, res, c.args)
			t.Fail()
		}
	}
}
