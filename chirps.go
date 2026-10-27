package main

import (
	"fmt"
)

func validateChirp(body string) error {
	if len(body) > 140 {
		return fmt.Errorf("body cannot exceed 140 characters")
	}
	return nil
}