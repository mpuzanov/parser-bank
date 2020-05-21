package errors

// ParserError для выдачи ошибок
type ParserError string

func (ee ParserError) Error() string {
	return string(ee)
}

var (
	// ErrListFormatEmpty Таблица форматов пуста
	ErrListFormatEmpty = ParserError("Таблица форматов пуста")
	// ErrFormat формат не подходит
	ErrFormat = ParserError("формат не подходит")
	// ErrFewFields мало полей
	ErrFewFields = ParserError("мало полей")
)
