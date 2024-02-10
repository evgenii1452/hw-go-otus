package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile        = errors.New("unsupported file")
	ErrOffsetExceedsFileSize  = errors.New("offset exceeds file size")
	ErrOffsetCannotBeNegative = errors.New("offset cannot be negative")
	ErrLimitCannotBeNegative  = errors.New("limit cannot be negative")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	sourceFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0o666)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	defer sourceFile.Close()

	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("get source file info: %w", err)
	}

	err = validate(fileInfo, offset, limit)
	if err != nil {
		return err
	}

	if limit == 0 || limit+offset > fileInfo.Size() {
		limit = fileInfo.Size() - offset
	}

	_, err = sourceFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("set offset in resource file: %w", err)
	}

	limitReader := io.LimitReader(sourceFile, limit)
	bar := pb.Full.Start64(limit)
	readerWithPB := bar.NewProxyReader(limitReader)

	err = copyFile(readerWithPB, toPath)
	if err != nil {
		return err
	}

	bar.Finish()

	return nil
}

func validate(fileInfo os.FileInfo, offset, limit int64) error {
	if offset < 0 {
		return ErrOffsetCannotBeNegative
	}

	if limit < 0 {
		return ErrLimitCannotBeNegative
	}

	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	if fileInfo.IsDir() {
		return ErrUnsupportedFile
	}

	return nil
}

func copyFile(resource io.Reader, copyPath string) error {
	file, err := os.Create(copyPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resource)
	if err != nil {
		defer os.Remove(copyPath)
		return fmt.Errorf("copy file: %w", err)
	}

	return nil
}
