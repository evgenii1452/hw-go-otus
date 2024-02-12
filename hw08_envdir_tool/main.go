package main

import (
	"errors"
	"os"
)

var ErrToFewArguments = errors.New("to few arguments")

func main() {
	if len(os.Args) < 3 {
		panic(ErrToFewArguments)
	}

	envs, err := ReadDir(os.Args[1])
	if err != nil {
		panic(err)
	}

	RunCmd(os.Args[2:], envs)
}
