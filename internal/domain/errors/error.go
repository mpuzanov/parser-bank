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

	// ErrCommissionNotFound мало полей
	ErrCommissionNotFound = ParserError("поле <Commission> не найдено")
	// ErrCommissionBadFormat мало полей
	ErrCommissionBadFormat = ParserError("поле <Commission> не правильный формат")
)
