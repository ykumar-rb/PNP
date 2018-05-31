package executor

import (
	"os/exec"
	"bytes"
	"io"
	"os"
	"fmt"
	"strings"
)

func ExecuteCommand (cmdString string) (errStr string, err error){
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("bash","-c", cmdString)

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err = cmd.Start()
	if err != nil {
		fmt.Printf("cmd.Start() failed with '%s'\n", err)
		return "", err
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}

	errStr = string(stderrBuf.Bytes())
	return errStr, err
}

func ExecuteServerInstructions(cmdString []string) (exeErr error) {
	var errStr string
	cmd := strings.Join(cmdString, " && ")
	errStr, exeErr = ExecuteCommand(cmd)
	if exeErr != nil {
		fmt.Printf("\nCommand <%v> failed to execute\nErrorString: %v\nError: %v\n", cmd, errStr, exeErr)
	}
	return exeErr
}

