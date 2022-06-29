package logcmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
)

func logReader(rd *bufio.Reader, log *logrus.Logger, prefix string, level logrus.Level) error {
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		log.Logf(level, "%s%s", prefix, line)
	}
}

func runAndLogCommand(
	cmd *exec.Cmd,
	log *logrus.Logger,
	stdoutLevel, stderrLevel logrus.Level,
) error {

	// prepare pipes
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("StdErrPipe() failed: %w", err)
	}
	defer stderr.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("StdOutPipe() failed: %w", err)
	}
	defer stdout.Close()

	// start command
	log.WithField("path", cmd.Path).WithField("args", cmd.Args).Info("starting command")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	// start logging
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		rd := bufio.NewReader(stdout)
		err = logReader(rd, log, "stdout> ", stdoutLevel)
		if err != nil {
			log.Warnf("failed to read from stdout: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		rd := bufio.NewReader(stderr)
		err = logReader(rd, log, "stderr> ", stderrLevel)
		if err != nil {
			log.Warnf("failed to read from stderr: %v", err)
		}
	}()

	// we need to wait for the pipes before waiting for the command
	// see: https://pkg.go.dev/os/exec#Cmd.StdoutPipe
	wg.Wait()

	return cmd.Wait()
}

func RunAndLogCommand(
	cmd *exec.Cmd,
	log *logrus.Logger,
) error {
	return runAndLogCommand(cmd, log, logrus.InfoLevel, logrus.WarnLevel)
}
