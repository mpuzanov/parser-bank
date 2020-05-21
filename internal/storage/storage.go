package storage

import (
	"encoding/json"
	"fmt"

	"github.com/mpuzanov/parser-bank/internal/domain/model"
)

// ListFormatBanks структура для хранения форматов файлов
type ListFormatBanks model.FormatBanks

// NewFormatBanks создание storage
func NewFormatBanks() *ListFormatBanks {
	return &ListFormatBanks{Db: make([]model.FormatBank, 0)}
}

// Open Загружаем форматы
func (s *ListFormatBanks) Open() error {
	err := json.Unmarshal([]byte(formatData), &s)
	if err != nil {
		return err
	}
	if len(s.Db) == 0 {
		return fmt.Errorf("варианты форматов реестров не загружены")
	}
	return nil
}
