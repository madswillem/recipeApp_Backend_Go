package tools

import (
	"bytes"
	"io"
	"os"
)

func ReadFileAsString(filePath string) (string, error) {
	println(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
