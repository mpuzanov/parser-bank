package payments

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/mpuzanov/parser-bank/internal/domain/model"
	"github.com/tealeg/xlsx"
	"go.uber.org/zap"
)

// ListPayments структура для хранения платежей
type ListPayments model.Payments

type fieldExcel struct {
	Name  string
	With  int
	Style int
	Type  string
}

var (
	//HeaderDoc список полей в заголовке
	HeaderDoc = []string{"Occ", "Address", "Date", "Value", "Commission", "Fio", "PaymentAccount"}
	headerMFC = []string{"№", "Fio", "Date", "Участок", "Occ", "Address", "Value", "Код_услуги", "Commission", "PaymentAccount"}

	headerMap = map[int]fieldExcel{
		0: {Name: "№", With: 10},
		1: {Name: "Fio", With: 20},
		2: {Name: "Date", With: 10, Type: "time.Time"},
		3: {Name: "Участок", With: 10},
		4: {Name: "Occ", With: 10, Type: "int"},
		5: {Name: "Address", With: 40},
		6: {Name: "Value", With: 10, Type: "float64"},
		7: {Name: "Код_услуги", With: 10},
		8: {Name: "Commission", With: 10, Type: "float64"},
		9: {Name: "PaymentAccount", With: 25},
	}

	//withHeader ширина колонок
	withHeader = make(map[string]int)
)

// SaveToExcel1 сохраняем данные в файл
func (s *ListPayments) SaveToExcel1(path, templateFile string) (string, error) {

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	if templateFile == "" {
		templateFile = "file*.xlsx"
	}
	tmpfile, err := ioutil.TempFile(path, templateFile)
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()
	fileName := tmpfile.Name()

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return "", err
	}
	headerFont := xlsx.NewFont(12, "Calibri")
	headerFont.Bold = true
	headerFont.Underline = false
	headerStyle := xlsx.NewStyle()
	headerStyle.Font = *headerFont

	dataFont := xlsx.NewFont(11, "Calibri")
	dataStyle := xlsx.NewStyle()
	dataStyle.Font = *dataFont

	//заполняем заголовок
	row = sheet.AddRow()
	for index := 0; index < len(headerMap); index++ {
		cell = row.AddCell()
		cell.Value = headerMap[index].Name
		cell.SetStyle(headerStyle)
	}

	//данные
	for index := 0; index < len(s.Db); index++ {
		row = sheet.AddRow()
		// добавляем поля в строке
		values := reflect.ValueOf(s.Db[index])
		//fields := reflect.TypeOf(s.Db[index])
		for i := 0; i < len(headerMap); i++ {
			cell = row.AddCell()
			f := values.FieldByName(strings.Title(headerMap[i].Name))
			if f.IsValid() {
				fieldValue := f.Interface()
				switch v := fieldValue.(type) {
				case float64:
					cell.SetFloatWithFormat(v, "#,##0.00")
				case int:
					cell.SetInt(int(v))
				case time.Time:
					cell.SetDate(v)
				default:
					cell.SetValue(v)
				}
				cell.SetStyle(dataStyle)
			}
		}
	}
	//Устанавливаем ширину колонок
	for i, col := range sheet.Cols {
		col.Width = float64(headerMap[i].With)
		//col.Width = float64(withHeader[headerMFC[i]])
	}

	err = file.Save(fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// SaveToExcel сохраняем данные в файл
func (s *ListPayments) SaveToExcel(path, templateFile string) (string, error) {

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	if templateFile == "" {
		templateFile = "file*.xlsx"
	}
	tmpfile, err := ioutil.TempFile(path, templateFile)
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()
	fileName := tmpfile.Name()

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return "", err
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
		return "", err
	}

	return fileName, nil
}

//SaveToExcel22 lib Excelize
func (s *ListPayments) SaveToExcel22(path, templateFile string) (string, error) {
	if templateFile == "" {
		templateFile = "file*.xlsx"
	}
	tmpfile, err := ioutil.TempFile(path, templateFile)
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()
	fileName := tmpfile.Name()
	zap.S().Infof("SaveToExcel2: %s", fileName)

	file := excelize.NewFile()

	sheetName := "Sheet1"
	indexSheet := file.NewSheet(sheetName)
	file.SetActiveSheet(indexSheet)

	expDate := "dd.MM.yyyy"
	styleDate, err := file.NewStyle(&excelize.Style{CustomNumFmt: &expDate})
	if err != nil {
		return "", err
	}
	styleHeader, err := file.NewStyle(`{"font":{"bold":true,"family":"Times New Roman","size":12}}`)
	if err != nil {
		return "", err
	}
	styleFloat, err := file.NewStyle(`{"number_format": 4}`)
	if err != nil {
		return "", err
	}
	//Зададим наименование колонок
	for index := 1; index <= len(headerMap); index++ {
		axis, err := excelize.CoordinatesToCellName(index, 1)
		if err != nil {
			return "", err
		}
		if err := file.SetCellValue(sheetName, axis, headerMap[index-1].Name); err != nil {
			return "", err
		}
		if err := file.SetCellStyle(sheetName, axis, axis, styleHeader); err != nil {
			return "", err
		}
		axis, _ = excelize.ColumnNumberToName(index)
		if err := file.SetColWidth(sheetName, axis, axis, float64(headerMap[index-1].With)); err != nil {
			return "", err
		}
	}

	//данные
	rowNo := 1
	for index := 0; index < len(s.Db); index++ {
		rowNo++
		// добавляем поля в строке
		values := reflect.ValueOf(s.Db[index])

		for i := 0; i < len(headerMap); i++ {
			axis, _ := excelize.CoordinatesToCellName(i+1, rowNo)
			f := values.FieldByName(strings.Title(headerMap[i].Name))
			if f.IsValid() {
				fieldValue := f.Interface()
				switch v := fieldValue.(type) {
				case float64:
					if err := file.SetCellFloat(sheetName, axis, v, 2, 64); err != nil {
						return "", err
					}
					if err := file.SetCellStyle(sheetName, axis, axis, styleFloat); err != nil {
						return "", err
					}
				case int:
					if err := file.SetCellInt(sheetName, axis, v); err != nil {
						return "", err
					}
				case string:
					if err := file.SetCellStr(sheetName, axis, v); err != nil {
						return "", err
					}
				case time.Time:
					if err := file.SetCellValue(sheetName, axis, v); err != nil {
						return "", err
					}
					if err := file.SetCellStyle(sheetName, axis, axis, styleDate); err != nil {
						return "", err
					}
				default:
					if err := file.SetCellValue(sheetName, axis, v); err != nil {
						return "", err
					}
				}
			}
		}
	}

	if err := file.SaveAs(fileName); err != nil {
		return "", err
	}
	return fileName, nil
}

//SaveToExcel2 lib Excelize
func (s *ListPayments) SaveToExcel2(path, templateFile string) (string, error) {
	if templateFile == "" {
		templateFile = "file*.xlsx"
	}
	tmpfile, err := ioutil.TempFile(path, templateFile)
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()
	fileName := tmpfile.Name()
	zap.S().Debugf("SaveToExcel2: %s", fileName)

	file := excelize.NewFile()

	sheetName := "Sheet1"
	indexSheet := file.NewSheet(sheetName)
	file.SetActiveSheet(indexSheet)

	expDate := "dd.MM.yyyy"
	styleDate, err := file.NewStyle(&excelize.Style{CustomNumFmt: &expDate})
	if err != nil {
		return "", err
	}
	styleHeader, err := file.NewStyle(`{"font":{"bold":true,"family":"Times New Roman","size":12}}`)
	if err != nil {
		return "", err
	}
	styleFloat, err := file.NewStyle(`{"number_format": 4}`)
	if err != nil {
		return "", err
	}
	//Зададим наименование колонок
	for index := 1; index <= len(HeaderDoc); index++ {
		axis, _ := excelize.CoordinatesToCellName(index, 1)
		if err := file.SetCellValue(sheetName, axis, HeaderDoc[index-1]); err != nil {
			return "", err
		}
		//err = file.SetCellStyle(sheetName, axis, axis, styleHeader)
	}
	if err := file.SetCellStyle(sheetName, "A1", "G1", styleHeader); err != nil {
		return "", err
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
			return "", err
		}

		colNo = 2 //Address
		axis, _ = excelize.CoordinatesToCellName(colNo, rowNo)
		if err := file.SetCellStr(sheetName, axis, s.Db[index].Address); err != nil {
			return "", err
		}

		colNo = 3 //Date
		axis, _ = excelize.CoordinatesToCellName(colNo, rowNo)
		if err := file.SetCellValue(sheetName, axis, s.Db[index].Date); err != nil {
			return "", err
		}
		if err := file.SetCellStyle(sheetName, axis, axis, styleDate); err != nil {
			return "", err
		}

		colNo = 4 //Value
		axis, _ = excelize.CoordinatesToCellName(colNo, rowNo)
		if err := file.SetCellFloat(sheetName, axis, s.Db[index].Value, 2, 64); err != nil {
			return "", err
		}
		if err := file.SetCellStyle(sheetName, axis, axis, styleFloat); err != nil {
			return "", err
		}

		colNo = 5 //Commission
		axis, _ = excelize.CoordinatesToCellName(colNo, rowNo)
		if err := file.SetCellFloat(sheetName, axis, s.Db[index].Commission, 2, 64); err != nil {
			return "", err
		}
		if err := file.SetCellStyle(sheetName, axis, axis, styleFloat); err != nil {
			return "", err
		}

		colNo = 6 //Fio
		axis, _ = excelize.CoordinatesToCellName(colNo, rowNo)
		if err := file.SetCellStr(sheetName, axis, s.Db[index].Fio); err != nil {
			return "", err
		}

		colNo = 7 //PaymentAccount
		axis, _ = excelize.CoordinatesToCellName(colNo, rowNo)
		if err := file.SetCellStr(sheetName, axis, s.Db[index].PaymentAccount); err != nil {
			return "", err
		}

	}

	if err := file.SetColWidth(sheetName, "A", "G", 15); err != nil {
		return "", err
	}
	if err := file.SetColWidth(sheetName, "B", "B", 40); err != nil {
		return "", err
	}
	if err := file.SetColWidth(sheetName, "G", "G", 25); err != nil {
		return "", err
	}

	if err := file.SaveAs(fileName); err != nil {
		return "", err
	}
	return fileName, nil
}

//SaveToExcelStream lib Excelize
func (s *ListPayments) SaveToExcelStream(path, templateFile string) (string, error) {
	if templateFile == "" {
		templateFile = "file*.xlsx"
	}
	tmpfile, err := ioutil.TempFile(path, templateFile)
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()
	fileName := tmpfile.Name()
	zap.S().Debugf("SaveToExcelStream: %s", fileName)
	file := excelize.NewFile()
	sheetName := "Sheet1"
	streamWriter, err := file.NewStreamWriter(sheetName)
	if err != nil {
		return "", err
	}

	expDate := "dd.MM.yyyy"
	styleDate, err := file.NewStyle(&excelize.Style{CustomNumFmt: &expDate})
	if err != nil {
		return "", err
	}
	expFloat := "#,##0.00"
	styleFloat, err := file.NewStyle(&excelize.Style{CustomNumFmt: &expFloat})
	if err != nil {
		return "", err
	}
	styleHeader, err := file.NewStyle(`{"font":{"bold":true,"family":"Times New Roman","size":12}}`)
	if err != nil {
		return "", err
	}

	CountCol := len(HeaderDoc)
	rowHeader := make([]interface{}, CountCol)
	for i := 0; i < CountCol; i++ {
		rowHeader[i] = excelize.Cell{StyleID: styleHeader, Value: HeaderDoc[i]}
	}
	if err := streamWriter.SetRow("A1", rowHeader); err != nil {
		return "", err
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
			return "", err
		}
	}

	if err := streamWriter.Flush(); err != nil {
		return "", err
	}

	// тормозит установка ширины колонок
	if err := file.SetColWidth(sheetName, "A", "G", 15); err != nil {
		return "", err
	}
	if err := file.SetColWidth(sheetName, "B", "B", 40); err != nil {
		return "", err
	}
	if err := file.SetColWidth(sheetName, "G", "G", 25); err != nil {
		return "", err
	}

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
		return "", err
	}
	return fileName, nil
}

//SaveToJSON save file to json format
func (s *ListPayments) SaveToJSON(path, templateFile string) (string, error) {
	if templateFile == "" {
		templateFile = "file*.xml"
	}
	tmpfile, err := ioutil.TempFile(path, templateFile)
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()
	fileName := tmpfile.Name()

	parsedJSON, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	jsonData, err := Prettyprint(parsedJSON)
	if err != nil {
		return "", err
	}
	jsonFile, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer jsonFile.Close()

	if _, err := jsonFile.Write(jsonData); err != nil {
		return "", err
	}
	jsonFile.Close()
	zap.S().Debugf("JSON data written to %s", fileName)

	return fileName, nil
}

//Prettyprint Делаем красивый json с отступами
func Prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "    ")
	return out.Bytes(), err
}

// SaveToXML save file to xml format
func (s *ListPayments) SaveToXML(path, templateFile string) (string, error) {
	if templateFile == "" {
		templateFile = "file*.xml"
	}
	tmpfile, err := ioutil.TempFile(path, templateFile)
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()
	fileName := tmpfile.Name()

	xmlFile, err := xml.MarshalIndent(s, "", " ")
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(fileName, xmlFile, 0644)
	if err != nil {
		return "", err
	}
	zap.S().Debugf("XML data written to %s", fileName)

	return fileName, nil
}
