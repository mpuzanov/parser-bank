package payments

import (
	"strconv"
	"unicode/utf8"

	"github.com/mpuzanov/parser-bank/internal/domain/model"
	"github.com/tealeg/xlsx"
)

// ListPayments структура для хранения платежей
type ListPayments model.Payments

// SaveToExcel сохраняем данные в файл
func (s *ListPayments) SaveToExcel(fileName string) error {

	//HeaderDoc список полей в заголовке
	var HeaderDoc = []string{"Occ", "Address", "Date", "Value", "Commission", "Fio", "PaymentAccount"}
	//withHeader ширина колонок
	var withHeader = make(map[string]int)

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()

	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return err
	}

	headerFont := xlsx.NewFont(12, "Calibri")
	headerFont.Bold = true
	headerFont.Underline = false
	headerStyle := xlsx.NewStyle()
	headerStyle.Font = *headerFont

	dataFont := xlsx.NewFont(11, "Calibri")
	dataStyle := xlsx.NewStyle()
	dataStyle.Font = *dataFont //*xlsx.DefaultFont()

	//Зададим наименование колонок
	row = sheet.AddRow()
	for index := 0; index < len(HeaderDoc); index++ {
		cell = row.AddCell()
		cell.Value = HeaderDoc[index]
		cell.SetStyle(headerStyle)
		withHeader[HeaderDoc[index]] = utf8.RuneCountInString(HeaderDoc[index])
	}
	//fmt.Println(withHeader)
	//данные
	for index := 0; index < len(s.Db); index++ {
		row = sheet.AddRow()
		// добавляем поля в строке
		// последовательность полей:  "Occ", "Address", "Date", "Value", "Commission", "Fio", "PaymentAccount"
		j := 0 //Occ
		cell = row.AddCell()
		cell.Value = strconv.Itoa(s.Db[index].Occ)
		cell.SetStyle(dataStyle)
		if utf8.RuneCountInString(cell.Value) > withHeader[HeaderDoc[j]] {
			withHeader[HeaderDoc[j]] = utf8.RuneCountInString(cell.Value)
		}
		j = 1 //Address
		cell = row.AddCell()
		cell.Value = s.Db[index].Address
		cell.SetStyle(dataStyle)
		if utf8.RuneCountInString(cell.Value) > withHeader[HeaderDoc[j]] {
			withHeader[HeaderDoc[j]] = utf8.RuneCountInString(cell.Value)
		}
		j = 2 //Date
		cell = row.AddCell()
		cell.Value = s.Db[index].Date.Format("2006-01-02")
		cell.SetStyle(dataStyle)
		if utf8.RuneCountInString(cell.Value) > withHeader[HeaderDoc[j]] {
			withHeader[HeaderDoc[j]] = utf8.RuneCountInString(cell.Value)
		}
		j = 3 //Value
		cell = row.AddCell()
		cell.Value = strconv.FormatFloat(s.Db[index].Value, 'f', -1, 64)
		cell.SetStyle(dataStyle)
		if utf8.RuneCountInString(cell.Value) > withHeader[HeaderDoc[j]] {
			withHeader[HeaderDoc[j]] = utf8.RuneCountInString(cell.Value)
		}
		j = 4 //Commission
		cell = row.AddCell()
		cell.Value = strconv.FormatFloat(s.Db[index].Commission, 'f', -1, 64)
		cell.SetStyle(dataStyle)
		if utf8.RuneCountInString(cell.Value) > withHeader[HeaderDoc[j]] {
			withHeader[HeaderDoc[j]] = utf8.RuneCountInString(cell.Value)
		}
		j = 5 //Fio
		cell = row.AddCell()
		cell.Value = s.Db[index].Fio
		cell.SetStyle(dataStyle)
		if utf8.RuneCountInString(cell.Value) > withHeader[HeaderDoc[j]] {
			withHeader[HeaderDoc[j]] = utf8.RuneCountInString(cell.Value)
		}
		j = 6 //PaymentAccount
		cell = row.AddCell()
		cell.Value = s.Db[index].PaymentAccount
		cell.SetStyle(dataStyle)
		if utf8.RuneCountInString(cell.Value) > withHeader[HeaderDoc[j]] {
			withHeader[HeaderDoc[j]] = utf8.RuneCountInString(cell.Value)
		}
	}
	//Устанавливаем ширину колонок
	//fmt.Println(withHeader)
	for i, col := range sheet.Cols {
		col.Width = float64(withHeader[HeaderDoc[i]])
	}

	err = file.Save(fileName)
	if err != nil {
		return err
	}

	return nil
}
