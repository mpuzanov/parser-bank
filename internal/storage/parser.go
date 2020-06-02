package storage

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mpuzanov/parser-bank/internal/domain/errors"
	"github.com/mpuzanov/parser-bank/internal/domain/model"
	"github.com/softlandia/cpd"
	"go.uber.org/zap"
)

// ReadFile возвращаем слайс платежей из заданного файла
func (s *ListFormatBanks) ReadFile(filePath string, logger *zap.Logger) ([]model.Payment, error) {

	//создаём выходную переменную
	values := make([]model.Payment, 0)

	//Определяем кодировку файла
	codePage, _ := cpd.FileCodepageDetect(filePath)
	logger.Debug("", zap.String("file", filePath), zap.String("codePage", codePage.String()))

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла реестра %s. %w", filePath, err)
	}
	//определяем формат файла реестра
	sf, err := s.detectFormatBank(b)
	if err != nil {
		return nil, fmt.Errorf("ошибка определения формата файла реестра. %w", err)
	}
	if sf == nil {
		return nil, fmt.Errorf("Формат файла реестра не определён")
	}
	logger.Debug("Выбрали", zap.String("Формат", sf.Name))

	//Открываем файл
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	//Перекодируем
	readerDecoder, err := cpd.NewReader(f, string(codePage))
	if err != nil {
		return nil, err
	}
	totalValue := 0.0
	totalCommission := 0.0
	//Открываем ридер с буфером
	scanner := bufio.NewScanner(readerDecoder)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if isHeaderFile(line, sf) {
			continue
		}

		val, err := s.getPaymentsVal(line, sf)
		if err != nil {
			return nil, fmt.Errorf("ошибка конвертации. %w", err)
		}
		values = append(values, val)
		totalValue += val.Value
		totalCommission += val.Commission
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	logger.Sugar().Debugf("Итоги: кол-во: %d, сумма: %8.2f, комиссия: %8.2f", len(values), totalValue, totalCommission)
	return values, nil
}

// detectFormatBank определить формат реестра платежей
func (s *ListFormatBanks) detectFormatBank(r []byte) (*model.FormatBank, error) {

	if len(s.Db) == 0 {
		return nil, errors.ErrListFormatEmpty
	}
	for _, fm := range s.Db {

		got, err := s.checkBankReestr(r, &fm)
		if err == nil && got {
			zap.S().Debug("подошёл", zap.String("формат", fm.Name))
			return &fm, nil
		}
		zap.S().Debug("не подходит", zap.String("формат", fm.Name), zap.Error(err))
	}
	return nil, nil
}

// checkBankReestr проверяем подходит ли формат
func (s *ListFormatBanks) checkBankReestr(r []byte, sf *model.FormatBank) (bool, error) {
	var err error
	res := false
	reader := bufio.NewReader(strings.NewReader(string(r)))
	zap.S().Debugf("Проверяем checkBankReestr: %s, Comma: `%s`, Comment: `%s`", sf.Name, sf.CharRazd, sf.CharZag)
	countStrOk := 0

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return res, err
		}
		line = strings.TrimSpace(line)

		if isHeaderFile(line, sf) {
			continue
		}
		_, err = s.getPaymentsVal(line, sf)
		if err != nil {
			return res, err
		}
		countStrOk++
		if countStrOk > 5 {
			break
		}
	}

	if countStrOk > 0 {
		zap.S().Debug("Выбор: ", sf.Name)
		res = true
	}

	return res, err
}

// getPaymentsVal конвертируем строку файла в структуру платежа
func (s *ListFormatBanks) getPaymentsVal(line string, sf *model.FormatBank) (model.Payment, error) {
	var res model.Payment
	var err error

	record := strings.Split(line, sf.CharRazd)
	zap.S().Debug("sf: ", sf)
	zap.S().Debug("record: ", record)
	//logger.LogSugar.Debug("sf: ", sf)
	//if sf.Name == "Почта_D7L1A3S5C6F2" {
	//logger.LogSugar.Debug("record: ", record)
	//}

	countField := len(record)
	// for i := 0; i < countField; i++ {
	// 	logger.LogSugar.Debugf("Field %d: %s", i, record[i])
	// }
	if countField < 4 {
		return res, errors.ErrFewFields
	}

	//logger.LogSugar.Debug("countField: ", countField)
	// проверяем дату
	if sf.DataPlatNo > countField {
		return res, errors.ErrFewFields
	}
	tmpStr := record[sf.DataPlatNo-1]
	layoutDate := fmt.Sprintf("02%s01%s2006", sf.Dateseparator, sf.Dateseparator)
	res.Date, err = time.Parse(layoutDate, tmpStr)
	if err != nil {
		return res, err
	}
	zap.S().Debug("date: ", res.Date)

	// проверяем сумму
	if sf.SummaNo > countField {
		return res, errors.ErrFewFields
	}
	tmpStr = strings.TrimSpace(record[sf.SummaNo-1])
	tmpStr = strings.ReplaceAll(tmpStr, ",", ".")
	// if tmpStr == "" {
	// 	return res, errors.ErrFormat
	// }
	res.Value, err = strconv.ParseFloat(tmpStr, 64)
	if err != nil {
		return res, err
	}
	zap.S().Debug("Summa: ", res.Value)

	// проверяем лицевой счёт
	if sf.LicNo > countField {
		return res, errors.ErrFewFields
	}

	tmpStr = record[sf.LicNo-1]
	if sf.LicName == "" {
		res.Occ, err = strconv.Atoi(tmpStr)
		if err != nil {
			return res, err
		}
	} else {
		fields := strings.Split(tmpStr, ":")
		if len(fields) < 2 {
			return res, errors.ErrFormat
		}
		res.Occ, err = strconv.Atoi(strings.TrimSpace(fields[1]))
		if err != nil {
			return res, err
		}
	}
	zap.S().Debug("occ: ", res.Occ)

	if sf.CommissNo > 0 {
		if sf.CommissNo > countField {
			return res, errors.ErrFormat
		}
		tmpStr := strings.TrimSpace(record[sf.CommissNo-1])
		tmpStr = strings.ReplaceAll(tmpStr, ",", ".")
		if tmpStr == "" {
			return res, errors.ErrCommissionNotFound
		}
		res.Commission, err = strconv.ParseFloat(tmpStr, 64)
		if err != nil {
			return res, err
		}
		if !isCommission(res.Commission, res.Value) {
			return res, fmt.Errorf("Commission=%v Value=%v %w", res.Commission, res.Value, errors.ErrCommissionBadFormat)
		}

		zap.S().Debug("Commission: ", res.Commission)
	}
	// проверяем адрес
	if sf.AddressNo > countField {
		return res, errors.ErrFormat
	}
	tmpStr = record[sf.AddressNo-1]
	if sf.LicName == "" {
		res.Address = tmpStr
	} else {
		fields := strings.Split(tmpStr, ":")
		if len(fields) < 2 {
			return res, errors.ErrFormat
		}
		res.Address = strings.TrimSpace(fields[1])
	}

	if sf.FioNo > 0 {
		if sf.FioNo > countField {
			return res, errors.ErrFormat
		}
		tmpStr := record[sf.FioNo-1]
		if sf.FioName != "" {
			fields := strings.Split(tmpStr, ":")
			if len(fields) < 2 {
				return res, errors.ErrFormat
			}
			tmpStr = fields[1]
		}
		res.Fio = strings.TrimSpace(tmpStr)
	}

	return res, err
}

// isCommission определение что v может являться комиссией платежа
func isCommission(v float64, amount float64) bool {
	result := true
	maxCommission := amount * 0.1
	if v >= maxCommission && v > 50 { // если больше 10% и 50 рублей то врятли это комиссия
		result = false
	}
	return result
}

// isHeaderFile строка является заголовком файла
func isHeaderFile(line string, sf *model.FormatBank) bool {
	result := true

	if len(line) == 0 {
		return result
	}

	if len(line) < 50 {
		return result
	}
	for _, charZag := range sf.CharZag {
		if strings.HasPrefix(line, charZag) {
			return result
		}
	}

	record := strings.Split(line, sf.CharRazd)
	if len(record) < 4 {
		return result
	}

	return false
}
