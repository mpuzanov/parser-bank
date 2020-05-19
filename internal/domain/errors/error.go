package errors

// ParserError для выдачи ошибок
type ParserError string

func (ee ParserError) Error() string {
	return string(ee)
}

var (
	// ErrFormat .
	ErrFormat = ParserError("формат не подходит")
	// ErrFewFields .
	ErrFewFields = ParserError("мало полей")
)
