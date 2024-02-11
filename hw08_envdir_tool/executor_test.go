package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	initEnvs := func(envs map[string]string) Environment {
		result := make(Environment)

		for k, v := range envs {
			env := EnvValue{
				Value:      v,
				NeedRemove: false,
			}
			if v == "" {
				env.NeedRemove = true
			}

			result[k] = env
		}

		return result
	}

	createLogFile := func() *os.File {
		logFile, err := os.Create("./logs")
		if err != nil {
			panic("Failed to create log file")
		}

		return logFile
	}

	getLog := func() string {
		logFile, err := os.Open("./logs")
		if err != nil {
			panic("Failed to open log file")
		}

		reader := bufio.NewReader(logFile)
		line, _, err := reader.ReadLine()
		if err != nil {
			panic("Failed to read log file")
		}

		return string(line)
	}

	t.Run("test success", func(t *testing.T) {
		basicStdout := os.Stdout
		logFile := createLogFile()
		defer os.Remove(logFile.Name())

		envs := initEnvs(map[string]string{"TEST_ENV": "TEST"})

		os.Stdout = logFile
		result := RunCmd([]string{"/bin/bash", "./testdata/checkVarExists.sh", "TEST_ENV"}, envs)
		os.Stdout = basicStdout

		require.Equal(t, 0, result)
		require.Equal(t, "TEST", getLog())
	})

	t.Run("test success variable not exists", func(t *testing.T) {
		basicStdout := os.Stdout
		logFile := createLogFile()
		defer os.Remove(logFile.Name())

		os.Stdout = logFile
		result := RunCmd([]string{"/bin/bash", "./testdata/checkVarExists.sh", "SOME_ENV"}, make(Environment))
		os.Stdout = basicStdout

		require.Equal(t, 0, result)
		require.Equal(t, "Variable not exists", getLog())
	})

	t.Run("test fail to few arguments", func(t *testing.T) {
		basicStdout := os.Stdout
		logFile := createLogFile()
		defer os.Remove(logFile.Name())

		os.Stdout = logFile
		result := RunCmd([]string{"/bin/bash", "./testdata/checkVarExists.sh"}, make(Environment))
		os.Stdout = basicStdout

		require.Equal(t, 1, result)
		require.Equal(t, "To few arguments", getLog())
	})

	t.Run("test fail", func(t *testing.T) {
		basicStdout := os.Stdout
		logFile := createLogFile()
		defer os.Remove(logFile.Name())

		os.Stdout = logFile
		envs := initEnvs(make(map[string]string))
		result := RunCmd([]string{"/bin/echo"}, envs)
		os.Stdout = basicStdout

		require.Equal(t, 1, result)
	})
}
