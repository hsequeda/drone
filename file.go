package drone

import (
	"fmt"
	"io"
	"os"
	"time"
)

func saveFile(src io.Reader) (filename string, err error) {
	if err = os.MkdirAll("./uploads", os.ModePerm); err != nil {
		return "", fmt.Errorf("create './uploads' dir: %w", err)
	}

	filename = fmt.Sprintf("./uploads/%d", time.Now().UnixNano())
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
