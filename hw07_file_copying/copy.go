package main

import (
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrFromOrToPathIsEmpty        = errors.New("from and to paths should be don't empty")
	ErrLimitOrOffsetIsNotPositive = errors.New("limit and offset should be positive numbers")
	ErrExcessOffsetFileSize       = errors.New("excess offset")
	ErrUndefinedFileSize          = errors.New("unsupported file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := validateParams(fromPath, toPath, offset, limit); err != nil {
		return err
	}

	fileFrom, err := os.OpenFile(fromPath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	defer func() {
		err := fileFrom.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	fileInfo, err := fileFrom.Stat()
	if err != nil {
		return err
	}

	if err := validateFileMetadata(fileInfo, offset); err != nil {
		return err
	}

	_, err = fileFrom.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer func() {
		err := fileTo.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if limit == 0 {
		limit = fileInfo.Size()
	}

	progressBarCount := fileInfo.Size() - offset
	if limit < progressBarCount {
		progressBarCount = limit
	}

	progressBar := prepareProgressBar(progressBarCount)
	progressBar.Start()

	barReader := progressBar.NewProxyReader(fileFrom)
	defer progressBar.Finish()

	_, err = io.CopyN(fileTo, barReader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

func prepareProgressBar(count int64) *pb.ProgressBar {
	progressBar := pb.New64(count)
	progressBar.SetRefreshRate(time.Millisecond)
	progressBar.Set(pb.Bytes, true)

	return progressBar
}

func validateParams(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrFromOrToPathIsEmpty
	}

	if limit < 0 || offset < 0 {
		return ErrLimitOrOffsetIsNotPositive
	}

	return nil
}

func validateFileMetadata(fileInfo os.FileInfo, offset int64) error {
	fileSize := fileInfo.Size()

	if offset > fileSize {
		return ErrExcessOffsetFileSize
	}

	if fileSize == 0. {
		return ErrUndefinedFileSize
	}

	return nil
}
