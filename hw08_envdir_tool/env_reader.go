package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func makeEnvValue(value string) EnvValue {
	return EnvValue{
		Value:      value,
		NeedRemove: false,
	}
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	env := make(Environment)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		envKey := file.Name()

		value, err := readFirstLine(fmt.Sprintf("%s/%s", dir, file.Name()))
		if err != nil {
			return nil, err
		}

		envVal := makeEnvValue(value)
		if strings.TrimSpace(value) == "" {
			envVal.NeedRemove = true
		}

		env[envKey] = envVal
	}

	return env, nil
}

func readFirstLine(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	line, _, err := reader.ReadLine()
	if err != nil && !errors.Is(io.EOF, err) {
		return "", fmt.Errorf("read line: %w", err)
	}

	line = bytes.ReplaceAll(line, []byte{0}, []byte{10})

	return string(line), nil
}
