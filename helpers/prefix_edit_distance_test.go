package helpers

import (
	"fmt"
	"testing"
)

func TestPED(t *testing.T) {
	dist := PrefixEditDistance("", "")
	fmt.Println(dist)
	str := "gaituri2023"
	pref := "gaitsuri"
	dist = PrefixEditDistance(str, pref)
	for _, line := range dist {
		fmt.Println(line)
	}
	i := 6
	j := 4

	fmt.Printf("%s for pref %s has %d", str[:i], pref[:j], dist[i][j])
	t.Fail()
}
