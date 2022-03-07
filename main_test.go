package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

func TestRootCmd(t *testing.T) {
	// Stub EC2 metadata service
	metadataStubURL := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tags := []string{"App", "Stack", "Stage", "gu:cdk:version"}

		if r.URL.Path == "/latest/api/token" {
			fmt.Fprintln(w, "12345")
			return
		}

		switch {
		case r.URL.Path == "/latest/api/token":
			fmt.Fprintln(w, "12345")
		case r.URL.Path == "/latest/meta-data/tags/instance":
			fmt.Fprintln(w, strings.Join(tags, "\n"))
		case strings.HasPrefix(r.URL.Path, "/latest/meta-data/tags/instance/"):
			tag := strings.TrimPrefix(r.URL.Path, "/latest/meta-data/tags/instance/")
			fmt.Fprintln(w, tag+"-VALUE")
		default: // individual tag
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Unknown request path: %s", r.URL.Path)
		}

	}))
	defer metadataStubURL.Close()

	cmd := RootCmd(metadataStubURL.URL)

	got := bytes.Buffer{}
	cmd.SetOut(&got)
	cmd.SetArgs([]string{"--logStreamARN", "foo", "--systemdUnit", "bar"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	want, _ := os.ReadFile("fluentbit.test.conf")

	if got.String() != string(want) {
		t.Error(diff.Diff(string(want), got.String()))
	}
}
