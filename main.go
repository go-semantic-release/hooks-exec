package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"

	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/mattn/go-shellwords"
)

var version = "dev"

type Exec struct {
	cmdOnSuccess   *template.Template
	cmdOnNoRelease *template.Template
	log            *log.Logger
}

func (e *Exec) Init(m map[string]string) error {
	var err error
	e.cmdOnSuccess, err = template.New("command").Parse(m["exec_on_success"])
	if err != nil {
		return fmt.Errorf("failed to parse exec_on_success: %w", err)
	}
	e.cmdOnNoRelease, err = template.New("command").Parse(m["exec_on_no_release"])
	if err != nil {
		return fmt.Errorf("failed to parse exec_on_no_release: %w", err)
	}
	return nil
}

func (e *Exec) Name() string {
	return "exec"
}

func (e *Exec) Version() string {
	return version
}

func (e *Exec) runCommand(cmdStr string) error {
	e.log.Println("#", cmdStr)
	cmdArgs, err := shellwords.Parse(cmdStr)
	if err != nil {
		return fmt.Errorf("failed to parse command: %w", err)
	}
	if len(cmdArgs) == 0 {
		e.log.Println("warning: command is empty")
		return nil
	}
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = os.Environ()
	logPipeRead, logPipeWrite := io.Pipe()
	defer logPipeWrite.Close()
	cmd.Stdout = logPipeWrite
	cmd.Stderr = logPipeWrite
	logLineScanner := bufio.NewScanner(logPipeRead)
	go func() {
		for logLineScanner.Scan() {
			e.log.Println(logLineScanner.Text())
		}
	}()

	return cmd.Run()
}

func executeTemplate(tpl *template.Template, data any) (string, error) {
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute command template: %w", err)
	}
	return buf.String(), nil
}

func (e *Exec) Success(config *hooks.SuccessHookConfig) error {
	cmdStr, err := executeTemplate(e.cmdOnSuccess, config)
	if err != nil {
		return err
	}
	if cmdStr == "" {
		return nil
	}
	e.log.Println("executing success command:")
	return e.runCommand(cmdStr)
}

func (e *Exec) NoRelease(config *hooks.NoReleaseConfig) error {
	cmdStr, err := executeTemplate(e.cmdOnNoRelease, config)
	if err != nil {
		return err
	}
	if cmdStr == "" {
		return nil
	}
	e.log.Println("executing no release command:")
	return e.runCommand(cmdStr)
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Hooks: func() hooks.Hooks {
			return &Exec{
				log: log.New(os.Stderr, "", 0),
			}
		},
	})
}
