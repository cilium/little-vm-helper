// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package slogger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerBasic(t *testing.T) {
	infoHandler := NewRecordingHandlerAtOrAbove(LevelInfo)
	warnHandler := NewRecordingHandlerAtOrAbove(LevelWarn)
	multiHandler := NewMultiHandler(infoHandler, warnHandler)

	log := NewWithHandler(multiHandler)

	log.Info("info message")
	log.Warn("warn message")
	log.Error("error message")

	infoMsgs := infoHandler.Messages()
	warnMsgs := warnHandler.Messages()

	assert.Contains(t, infoMsgs, "info message")
	assert.Contains(t, infoMsgs, "warn message")
	assert.Contains(t, infoMsgs, "error message")

	assert.NotContains(t, warnMsgs, "info message")
	assert.Contains(t, warnMsgs, "warn message")
	assert.Contains(t, warnMsgs, "error message")
}

func TestLoggerFormatted(t *testing.T) {
	handler := NewRecordingHandlerAtOrAbove(LevelInfo)
	log := NewWithHandler(handler)

	log.Infof("hello %s", "world")
	log.Warnf("count: %d", 42)

	msgs := handler.Messages()
	assert.Contains(t, msgs, "hello world")
	assert.Contains(t, msgs, "count: 42")
}

func TestLoggerWithField(t *testing.T) {
	handler := NewRecordingHandlerAtOrAbove(LevelInfo)
	log := NewWithHandler(handler)

	log.WithField("key", "value").Info("message with field")

	records := handler.Records()
	assert.Equal(t, 1, len(records))

	// Check that the attribute was recorded
	var foundAttr bool
	records[0].Attrs(func(attr slog.Attr) bool {
		if attr.Key == "key" && attr.Value.String() == "value" {
			foundAttr = true
		}
		return true
	})
	assert.True(t, foundAttr, "expected to find key=value attribute")
	assert.Equal(t, "message with field", records[0].Message)
}

func TestLoggerWithFields(t *testing.T) {
	handler := NewRecordingHandlerAtOrAbove(LevelInfo)
	log := NewWithHandler(handler)

	fields := map[string]any{
		"key1": "value1",
		"key2": 123,
	}
	log.WithFields(fields).Info("message with fields")

	records := handler.Records()
	assert.Equal(t, 1, len(records))
	assert.Equal(t, "message with fields", records[0].Message)
}

func TestLoggerWithError(t *testing.T) {
	handler := NewRecordingHandlerAtOrAbove(LevelInfo)
	log := NewWithHandler(handler)

	err := assert.AnError
	log.WithError(err).Error("something failed")

	records := handler.Records()
	assert.Equal(t, 1, len(records))
	assert.Equal(t, "something failed", records[0].Message)
}

func TestLoggerChaining(t *testing.T) {
	handler := NewRecordingHandlerAtOrAbove(LevelInfo)
	log := NewWithHandler(handler)

	log.WithField("request_id", "123").
		WithField("user", "alice").
		WithError(assert.AnError).
		Info("chained fields")

	records := handler.Records()
	assert.Equal(t, 1, len(records))
	assert.Equal(t, "chained fields", records[0].Message)
}

func TestDiscardLogger(t *testing.T) {
	log := NewDiscard()

	// Should not panic
	log.Info("discarded")
	log.Warn("also discarded")
	log.Error("discarded too")
	log.WithField("key", "value").Info("discarded with field")
}

func TestRecordingHandlerClear(t *testing.T) {
	handler := NewRecordingHandlerAtOrAbove(LevelInfo)
	log := NewWithHandler(handler)

	log.Info("message 1")
	log.Info("message 2")

	assert.Equal(t, 2, len(handler.Messages()))

	handler.Clear()
	assert.Equal(t, 0, len(handler.Messages()))

	log.Info("message 3")
	assert.Equal(t, 1, len(handler.Messages()))
}

func TestLoggerLevels(t *testing.T) {
	tests := []struct {
		name      string
		level     Level
		logFunc   func(Logger)
		shouldLog bool
	}{
		{"debug at info level", LevelInfo, func(l Logger) { l.Debug("test") }, false},
		{"info at info level", LevelInfo, func(l Logger) { l.Info("test") }, true},
		{"warn at info level", LevelInfo, func(l Logger) { l.Warn("test") }, true},
		{"error at info level", LevelInfo, func(l Logger) { l.Error("test") }, true},
		{"debug at debug level", LevelDebug, func(l Logger) { l.Debug("test") }, true},
		{"info at warn level", LevelWarn, func(l Logger) { l.Info("test") }, false},
		{"warn at warn level", LevelWarn, func(l Logger) { l.Warn("test") }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewRecordingHandlerAtOrAbove(tt.level)
			log := NewWithHandler(handler)
			tt.logFunc(log)

			if tt.shouldLog {
				assert.Equal(t, 1, len(handler.Messages()), "expected message to be logged")
			} else {
				assert.Equal(t, 0, len(handler.Messages()), "expected message NOT to be logged")
			}
		})
	}
}
