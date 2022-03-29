package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

func TestRootCmd(t *testing.T) {
	cmd := RootCmd()

	got := bytes.Buffer{}
	cmd.SetOut(&got)
	cmd.SetArgs([]string{"--kinesisStreamName", "foo", "--tags", "App=bar,Stack=test,Stage=CODE"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	want, _ := os.ReadFile("fluentbit/fluentbit.test.conf")

	if got.String() != string(want) {
		t.Error(diff.Diff(string(want), got.String()))
	}
}

func TestOutput(t *testing.T) {
	// docker container with a service to target (just outputs a line every second)
	//
}
