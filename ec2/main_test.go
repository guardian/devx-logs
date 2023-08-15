package main

import (
	"bytes"
	"fmt"
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

	fluentbitConfig, _ := os.ReadFile("fluentbit/test-configs/fluentbit-cloud-init-only.test.conf")
	want := fmt.Sprintf("Main config:\n%s\nApplication config:%s", string(fluentbitConfig), "")

	if got.String() != want {
		t.Error(diff.Diff(want, got.String()))
	}
}

func TestGenerateConfigs(t *testing.T) {

	cloudInitOnly, _ := os.ReadFile("fluentbit/test-configs/fluentbit-cloud-init-only.test.conf")
	cloudInitWithInclude, _ := os.ReadFile("fluentbit/test-configs/fluentbit-cloud-init-and-app.test.conf")
	applicationOnly, _ := os.ReadFile("fluentbit/test-configs/fluentbit-app-only.test.conf")
	application, _ := os.ReadFile("fluentbit/test-configs/application-logs.test.conf")

	withoutApplicationLogs := FluentbitConfig{MainConfigFile: string(cloudInitOnly)}
	withApplicationLogs := FluentbitConfig{MainConfigFile: string(cloudInitWithInclude), ApplicationConfigFile: string(application)}
	withoutCloudInitLogs := FluentbitConfig{MainConfigFile: string(applicationOnly)}

	var tests = []struct {
		tagsArg string
		want    FluentbitConfig
	}{
		{"app=bar,stack=test,stage=CODE", withoutApplicationLogs},
		{"app=bar,stack=test,stage=CODE,SystemdUnit=bar.service", withApplicationLogs},
		{"app=bar,stack=test,stage=CODE,SystemdUnit=bar.service,DisableCloudInitLogs=true", withoutCloudInitLogs},
	}

	for _, testCase := range tests {
		got := generateConfigs(testCase.tagsArg, "foo", false)

		if got.MainConfigFile != testCase.want.MainConfigFile {
			t.Error(diff.Diff(got.MainConfigFile, testCase.want.MainConfigFile))
		}
		if got.ApplicationConfigFile != testCase.want.ApplicationConfigFile {
			t.Error(diff.Diff(got.ApplicationConfigFile, testCase.want.ApplicationConfigFile))
		}
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
