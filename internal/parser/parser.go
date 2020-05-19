package parser

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

var layoutDate = "2006-01-02"

// ReadFile .
func ReadFile(filePath string, fbs *model.FormatBanks, logger *zap.Logger) ([]model.Payments, error) {

	//создаём выходную переменную
	values := make([]model.Payments, 0)

	//Определяем кодировку файла
	codePage, _ := cpd.FileCodepageDetect(filePath)
	logger.Info("", zap.String("file", filePath), zap.String("codePage", codePage.String()))

	b, err := ioutil.ReadFile(filePath)
	//определяем формат файла реестра
	sf, err := detectFormatBank(b, fbs, logger)
	if err != nil {
		return nil, fmt.Errorf("ошибка определения формата файла реестра. %w", err)
	}
	if sf == nil {
		return nil, fmt.Errorf("Формат файла реестра не определён")
	}
	logger.Info("Выбрали", zap.String("Формат", sf.Name))

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
		if len(line) == 0 {
			continue
		}

		line = strings.TrimSpace(line)
		if len(line) < 50 {
			continue
		}
		if strings.HasPrefix(line, sf.CharZag) {
			continue
		}

		val, err := getPaymentsVal(line, sf, logger)
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

	logger.Sugar().Infof("Итоги: кол-во: %d, сумма: %8.2f, комиссия: %8.2f", len(values), totalValue, totalCommission)
	return values, nil
}

// detectFormatBank определить формат реестра платежей
func detectFormatBank(r []byte, fbs *model.FormatBanks, logger *zap.Logger) (*model.FormatBank, error) {
	var res model.FormatBank

	if len(fbs.FormatBanks) == 0 {
		return &res, fmt.Errorf("Таблица форматов пуста")
	}
	//fmt.Println(r)
	for _, fm := range fbs.FormatBanks {

		logger.Debug("Проверяем: ", zap.String("формат", fm.Name))
		got, err := checkBankReestr(r, &fm, logger)
		if err == nil && got {
			zap.S().Debug("формат подошёл")
			res = fm
			break
		} else {
			logger.Debug("не подходит", zap.Error(err))
		}
	}
	return &res, nil
}

func checkBankReestr(r []byte, sf *model.FormatBank, logger *zap.Logger) (bool, error) {
	var err error
	res := false
	reader := bufio.NewReader(strings.NewReader(string(r)))
	logger.Sugar().Debugf("Проверяем checkBankReestr: %s, Comma: `%s`, Comment: `%s`", sf.Name, sf.CharRazd, sf.CharZag)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return res, err
		}
		line = strings.TrimSpace(line)
		if len(line) < 50 {
			continue
		}
		if strings.HasPrefix(line, sf.CharZag) {
			continue
		}
		_, err = getPaymentsVal(line, sf, logger)
		if err != nil {
			return res, err
		}
		logger.Sugar().Debug("Выбор: ", sf.Name)
		res = true
		break
	}
	return res, err
}

// getPaymentsVal конвертируем строку файла в структуру платежа
func getPaymentsVal(line string, sf *model.FormatBank, logger *zap.Logger) (model.Payments, error) {
	var res model.Payments
	var err error

	record := strings.Split(line, sf.CharRazd)
	logger.Sugar().Debug("record: ", record)
	countField := len(record)
	// for i := 0; i < countField; i++ {
	// 	logger.Sugar().Debugf("Field %d: %s", i, record[i])
	// }
	if countField < 4 {
		return res, errors.ErrFewFields
	}

	//logger.Sugar().Debug("countField: ", countField)
	// проверяем дату
	if sf.DataPlatNo < countField {
		tmpStr := record[sf.DataPlatNo-1]
		layoutDate := fmt.Sprintf("02%s01%s2006", sf.Dateseparator, sf.Dateseparator)
		res.Date, err = time.Parse(layoutDate, tmpStr)
		if err != nil {
			return res, err
		}
		logger.Sugar().Debug("date ", res.Date)
	} else {
		return res, err
	}
	// проверяем сумму
	if sf.SummaNo < countField {
		tmpStr := strings.TrimSpace(record[sf.SummaNo-1])
		tmpStr = strings.ReplaceAll(tmpStr, ",", ".")
		res.Value, err = strconv.ParseFloat(tmpStr, 64)
		if err != nil {
			return res, err
		}
		logger.Sugar().Debug("value ", res.Value)
	} else {
		return res, err
	}
	// проверяем лицевой счёт
	if sf.LicNo < countField {
		tmpStr := record[sf.LicNo-1]
		if sf.LicName == "" {
			res.Occ, err = strconv.Atoi(tmpStr)
			if err != nil {
				return res, err
			}
		} else {
			fields := strings.Split(tmpStr, ":")
			//fmt.Println("len fields: ", len(fields))
			if len(fields) < 2 {
				return res, err
			}
			res.Occ, err = strconv.Atoi(strings.TrimSpace(fields[1]))
			if err != nil {
				return res, err
			}
		}
		logger.Sugar().Debug("occ ", res.Occ)
	} else {
		return res, err
	}
	if sf.CommissNo > 0 {
		if sf.CommissNo < countField {
			tmpStr := strings.TrimSpace(record[sf.CommissNo-1])
			tmpStr = strings.ReplaceAll(tmpStr, ",", ".")
			res.Commission, err = strconv.ParseFloat(tmpStr, 64)
			if err != nil {
				return res, err
			}
			logger.Sugar().Debug("Commission: ", res.Commission)
		} else {
			return res, errors.ErrFormat
		}
	}
	// проверяем адрес
	if sf.AddresNo < countField {
		tmpStr := record[sf.AddresNo-1]
		if sf.LicName == "" {
			res.Address = tmpStr
		} else {
			fields := strings.Split(tmpStr, ":")
			if len(fields) < 2 {
				return res, err
			}
			res.Address = strings.TrimSpace(fields[1])
		}
	} else {
		return res, errors.ErrFormat
	}
	if sf.FioNo > 0 {
		if sf.FioNo < countField {
			tmpStr := record[sf.FioNo-1]
			res.Fio = strings.TrimSpace(tmpStr)
		}
	}
	//logger.Sugar().Info("getPaymentsVal", zap.Any("Fields", res))

	return res, err
}
