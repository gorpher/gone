package authed

import "testing"

func TestCreateToken(t *testing.T) {
	authed := NewAuthed()
	token, refresh, err := authed.CreateToken(&UserSession{
		ID: "123",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	t.Log(refresh)
	payload, err := authed.VerifyToken(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(payload.GetToken())
}
