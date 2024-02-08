package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testFilesDir = "./testdata/"
	testFilePath = testFilesDir + "input.txt"
)

func TestCopy(t *testing.T) {
	createFile := func(content []byte) *os.File {
		file, err := os.CreateTemp("./", "")
		if err != nil {
			panic("Temporary file creating error")
		}
		defer file.Close()

		_, err = file.Write(content)
		if err != nil {
			panic("Error writing to temporary file")
		}

		return file
	}

	copyFunction := func(content []byte, offset, limit int64) []byte {
		fileCopy := "./copy.txt"
		file := createFile(content)
		defer os.Remove(file.Name())
		defer os.Remove(fileCopy)

		err := Copy(file.Name(), fileCopy, offset, limit)
		if err != nil {
			panic(err)
		}

		copiedContent, err := os.ReadFile(fileCopy)
		if err != nil {
			panic("Failed reading file copy")
		}

		return copiedContent
	}

	t.Run("success copy", func(t *testing.T) {
		result := copyFunction([]byte("test file"), 0, 0)
		require.Equal(t, []byte("test file"), result)
	})

	t.Run("success copy with offset", func(t *testing.T) {
		result := copyFunction([]byte("test file"), 1, 0)
		require.Equal(t, []byte("est file"), result)
	})

	t.Run("success copy with limit", func(t *testing.T) {
		result := copyFunction([]byte("test file"), 0, 5)
		require.Equal(t, []byte("test "), result)
	})

	t.Run("success copy with offset and limit", func(t *testing.T) {
		result := copyFunction([]byte("test file"), 2, 5)
		require.Equal(t, []byte("st fi"), result)
	})

	t.Run("success copy if limit grater file size", func(t *testing.T) {
		result := copyFunction([]byte("test file"), 2, 100)
		require.Equal(t, []byte("st file"), result)
	})
}

func TestFailCopy(t *testing.T) {
	t.Run("error if offset is negative", func(t *testing.T) {
		err := Copy(testFilePath, "test.txt", -5, 10)
		require.Truef(t, errors.Is(err, ErrOffsetCannotBeNegative), "actual error - %q", err)
	})

	t.Run("error if limit is negative", func(t *testing.T) {
		err := Copy(testFilePath, "test.txt", 10, -5)
		require.Truef(t, errors.Is(err, ErrLimitCannotBeNegative), "actual error - %q", err)
	})

	t.Run("error if input path is dir", func(t *testing.T) {
		err := Copy(testFilesDir, "test.txt", 10, 10)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error - %q", err)
	})

	t.Run("error if offset is grater than file size", func(t *testing.T) {
		file, err := os.OpenFile(testFilePath, os.O_RDONLY, 0o666)
		if err != nil {
			panic("Test file not found!")
		}
		fileInfo, _ := file.Stat()

		err = Copy(testFilePath, "test.txt", fileInfo.Size()+1, 10)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error - %q", err)
	})

	t.Run("error if offset is grater than file size", func(t *testing.T) {
		file, err := os.OpenFile(testFilePath, os.O_RDONLY, 0o666)
		if err != nil {
			panic("Test file not found!")
		}
		fileInfo, _ := file.Stat()

		err = Copy(testFilePath, "test.txt", fileInfo.Size()+1, 10)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error - %q", err)
	})
}
