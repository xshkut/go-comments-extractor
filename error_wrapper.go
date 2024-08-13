package main

import "fmt"

func wrapIfError(err error, msg string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", msg, err)
	}

	return nil
}
