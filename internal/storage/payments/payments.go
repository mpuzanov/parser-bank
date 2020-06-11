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

	"github.com/mpuzanov/parser-bank/internal/domain/model"
	"github.com/tealeg/xlsx"
	"go.uber.org/zap"
)

// ListPayments структура для хранения платежей
type ListPayments model.Payments

type fieldExcel struct {
	Name     string
	Position int
	With     int
	Style    int
	Type     string
}

var (
	headerMap = map[int]*fieldExcel{
		0: {Name: "№", With: 10},
		1: {Name: "Fio", With: 20},
		2: {Name: "Date", With: 10},
		3: {Name: "Участок", With: 10},
		4: {Name: "Occ", With: 10},
		5: {Name: "Address", With: 40},
		6: {Name: "Value", With: 10},
		7: {Name: "Код_услуги", With: 10},
		8: {Name: "Commission", With: 10},
		9: {Name: "PaymentAccount", With: 25},
	}
	headerName = make(map[string]*fieldExcel, 9)
)

func init() {
	for i := 0; i < len(headerMap); i++ {
		headerName[headerMap[i].Name] = &fieldExcel{Name: headerMap[i].Name, Position: i, With: headerMap[i].With}
	}
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
					cell.SetDate(time.Time(v))
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
	}

	err = file.Save(fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// SaveToExcel2 сохраняем данные в файл
func (s *ListPayments) SaveToExcel2(path, templateFile string) (string, error) {

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
	for index := 0; index < len(headerMap); index++ {
		cell = row.AddCell()
		cell.Value = headerMap[index].Name
		cell.SetStyle(headerStyle)
		if utf8.RuneCountInString(headerMap[index].Name) > headerMap[index].With {
			headerMap[index].With = utf8.RuneCountInString(headerMap[index].Name)
		}
	}

	//данные
	for index := 0; index < len(s.Db); index++ {
		row = sheet.AddRow()
		// добавляем поля в строке
		// последовательность полей:
		// "№", "Fio", "Date", "Участок", "Occ", "Address", "Value", "Код_услуги", "Commission", "PaymentAccount"
		//"№"
		cell = row.AddCell()
		cell.SetInt(index + 1)
		//Fio
		cell = row.AddCell()
		cell.Value = s.Db[index].Fio
		cell.SetStyle(dataStyle)
		if utf8.RuneCountInString(cell.Value) > headerMap[headerName["Fio"].Position].With {
			headerMap[headerName["Fio"].Position].With = utf8.RuneCountInString(cell.Value)
		}
		//Date
		cell = row.AddCell()
		//cell.Value = s.Db[index].Date.Format("2006-01-02")
		cell.SetDate(s.Db[index].Date)
		cell.SetStyle(dataStyle)
		//"Участок"
		row.AddCell()
		//Occ
		cell = row.AddCell()
		cell.SetInt(s.Db[index].Occ)
		cell.SetStyle(dataStyle)
		//Address
		cell = row.AddCell()
		cell.Value = s.Db[index].Address
		cell.SetStyle(dataStyle)
		if utf8.RuneCountInString(cell.Value) > headerMap[headerName["Address"].Position].With {
			headerMap[headerName["Address"].Position].With = utf8.RuneCountInString(cell.Value)
		}
		//Value
		cell = row.AddCell()
		cell.SetFloatWithFormat(s.Db[index].Value, "#,##0.00")
		cell.SetStyle(dataStyle)
		//"Код_услуги"
		row.AddCell()
		//Commission
		cell = row.AddCell()
		cell.SetFloatWithFormat(s.Db[index].Commission, "#,##0.00")
		cell.SetStyle(dataStyle)
		//PaymentAccount
		cell = row.AddCell()
		cell.Value = s.Db[index].PaymentAccount
		cell.SetStyle(dataStyle)
	}
	//Устанавливаем ширину колонок
	for i, col := range sheet.Cols {
		col.Width = float64(headerMap[i].With)
	}

	err = file.Save(fileName)
	if err != nil {
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
