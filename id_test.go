package gone

import (
	"testing"
)

func TestIDInt64(t *testing.T) {
	if ID.SInt64() == 0 {
		t.Error("Generate int64 ID error")
	}
}

func TestIDString(t *testing.T) {
	t.Log(ID.SString())
	if ID.SString() == "" {
		t.Error("Generate string ID error")
	}
}

func TestId_XID(t *testing.T) {
	t.Log(ID.XString())
	if ID.XString() == "" {
		t.Error("Generate string ID error")
	}
	t.Log(ID.XID().String())
	if ID.XID().String() == "" {
		t.Error("Generate string ID error")
	}
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
}

func TestId_RandString(t *testing.T) {
	i := ID.RandString(11)
	t.Log(i)
	if i == "" || len(i) != 11 {
		t.Error("Generate string ID error")
	}
}
