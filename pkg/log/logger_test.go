// Copyright 2025 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf, Debug)

	logger.Info("test message", "key", "value", "number", 42)

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected log output to contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, "key=value") {
		t.Errorf("Expected log output to contain 'key=value', got: %s", output)
	}
	if !strings.Contains(output, "number=42") {
		t.Errorf("Expected log output to contain 'number=42', got: %s", output)
	}
}

func TestLogLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf, Warn)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Error("Debug message should not appear with Warn level")
	}
	if strings.Contains(output, "info message") {
		t.Error("Info message should not appear with Warn level")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("Warn message should appear with Warn level")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Error message should appear with Warn level")
	}
}

func TestJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithFormat(&buf, Info, FormatJSON)

	logger.Info("test message", "key", "value")

	output := buf.String()
	if !strings.Contains(output, `"msg":"test message"`) {
		t.Errorf("Expected JSON log output to contain '\"msg\":\"test message\"', got: %s", output)
	}
	if !strings.Contains(output, `"key":"value"`) {
		t.Errorf("Expected JSON log output to contain '\"key\":\"value\"', got: %s", output)
	}
	if !strings.Contains(output, `"level":"INFO"`) {
		t.Errorf("Expected JSON log output to contain '\"level\":\"INFO\"', got: %s", output)
	}
}

func TestWith(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf, Info).With("component", "test")

	logger.Info("test message", "key", "value")

	output := buf.String()
	if !strings.Contains(output, "component=test") {
		t.Errorf("Expected log output to contain 'component=test', got: %s", output)
	}
	if !strings.Contains(output, "key=value") {
		t.Errorf("Expected log output to contain 'key=value', got: %s", output)
	}
}

func TestWithScope(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf, Info).WithScope("lsp/pkg/notedownls")

	logger.Info("test message", "key", "value")

	output := buf.String()
	if !strings.Contains(output, "scope=lsp/pkg/notedownls") {
		t.Errorf("Expected log output to contain 'scope=lsp/pkg/notedownls', got: %s", output)
	}
	if !strings.Contains(output, "key=value") {
		t.Errorf("Expected log output to contain 'key=value', got: %s", output)
	}
}
