package main

import (
	"fmt"
	"os"
	"os/exec"
)

const ErrorCode = 1

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 2 {
		fmt.Printf("Error: Too few atguments")
		return ErrorCode
	}
	cmdName := cmd[0]
	args := make([]string, len(cmd)-1)

	if len(cmd) > 1 {
		args = cmd[1:]
	}

	command := exec.Command(cmdName, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	for k, v := range env {
		if v.NeedRemove {
			err := os.Unsetenv(k)
			if err != nil {
				fmt.Printf("Failed unset variable %s: %s \n", k, err.Error())
				return ErrorCode
			}

			continue
		}

		err := os.Setenv(k, v.Value)
		if err != nil {
			fmt.Printf("Failed set valiable %s: %s \n", k, err.Error())
			return ErrorCode
		}
	}

	err := command.Run()
	if err != nil {
		fmt.Println("Failed run command: %w", err)
		return ErrorCode
	}

	return command.ProcessState.ExitCode()
}
