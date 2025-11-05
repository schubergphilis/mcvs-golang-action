package data

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type Storer interface {
	DownloadFileFromURL(url string) ([]byte, error)
	WriteToFile(dest string, data []byte) error
	WriteTextFile(dest string, content string) error
}

type FileStorer struct {
	Logger *slog.Logger
}

func (fs *FileStorer) DownloadFileFromURL(url string) ([]byte, error) {
	fs.Logger.Info("Requesting file", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		fs.Logger.Error("HTTP GET failed", "err", err)

		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fs.Logger.Warn("Non-200 response", "status", resp.Status)

		return nil, errors.New("failed to download file: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fs.Logger.Error("Failed to read response body", "err", err)

		return nil, err
	}

	fs.Logger.Info("File downloaded", "bytes", len(body))

	return body, nil
}

func (fs *FileStorer) WriteToFile(dest string, data []byte) error {
	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fs.Logger.Error("Failed to create directory", "err", err, "dir", dir)

		return err
	}

	err := os.WriteFile(dest, data, 0o644)
	if err != nil {
		fs.Logger.Error("Failed to write file", "err", err, "dest", dest)

		return err
	}

	fs.Logger.Info("File written", "dest", dest, "size", len(data))

	return nil
}

func (fs *FileStorer) WriteTextFile(dest string, content string) error {
	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fs.Logger.Error("Failed to create directory", "err", err, "dir", dir)

		return err
	}

	err := os.WriteFile(dest, []byte(content), 0o644)
	if err != nil {
		fs.Logger.Error("Failed to write text file", "err", err, "dest", dest)

		return err
	}

	fs.Logger.Info("Text file written", "dest", dest, "size", len(content))

	return nil
}
