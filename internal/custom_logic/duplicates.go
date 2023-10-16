package custom_logic

import (
	"errors"
	"regexp"
)

var pattern = regexp.MustCompile(`$.([0-9])^`)

func getDuplicateKey(message string) (string, error) {

	matches := pattern.FindStringSubmatch(message)
	if len(matches) < 1 {
		return "", errors.New("cannot find a unique key")
	}
	return matches[0], nil

}
