package printline

import (
	"testing"
)

// func prepare() (pl *PrintLine, buff bytes.Buffer) {
// 	fmt.Println(&buff)
// 	pl = New(&buff)
// 	return
// }

func TestAppendSpace(t *testing.T) {
	res := appendSpace("", 5)
	if res != "     " {
		t.Logf("expect 5 spaces, got: '%s'", res)
		t.Fail()
	}
	res = appendSpace("awesome", 2)
	if res != "awesome" {
		t.Logf("expect awesome, got: '%s'", res)
		t.Fail()
	}
}
