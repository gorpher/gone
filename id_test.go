package gone

import "testing"

func TestIDInt64(t *testing.T) {
	if IDInt64() == 0 {
		t.Error("Generate int64 ID error")
	}
}
func TestIDString(t *testing.T) {
	t.Log(IDString())
	if IDString() == "" {
		t.Error("Generate string ID error")
	}
}
