package utils

import "fmt"

func HandleError(operation string, err error) error {
	if err == nil {
		return nil
	}
	Logf("[ERROR] %s: %v", operation, err)
	return fmt.Errorf("%s: %w", operation, err)
}

func HandleErrorf(operation string, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	Logf("[ERROR] %s: %v", operation, err)
	return err
}

func WrapError(operation string, err error, context string) error {
	if err == nil {
		return nil
	}
	Logf("[ERROR] %s (%s): %v", operation, context, err)
	return fmt.Errorf("%s (%s): %w", operation, context, err)
}

func Close() {
	CloseLogger()
}

