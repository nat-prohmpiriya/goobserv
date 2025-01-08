package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/nat-prohmpiriya/goobserv/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestStdoutOutput(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Restore stdout when done
	defer func() {
		os.Stdout = oldStdout
	}()

	output := NewStdoutOutput()

	// Create test entries
	now := time.Now()
	entries := []*core.Entry{
		{
			Time:      now,
			Level:     core.LevelInfo,
			Message:   "Test message 1",
			TraceID:   "trace-1",
			SpanID:    "span-1",
			RequestID: "req-1",
			Data: map[string]interface{}{
				"key1": "value1",
			},
		},
		{
			Time:      now,
			Level:     core.LevelError,
			Message:   "Test message 2",
			TraceID:   "trace-2",
			SpanID:    "span-2",
			RequestID: "req-2",
			Data: map[string]interface{}{
				"key2": "value2",
			},
		},
	}

	// Write entries
	err := output.Write(entries)
	assert.NoError(t, err, "Write should not return error")

	// Close pipe
	w.Close()

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Split output into lines
	lines := bytes.Split(buf.Bytes(), []byte("\n"))
	assert.Equal(t, 3, len(lines), "Should have 2 lines plus empty line")

	// Parse and verify each line
	for i, line := range lines[:2] {
		var logEntry map[string]interface{}
		err := json.Unmarshal(line, &logEntry)
		assert.NoError(t, err, "Should be valid JSON")

		// Verify fields
		assert.Equal(t, entries[i].Time.Format(time.RFC3339), logEntry["timestamp"])
		assert.Equal(t, entries[i].Level.String(), logEntry["level"])
		assert.Equal(t, entries[i].Message, logEntry["message"])

		data := logEntry["data"].(map[string]interface{})
		assert.Equal(t, entries[i].TraceID, data["trace_id"])
		assert.Equal(t, entries[i].SpanID, data["span_id"])
		assert.Equal(t, entries[i].RequestID, data["request_id"])
		for k, v := range entries[i].Data {
			assert.Equal(t, v, data[k])
		}
	}

	// Test Flush and Close
	assert.NoError(t, output.Flush(), "Flush should not return error")
	assert.NoError(t, output.Close(), "Close should not return error")
}
