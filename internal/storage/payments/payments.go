package payments

import (
	"unicode/utf8"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/mpuzanov/parser-bank/internal/domain/model"
	"github.com/tealeg/xlsx"
	"go.uber.org/zap"
)

// ListPayments структура для хранения платежей
type ListPayments model.Payments

var (
	//HeaderDoc список полей в заголовке
	HeaderDoc = []string{"Occ", "Address", "Date", "Value", "Commission", "Fio", "PaymentAccount"}

	//withHeader ширина колонок
	withHeader = make(map[string]int)
)

// SaveToExcel сохраняем данные в файл
func (s *ListPayments) SaveToExcel(fileName string) error {

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()

	sheet, err = file.AddSheet("Платежи")
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

	//данные
	for index := 0; index < len(s.Db); index++ {
		row = sheet.AddRow()
		// добавляем поля в строке
		// последовательность полей:  "Occ", "Address", "Date", "Value", "Commission", "Fio", "PaymentAccount"
		j := 0 //Occ
		cell = row.AddCell()
		cell.SetInt(s.Db[index].Occ)
		cell.SetStyle(dataStyle)
		withHeader[HeaderDoc[j]] = 10

		j = 1 //Address
		cell = row.AddCell()
		cell.Value = s.Db[index].Address
		cell.SetStyle(dataStyle)
		if utf8.RuneCountInString(cell.Value) > withHeader[HeaderDoc[j]] {
			withHeader[HeaderDoc[j]] = utf8.RuneCountInString(cell.Value)
		}

		j = 2 //Date
		cell = row.AddCell()
		//cell.Value = s.Db[index].Date.Format("2006-01-02")
		cell.SetDate(s.Db[index].Date)
		cell.SetStyle(dataStyle)
		withHeader[HeaderDoc[j]] = 10

		j = 3 //Value
		cell = row.AddCell()
		cell.SetFloatWithFormat(s.Db[index].Value, "#,##0.00")
		cell.SetStyle(dataStyle)
		withHeader[HeaderDoc[j]] = 10

		j = 4 //Commission
		cell = row.AddCell()
		//cell.Value = strconv.FormatFloat(s.Db[index].Commission, 'f', -1, 64)
		//cell.SetFloat(s.Db[index].Commission)
		cell.SetFloatWithFormat(s.Db[index].Commission, "#,##0.00")
		cell.SetStyle(dataStyle)
		withHeader[HeaderDoc[j]] = 10

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
		withHeader[HeaderDoc[j]] = 20
	}
	//Устанавливаем ширину колонок
	for i, col := range sheet.Cols {
		col.Width = float64(withHeader[HeaderDoc[i]])
	}

	err = file.Save(fileName)
	if err != nil {
		return err
	}

	return nil
}

//SaveToExcel2 lib Excelize
func (s *ListPayments) SaveToExcel2(fileName string) error {
	zap.S().Debugf("SaveToExcel2: %s", fileName)

	file := excelize.NewFile()

	sheetName := "Sheet1"
	indexSheet := file.NewSheet(sheetName)
	file.SetActiveSheet(indexSheet)

	expDate := "dd.MM.yyyy"
	styleDate, err := file.NewStyle(&excelize.Style{CustomNumFmt: &expDate})
	if err != nil {
		return err
	}
	styleHeader, err := file.NewStyle(`{"font":{"bold":true,"family":"Times New Roman","size":12}}`)
	if err != nil {
		return err
	}
	styleFloat, err := file.NewStyle(`{"number_format": 4}`)
	if err != nil {
		return err
	}
	//Зададим наименование колонок
	for index := 1; index <= len(HeaderDoc); index++ {
		axis, _ := excelize.CoordinatesToCellName(index, 1)
		if err := file.SetCellValue(sheetName, axis, HeaderDoc[index-1]); err != nil {
			return err
		}
		//err = file.SetCellStyle(sheetName, axis, axis, styleHeader)
	}
	if err := file.SetCellStyle(sheetName, "A1", "G1", styleHeader); err != nil {
		return err
	}
	//данные
	rowNo := 1
	for index := 0; index < len(s.Db); index++ {
		rowNo++
		// добавляем поля в строке
		// последовательность полей:  "Occ", "Address", "Date", "Value", "Commission", "Fio", "PaymentAccount"
		colNo := 1 //Occ
		axis, _ := excelize.CoordinatesToCellName(colNo, rowNo)
		if err := file.SetCellInt(sheetName, axis, s.Db[index].Occ); err != nil {
			return err
		}

		colNo = 2 //Address
		axis, _ = excelize.CoordinatesToCellName(colNo, rowNo)
		if err := file.SetCellStr(sheetName, axis, s.Db[index].Address); err != nil {
			return err
		}

		colNo = 3 //Date
		axis, _ = excelize.CoordinatesToCellName(colNo, rowNo)
		file.SetCellValue(sheetName, axis, s.Db[index].Date)
		if err := file.SetCellStyle(sheetName, axis, axis, styleDate); err != nil {
			return err
		}

		colNo = 4 //Value
		axis, _ = excelize.CoordinatesToCellName(colNo, rowNo)
		if err := file.SetCellFloat(sheetName, axis, s.Db[index].Value, 2, 64); err != nil {
			return err
		}
		if err := file.SetCellStyle(sheetName, axis, axis, styleFloat); err != nil {
			return err
		}

		colNo = 5 //Commission
		axis, _ = excelize.CoordinatesToCellName(colNo, rowNo)
		if err := file.SetCellFloat(sheetName, axis, s.Db[index].Commission, 2, 64); err != nil {
			return err
		}
		if err := file.SetCellStyle(sheetName, axis, axis, styleFloat); err != nil {
			return err
		}

		colNo = 6 //Fio
		axis, _ = excelize.CoordinatesToCellName(colNo, rowNo)
		if err := file.SetCellStr(sheetName, axis, s.Db[index].Fio); err != nil {
			return err
		}

		colNo = 7 //PaymentAccount
		axis, _ = excelize.CoordinatesToCellName(colNo, rowNo)
		if err := file.SetCellStr(sheetName, axis, s.Db[index].PaymentAccount); err != nil {
			return err
		}

	}

	if err := file.SetColWidth(sheetName, "A", "G", 15); err != nil {
		return err
	}
	if err := file.SetColWidth(sheetName, "B", "B", 40); err != nil {
		return err
	}
	if err := file.SetColWidth(sheetName, "G", "G", 25); err != nil {
		return err
	}

	if err := file.SaveAs(fileName); err != nil {
		return err
	}
	return nil
}

//SaveToExcelStream lib Excelize
func (s *ListPayments) SaveToExcelStream(fileName string) error {
	zap.S().Debugf("SaveToExcelStream: %s", fileName)
	file := excelize.NewFile()
	sheetName := "Sheet1"
	streamWriter, err := file.NewStreamWriter(sheetName)
	if err != nil {
		return err
	}

	expDate := "dd.MM.yyyy"
	styleDate, err := file.NewStyle(&excelize.Style{CustomNumFmt: &expDate})
	if err != nil {
		return err
	}
	expFloat := "#,##0.00"
	styleFloat, err := file.NewStyle(&excelize.Style{CustomNumFmt: &expFloat})
	if err != nil {
		return err
	}
	styleHeader, err := file.NewStyle(`{"font":{"bold":true,"family":"Times New Roman","size":12}}`)
	if err != nil {
		return err
	}

	CountCol := len(HeaderDoc)
	rowHeader := make([]interface{}, CountCol)
	for i := 0; i < CountCol; i++ {
		rowHeader[i] = excelize.Cell{StyleID: styleHeader, Value: HeaderDoc[i]}
	}
	if err := streamWriter.SetRow("A1", rowHeader); err != nil {
		return err
	}

	//данные
	rowNo := 1
	for index := 0; index < len(s.Db); index++ {
		row := make([]interface{}, CountCol)
		rowNo++

		row[0] = s.Db[index].Occ
		row[1] = s.Db[index].Address
		row[2] = excelize.Cell{StyleID: styleDate, Value: s.Db[index].Date}
		row[3] = excelize.Cell{StyleID: styleFloat, Value: s.Db[index].Value}
		row[4] = excelize.Cell{StyleID: styleFloat, Value: s.Db[index].Commission}
		row[5] = s.Db[index].Fio
		row[6] = s.Db[index].PaymentAccount

		cell, _ := excelize.CoordinatesToCellName(1, rowNo)
		if err := streamWriter.SetRow(cell, row); err != nil {
			return err
		}
	}

	if err := streamWriter.Flush(); err != nil {
		return err
	}

	file.SetColWidth(sheetName, "A", "G", 15) // тормозит
	file.SetColWidth(sheetName, "B", "B", 40)
	file.SetColWidth(sheetName, "G", "G", 25)

	// if err = file.SetCellStyle(sheetName, "A1", "G1", styleHeader); err != nil {
	// 	return err
	// }
	// if err = file.SetCellStyle(sheetName, "C2", fmt.Sprintf("C%d", len(s.Db)+1), styleDate); err != nil {
	// 	return err
	// }
	// if err := file.SetCellStyle(sheetName, "D2", fmt.Sprintf("D%d", len(s.Db)+1), styleFloat); err != nil {
	// 	return err
	// }
	// if err = file.SetCellStyle(sheetName, "E2", fmt.Sprintf("E%d", len(s.Db)+1), styleFloat); err != nil {
	// 	return err
	// }
	// if err = file.SetCellStyle(sheetName, "A1", "G1", styleHeader); err != nil {
	// 	return err
	// }

	if err := file.SaveAs(fileName); err != nil {
		return err
	}
	return nil
}
