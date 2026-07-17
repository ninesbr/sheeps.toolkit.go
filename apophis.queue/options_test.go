package apophis

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestConsumerParallelismOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		configure func(int) func(*options)
	}{
		{name: "correct spelling", configure: WithConsumerParallelism},
		{name: "deprecated spelling", configure: WithConsumerParralelism},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ops := NewOptions(test.configure(4))
			if ops.consumerParallelism != 4 {
				t.Fatalf("consumer parallelism = %d, want 4", ops.consumerParallelism)
			}
		})
	}
}

func TestInsecureOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		configure func(bool) func(*options)
	}{
		{name: "correct spelling", configure: WithInsecure},
		{name: "deprecated spelling", configure: WithInsecured},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ops := NewOptions(test.configure(true))
			if !ops.insecure {
				t.Fatal("insecure = false, want true")
			}
		})
	}
}

func TestOptionsRejectInvalidRequestTimeout(t *testing.T) {
	t.Parallel()

	ops := NewOptions(
		WithHost("localhost"),
		WithPort(50051),
		WithQueueName("queue-1"),
		WithRequestTimeout(0),
	)
	if err := ops.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want invalid request timeout error")
	}
}

func TestOptionsRejectInvalidConsumerSettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		configure func(*options)
		wantError string
	}{
		{
			name:      "zero reconnect interval",
			configure: WithReconnectInterval(0),
			wantError: "reconnect interval must be greater than zero",
		},
		{
			name:      "negative reconnect interval",
			configure: WithReconnectInterval(-time.Second),
			wantError: "reconnect interval must be greater than zero",
		},
		{
			name:      "zero auto commit time",
			configure: WithAutoCommitTime(0),
			wantError: "auto commit time must be greater than zero",
		},
		{
			name:      "negative auto commit time",
			configure: WithAutoCommitTime(-time.Second),
			wantError: "auto commit time must be greater than zero",
		},
		{
			name:      "zero consumer parallelism",
			configure: WithConsumerParallelism(0),
			wantError: "consumer parallelism must be greater than zero",
		},
		{
			name:      "negative consumer parallelism",
			configure: WithConsumerParallelism(-1),
			wantError: "consumer parallelism must be greater than zero",
		},
	}

	if strconv.IntSize > 32 {
		overflow := maxInt32 + 1
		tests = append(tests, struct {
			name      string
			configure func(*options)
			wantError string
		}{
			name:      "consumer parallelism exceeds int32",
			configure: WithConsumerParallelism(int(overflow)),
			wantError: "consumer parallelism exceeds int32",
		})
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ops := NewOptions(
				WithHost("localhost"),
				WithPort(50051),
				WithQueueName("queue-1"),
				test.configure,
			)
			err := ops.Validate()
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("Validate() error = %v, want error containing %q", err, test.wantError)
			}
		})
	}
}
