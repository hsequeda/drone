package http

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (h *DroneController) saveFile(src io.Reader) (filename string, err error) {
	if err = os.MkdirAll(h.uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("create %q dir: %w", h.uploadDir, err)
	}

	randomName := strconv.FormatInt(time.Now().UnixNano(), 10)
	filename = filepath.Join(h.uploadDir, randomName)
	dst, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("create new file: %w", err)
	}

	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("fill new file with data: %w", err)
	}

	return filename, nil
}
