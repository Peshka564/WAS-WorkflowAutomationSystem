package errs

type InvalidInputError struct {
}

func (err InvalidInputError) Error() string {
	return "Invalid Input"
}