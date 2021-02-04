package gone

import "testing"

func TestMacAddr(t *testing.T) {
	if _, err := MacAddr(); err != nil {
		t.Error(err)
		return
	}
}
