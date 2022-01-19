package script_test

import (
	"context"
	"testing"

	"github.com/gorpher/gone/script"
)

func TestExec(t *testing.T) {
	expect := "hello"
	ouput, err := script.Exec(&script.Options{
		CancelCtx: context.Background(),
		Command:   "echo",
		CliArgs:   []string{expect},
		Env:       map[string]string{},
	})
	if err != nil {
		t.Fatal(err)
	}
	got := ouput.Stdout.String()
	if got != expect+"\n" {
		t.Fatalf("expect %s, but get %s", expect, got)
	}
}
