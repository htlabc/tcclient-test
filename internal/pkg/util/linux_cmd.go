package util

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

func RunLinuxCmd(cmdpath string, args ...string) (string, error) {
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	fmt.Println(fmt.Sprintf("Execute shell command: %s ,args: %s", cmdpath, args))
	cmd := exec.Command(cmdpath, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if outStr != "" {
		fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	}
	if err != nil || errStr != "" {
		return "", errors.New(errStr)
	}

	return outStr, nil

}

func RunRemoteShellcmd(ip, shellcmd string) (string, error) {
	return RunLinuxCmd("sshpass", "-p", "tc123456", "ssh", "-o", "StrictHostKeychecking=no", "-o", "ConnectTimeout=3", fmt.Sprintf("root@%s", ip), shellcmd)
}
