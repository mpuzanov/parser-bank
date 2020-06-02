package storage

import (
	"encoding/json"
	"fmt"
	"sort"

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

	for _, fb := range FormatDataMap {
		s.Db = append(s.Db, fb)
	}
	// сортируем по возрастанию
	sort.Slice(s.Db, func(i, j int) bool {
		return s.Db[i].Priority < s.Db[j].Priority
	})
	if len(s.Db) == 0 {
		return fmt.Errorf("варианты форматов реестров не загружены")
	}
	return nil
}

// Open2 Загружаем форматы
func (s *ListFormatBanks) Open2() error {
	err := json.Unmarshal([]byte(formatData), &s)
	if err != nil {
		return err
	}
	if len(s.Db) == 0 {
		return fmt.Errorf("варианты форматов реестров не загружены")
	}
	return nil
}
