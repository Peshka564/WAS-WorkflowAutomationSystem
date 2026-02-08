package errs

type NotFoundError struct {
}

func (err NotFoundError) Error() string {
	return "Entity not found"
}