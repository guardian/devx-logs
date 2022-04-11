package main

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

func TestRootCmd(t *testing.T) {
	cmd := RootCmd()

	got := bytes.Buffer{}
	cmd.SetOut(&got)
	cmd.SetArgs([]string{"--kinesisStreamName", "foo", "--tags", "app=bar,stack=test,stage=CODE", "--dry-run"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	want, _ := os.ReadFile("fluentbit/fluentbit.test.conf")

	if got.String() != string(want) {
		t.Error(diff.Diff(string(want), got.String()))
	}
}

func TestNormaliseTags(t *testing.T) {
	tags := map[string]string{"App": "foo", "Stage": "PROD", "Stack": "deploy", "Name": "foo"}

	got := normaliseTags(tags)
	want := map[string]string{"app": "foo", "stage": "PROD", "stack": "deploy", "Name": "foo"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v and want %v", got, want)
	}
}

func TestOutput(t *testing.T) {
	// docker container with a service to target (just outputs a line every second)
	//
}
