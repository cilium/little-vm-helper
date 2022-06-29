package logcmd

import (
	"io/ioutil"
	"os/exec"
	"testing"

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

	err := runAndLogCommand(cmd, log, logrus.InfoLevel, logrus.WarnLevel)
	assert.Nil(t, err)
	assert.Equal(t, []string{"stderr> FOO\n"}, warnRecorder.Messages)
	assert.Equal(t, []string{"starting command", "stdout> LALA\n"}, infoRecorder.Messages)
}
