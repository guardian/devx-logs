package main

import (
	"context"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/vmihailenco/msgpack/v5"
)

type logConsumer struct {
	mu   sync.Mutex
	logs []string
}

func (c *logConsumer) Accept(log testcontainers.Log) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logs = append(c.logs, string(log.Content))
}

func (c *logConsumer) String() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return strings.Join(c.logs, "")
}

func TestFluentBit(t *testing.T) {
	ctx := context.Background()

	consumer := &logConsumer{}

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    ".",
			Dockerfile: "Dockerfile",
		},

		ExposedPorts: []string{"24224/tcp"},

		Env: map[string]string{
			"STACK":               "test",
			"STAGE":               "test",
			"APP":                 "test",
			"GU_REPO":             "test/repo",
			"TASK_NAME":           "test-task",
			"ECS_CLUSTER":         "test-cluster",
			"ECS_TASK_ARN":        "test-task-arn",
			"ECS_TASK_DEFINITION": "test-task-definition",

			// Test-only override:
			"OUTPUT_PLUGIN": "stdout",
		},
		//WaitingFor: wait.ForListeningPort("24224/tcp"),
	}
	container, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		t.Fatalf("starting Fluent Bit: %v", err)
	}
	defer container.Terminate(ctx)

	container.FollowOutput(consumer)

	if err := container.StartLogProducer(ctx); err != nil {
		t.Fatalf("starting log producer: %v", err)
	}
	defer container.StopLogProducer()

	t.Log("Fluent Bit container started")

	// Give Fluent Bit a moment to start and emit any startup errors.
	time.Sleep(2 * time.Second)

	output := consumer.String()
	t.Logf("Fluent Bit output:\n%s", output)

	if strings.Contains(output, "[error]") {
		t.Fatalf("Fluent Bit reported an error:\n%s", output)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := container.MappedPort(ctx, "24224/tcp")
	if err != nil {
		t.Fatal(err)
	}

	address := net.JoinHostPort(host, port.Port())

	sendForwardRecord(t, address)

	requireOutput(t, consumer, `"hello":"world"`)
}

func sendForwardRecord(t *testing.T, address string) {
	t.Helper()

	record := []interface{}{
		"application-logs",
		[]interface{}{
			[]interface{}{
				time.Now().Unix(),
				map[string]interface{}{
					"MESSAGE": `{"hello":"world","test":true}`,
				},
			},
		},
	}

	payload, err := msgpack.Marshal(record)
	if err != nil {
		t.Fatalf("encoding Forward record: %v", err)
	}

	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		t.Fatalf("connecting to Fluent Bit: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Write(payload); err != nil {
		t.Fatalf("sending Forward record: %v", err)
	}
}

func requireOutput(t *testing.T, consumer *logConsumer, expected string) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)

	for time.Now().Before(deadline) {
		if strings.Contains(consumer.String(), expected) {
			return
		}

		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf(
		"did not find %q in Fluent Bit output:\n%s",
		expected,
		consumer.String(),
	)
}
