package apperr

type Error struct {
	code    int
	message string
}

func New(code int, text string) error {
	return &Error{code, text}
}

func (e Error) ErrorCode() int {
	return e.code
}

func (e Error) Error() string {
	return e.message
}
