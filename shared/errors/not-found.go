package errs

import "fmt"

type NotFoundError struct {
	EntityName string
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", err.EntityName)
}