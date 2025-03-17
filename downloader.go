package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cheggaaa/pb/v3"
)

func Download(url string, refererURL string, outputPath string) error {
	totalFileSize, err := getTotalFileSize(url, refererURL)
	if err != nil {
		return err
	}

	tempFilePath := outputPath + ".download"
	var downloadedFileSize int64

	tempFileStat, err := os.Stat(tempFilePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		downloadedFileSize = tempFileStat.Size()
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", refererURL)
	if downloadedFileSize > 0 {
		req.Header.Set("Range", "bytes="+strconv.FormatInt(downloadedFileSize, 10)+"-")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http status code: %v, body: %s", resp.StatusCode, body)
	}

	tempFile, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	filename := filepath.Base(outputPath)
	processBarTempate := `{{counters .}} {{bar . "[" "=" ">" "-" "]"}} {{speed .}} {{percent . | green}} {{rtime .}}` + filename
	bar := pb.New64(totalFileSize).
		Set(pb.Bytes, true).
		SetMaxWidth(1000).
		SetTemplateString(processBarTempate).
		SetCurrent(downloadedFileSize)
	bar.Start()

	if _, err := tempFile.Seek(downloadedFileSize, 0); err != nil {
		return err
	}

	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := tempFile.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		bar.Add(n)
	}

	bar.Finish()

	if err := tempFile.Sync(); err != nil {
		return err
	}
	return os.Rename(tempFilePath, outputPath)
}

func getTotalFileSize(url string, refererURL string) (int64, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Referer", refererURL)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	contentLengthStr := resp.Header.Get("Content-Length")
	var totalSize int64
	if contentLengthStr != "" {
		totalSize, err = strconv.ParseInt(contentLengthStr, 10, 64)
		if err != nil {
			return 0, err
		}
	}
	return totalSize, nil
}
