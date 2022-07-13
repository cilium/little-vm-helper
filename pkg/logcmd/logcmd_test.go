package logcmd

import (
	"context"
	"io/ioutil"
	"os/exec"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type logrusRecorder struct {
	Level    logrus.Level
	Messages []string
}

func (r *logrusRecorder) Levels() []logrus.Level {
	return []logrus.Level{r.Level}
}

func (r *logrusRecorder) Fire(e *logrus.Entry) error {
	r.Messages = append(r.Messages, e.Message)
	return nil
}

func TestLogcmd(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "echo FOO>&2; echo LALA")
	log := logrus.New()
	log.SetOutput(ioutil.Discard)

	infoRecorder := logrusRecorder{Level: logrus.InfoLevel}
	warnRecorder := logrusRecorder{Level: logrus.WarnLevel}
	log.AddHook(&infoRecorder)
	log.AddHook(&warnRecorder)

	err := runAndLogCommand(nil, cmd, log, logrus.InfoLevel, logrus.WarnLevel)
	assert.Nil(t, err)
	assert.Equal(t, []string{"stderr> FOO\n"}, warnRecorder.Messages)
	assert.Equal(t, []string{"starting command", "stdout> LALA\n"}, infoRecorder.Messages)
}

func TestLogcmdFail(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "exit 1")
	log := logrus.New()
	log.SetOutput(ioutil.Discard)

	infoRecorder := logrusRecorder{Level: logrus.InfoLevel}
	warnRecorder := logrusRecorder{Level: logrus.WarnLevel}
	log.AddHook(&infoRecorder)
	log.AddHook(&warnRecorder)

	err := runAndLogCommand(nil, cmd, log, logrus.InfoLevel, logrus.WarnLevel)
	assert.Error(t, err)
	assert.Nil(t, warnRecorder.Messages)
	assert.Equal(t, []string{"starting command"}, infoRecorder.Messages)
}

func TestLogcmdTimeout(t *testing.T) {
	log := logrus.New()
	log.SetOutput(ioutil.Discard)

	infoRecorder := logrusRecorder{Level: logrus.InfoLevel}
	warnRecorder := logrusRecorder{Level: logrus.WarnLevel}
	log.AddHook(&infoRecorder)
	log.AddHook(&warnRecorder)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	err := RunAndLogCommandContext(ctx, log, "/bin/sh", "-c", "sleep inf")
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestLogcmdNoTimeout(t *testing.T) {
	log := logrus.New()
	log.SetOutput(ioutil.Discard)
	infoRecorder := logrusRecorder{Level: logrus.InfoLevel}
	warnRecorder := logrusRecorder{Level: logrus.WarnLevel}
	log.AddHook(&infoRecorder)
	log.AddHook(&warnRecorder)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := RunAndLogCommandContext(ctx, log, "/bin/sh", "-c", "sleep .1s")
	assert.Nil(t, err)
}

func TestRunAndLogCommandsContext(t *testing.T) {
	log := logrus.New()
	log.SetOutput(ioutil.Discard)
	infoRecorder := logrusRecorder{Level: logrus.InfoLevel}
	warnRecorder := logrusRecorder{Level: logrus.WarnLevel}
	log.AddHook(&infoRecorder)
	log.AddHook(&warnRecorder)

	err := RunAndLogCommandsContext(context.Background(), log,
		[]string{"/bin/sh", "-c", "echo FOO"},
		[]string{"/bin/sh", "-c", "echo LALA>&2"},
	)
	assert.Nil(t, err)
	assert.Equal(t, []string{"starting command", "stdout> FOO\n", "starting command"}, infoRecorder.Messages)
	assert.Equal(t, []string{"stderr> LALA\n"}, warnRecorder.Messages)
}
