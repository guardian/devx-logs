package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
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
	application, _ := os.ReadFile("fluentbit/test-configs/application-logs.test.conf")

	withoutApplicationLogs := FluentbitConfig{MainConfigFile: string(cloudInitOnly)}
	withApplicationLogs := FluentbitConfig{MainConfigFile: string(cloudInitWithInclude), ApplicationConfigFile: string(application)}

	var tests = []struct {
		tagsArg string
		want    FluentbitConfig
	}{
		{"app=bar,stack=test,stage=CODE", withoutApplicationLogs},
		{"app=bar,stack=test,stage=CODE,SystemdUnit=bar.service", withApplicationLogs},
	}

	for _, testCase := range tests {
		got := generateConfigs(testCase.tagsArg, "foo")

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

func TestWriteConfigs(t *testing.T) {
	const modern = "fluent-bit/fluent-bit.conf"
	const legacy = "td-agent-bit/td-agent-bit.conf"
	for _, tc := range []struct {
		name      string
		installed []string
		selected  string
	}{
		{"modern", []string{modern}, modern},
		{"legacy", []string{legacy}, legacy},
		{"prefer modern", []string{legacy, modern}, modern},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			for _, path := range tc.installed {
				fullPath := filepath.Join(root, path)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(fullPath, []byte("original"), 0644); err != nil {
					t.Fatal(err)
				}
			}
			config := generateConfigs("app=bar,SystemdUnit=bar.service", "foo")
			for range 2 {
				if err := writeConfigs(config, root); err != nil {
					t.Fatal(err)
				}
			}
			for _, path := range tc.installed {
				want := "original"
				if path == tc.selected {
					want = config.MainConfigFile
				}
				got, err := os.ReadFile(filepath.Join(root, path))
				if err != nil || string(got) != want {
					t.Fatalf("config %s: got %q, want %q, error %v", path, got, want, err)
				}
			}
			applicationPath := filepath.Join(root, filepath.Dir(tc.selected), "application-logs.conf")
			got, err := os.ReadFile(applicationPath)
			if err != nil || string(got) != config.ApplicationConfigFile {
				t.Fatalf("application config: got %q, want %q, error %v", got, config.ApplicationConfigFile, err)
			}
			// Removing application log shipping must clear the previous config.
			if err := writeConfigs(FluentbitConfig{MainConfigFile: "cloud-init only"}, root); err != nil {
				t.Fatal(err)
			}
			got, err = os.ReadFile(applicationPath)
			if err != nil || len(got) != 0 {
				t.Fatalf("application config was not cleared: %q, error %v", got, err)
			}
		})
	}
}

func TestWriteConfigsWithoutInstallation(t *testing.T) {
	root := t.TempDir()
	if err := writeConfigs(FluentbitConfig{}, root); err == nil {
		t.Fatal("expected an error when neither package is installed")
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatalf("unexpected files created: %v, error %v", entries, err)
	}
}

func TestWriteConfigsPreservesMainOnWriteFailure(t *testing.T) {
	root := t.TempDir()
	modernDir := filepath.Join(root, "fluent-bit")
	if err := os.MkdirAll(filepath.Join(modernDir, "application-logs.conf"), 0755); err != nil {
		t.Fatal(err)
	}
	mainPath := filepath.Join(modernDir, "fluent-bit.conf")
	if err := os.WriteFile(mainPath, []byte("original"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := writeConfigs(FluentbitConfig{MainConfigFile: "replacement"}, root); err == nil {
		t.Fatal("expected a write error")
	}
	got, err := os.ReadFile(mainPath)
	if err != nil || string(got) != "original" {
		t.Fatalf("main config changed after a failed application config write: %q, error %v", got, err)
	}
}

func TestOutput(t *testing.T) {
	// docker container with a service to target (just outputs a line every second)
	//
}
