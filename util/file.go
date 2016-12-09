package util

import (
	"io/ioutil"
	"os"
)

func ReadLatin1File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	buffer := make([]rune, len(content))
	for i, b := range content {
		buffer[i] = rune(b)
	}

	return string(buffer), nil
}
