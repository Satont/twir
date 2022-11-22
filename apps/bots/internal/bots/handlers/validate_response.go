package handlers

import (
	"errors"
	"strings"
)

func ValidateResponseSlashes(response string) error {
	if !strings.HasPrefix(response, "/") || strings.HasPrefix(response, "/me") || strings.HasPrefix(response, "/announce") {
		return nil
	} else {
		return errors.New("slash commands except /me and /announce is disallowed. This response wont be ever sended")
	}
}