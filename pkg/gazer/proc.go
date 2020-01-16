/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"

	"github.com/mattn/go-shellwords"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/time"
)

func executeCommandOrTimeout(cmd *exec.Cmd, timeoutMill int) error {
	exec := executeCommandAsync(cmd)
	if timeoutMill <= 0 {
		e := <-exec
		return e
	}

	timeout := time.After(timeoutMill)
	var err error
	finished := false
	for {
		if finished {
			break
		}
		select {
		case <-timeout:
			if cmd.Process == nil {
				timeout = time.After(5)
				break
			}
			kill(cmd, "Timeout")
			finished = true
		case err = <-exec:
			finished = true
		}
	}
	return err
}

func executeCommandAsync(cmd *exec.Cmd) <-chan error {
	ch := make(chan error)

	go func() {
		if cmd == nil {
			ch <- errors.New("Failed: cmd is nil")
			return
		}
		err := executeCommand(cmd)
		if err != nil {
			ch <- err
			return
		}
		ch <- nil
	}()
	return ch
}

func executeCommand(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	if cmd.Process != nil {
		logger.Info("Pid: %d", cmd.Process.Pid)
	} else {
		logger.Info("Pid: ????")
	}
	err := cmd.Wait()
	return err
}

func kill(cmd *exec.Cmd, reason string) {
	err := cmd.Process.Signal(os.Interrupt)
	if err != nil {
		logger.NoticeObject(err)
	}
	logger.Notice("%s: %d has been killed", reason, cmd.Process.Pid)
}

func createCommand(commandString string) *exec.Cmd {
	perser := shellwords.NewParser()
	// perser.ParseBacktick = true
	// perser.ParseEnv = true
	args, err := perser.Parse(commandString)
	if err != nil {
		return nil
	}
	if len(args) == 1 {
		return exec.Command(args[0])
	}
	return exec.Command(args[0], args[1:]...)
}

func sigIntChannel() chan struct{} {
	ch := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		close(ch)
	}()
	return ch
}
