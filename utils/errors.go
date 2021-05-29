package utils

type ReadableError struct {
	Err error
	Msg string
}

func (e *ReadableError) Error() string {
	return e.Msg
}
