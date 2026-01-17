// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package logcmd

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/cilium/little-vm-helper/pkg/slogger"
	"github.com/stretchr/testify/assert"
)

func TestLogcmd(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "echo FOO>&2; echo LALA")

	infoRecorder := slogger.NewRecordingHandler(slogger.LevelInfo)
	warnRecorder := slogger.NewRecordingHandler(slogger.LevelWarn)
	multiHandler := slogger.NewMultiHandler(infoRecorder, warnRecorder)
	log := slogger.NewWithHandler(multiHandler)

	logStdout := getLogfForLevel(log, slogger.LevelInfo)
	logStderr := getLogfForLevel(log, slogger.LevelWarn)
	err := runAndLogCommand(nil, cmd, logStdout, logStderr)
	assert.Nil(t, err)
	assert.Equal(t, []string{"stderr> FOO\n"}, warnRecorder.Messages())
	assert.Equal(t, []string{"stdout> LALA\n"}, infoRecorder.Messages())
}

func TestLogcmdFail(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "exit 1")

	infoRecorder := slogger.NewRecordingHandler(slogger.LevelInfo)
	warnRecorder := slogger.NewRecordingHandler(slogger.LevelWarn)
	multiHandler := slogger.NewMultiHandler(infoRecorder, warnRecorder)
	log := slogger.NewWithHandler(multiHandler)

	logStdout := getLogfForLevel(log, slogger.LevelInfo)
	logStderr := getLogfForLevel(log, slogger.LevelWarn)
	err := runAndLogCommand(nil, cmd, logStdout, logStderr)
	assert.Error(t, err)
	assert.Equal(t, 0, len(warnRecorder.Messages()))
	assert.Equal(t, 0, len(infoRecorder.Messages()))
}

func TestLogcmdTimeout(t *testing.T) {
	infoRecorder := slogger.NewRecordingHandler(slogger.LevelInfo)
	warnRecorder := slogger.NewRecordingHandler(slogger.LevelWarn)
	multiHandler := slogger.NewMultiHandler(infoRecorder, warnRecorder)
	log := slogger.NewWithHandler(multiHandler)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	err := RunAndLogCommandContext(ctx, log, "/bin/sh", "-c", "sleep inf")
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestLogcmdNoTimeout(t *testing.T) {
	infoRecorder := slogger.NewRecordingHandler(slogger.LevelInfo)
	warnRecorder := slogger.NewRecordingHandler(slogger.LevelWarn)
	multiHandler := slogger.NewMultiHandler(infoRecorder, warnRecorder)
	log := slogger.NewWithHandler(multiHandler)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := RunAndLogCommandContext(ctx, log, "/bin/sh", "-c", "sleep .1s")
	assert.Nil(t, err)
}

func TestRunAndLogCommandsContext(t *testing.T) {
	infoRecorder := slogger.NewRecordingHandler(slogger.LevelInfo)
	warnRecorder := slogger.NewRecordingHandler(slogger.LevelWarn)
	multiHandler := slogger.NewMultiHandler(infoRecorder, warnRecorder)
	log := slogger.NewWithHandler(multiHandler)

	err := RunAndLogCommandsContext(context.Background(), log,
		[]string{"/bin/sh", "-c", "echo FOO"},
		[]string{"/bin/sh", "-c", "echo LALA>&2"},
	)
	assert.Nil(t, err)
	assert.Equal(t, []string{"starting command", "stdout> FOO\n", "starting command"}, infoRecorder.Messages())
	assert.Equal(t, []string{"stderr> LALA\n"}, warnRecorder.Messages())
}
