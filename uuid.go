package goutils

import (
	"strings"

	"github.com/google/uuid"
)

func UUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	code := strings.ReplaceAll(id.String(), "-", "")
	return code, nil
}
