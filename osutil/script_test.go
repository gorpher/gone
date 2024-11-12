package osutil_test

import (
	"context"
	"github.com/gorpher/gone/osutil"
	"testing"
)

func TestExec(t *testing.T) {
	expect := "hello"
	ouput, err := osutil.Exec(&osutil.Options{
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
