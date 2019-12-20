package terraformcli

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/kernel164/awspolicy/internal/capture"
	"github.com/kernel164/awspolicy/internal/util"
)

type terraformCli struct {
	logPath string
}

func New() capture.Capture {
	logPath := path.Join(util.Getenv("TERRAFORM_CLI_LOG_DIR", "/tmp"), fmt.Sprintf("terraform_%v.log", time.Now().UnixNano()))
	return &terraformCli{logPath: logPath}
}

func (t *terraformCli) Run(params *capture.Params) (*capture.Results, error) {
	// remove log after use.
	defer os.Remove(t.logPath)

	// Disable output buffering, enable streaming
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}

	// environ
	envs := os.Environ()
	envs = append(envs, "TF_LOG=TRACE")
	envs = append(envs, "TF_LOG_PATH="+t.logPath)

	// Create Cmd with options
	// change args to add -no-color
	tfCmd := cmd.NewCmdOptions(cmdOptions, util.Getenv("TERRAFORM_CLI_PATH", "terraform"), params.Args...)
	tfCmd.Env = envs
	tfCmd.Dir = util.Getenv("TERRAFORM_CLI_CD_DIR", "")

	// Run and wait for Cmd to return, discard Status
	go func() {
		for {
			select {
			// Print STDOUT and STDERR lines streaming from Cmd
			case line := <-tfCmd.Stdout:
				fmt.Println(line)
			case line := <-tfCmd.Stderr:
				fmt.Fprintln(os.Stderr, line)
			}
		}
	}()

	status := <-tfCmd.Start()
	time.Sleep(100 * time.Millisecond) // time for goroutine to flush cmd's stdout/stderr
	if status.Error != nil {
		return nil, status.Error
	}

	return t.processLogs()
}

func (t *terraformCli) processLogs() (*capture.Results, error) {
	actions := map[string]struct{}{}
	data := &capture.Results{Actions: actions}

	// open log file
	file, err := os.Open(t.logPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// scan log file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "[aws-sdk-go] DEBUG: Request") {
			found := false
			for _, part := range strings.Split(line, " ") {
				if found {
					key := strings.Replace(strings.TrimSpace(part), "/", ":", -1)
					actions[key] = struct{}{}
					break
				}
				if part == "Request" {
					found = true
				}
			}
		}
	}

	// scan error
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// return results.
	return data, nil
}
