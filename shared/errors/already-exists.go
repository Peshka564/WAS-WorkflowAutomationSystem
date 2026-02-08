package errs

import "fmt"

type AlreadyExists struct {
	EntityName string
}

func (err AlreadyExists) Error() string {
	return fmt.Sprintf("%s already exists", err.EntityName)
}