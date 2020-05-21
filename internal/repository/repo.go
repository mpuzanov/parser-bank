package repository

import (
	"github.com/mpuzanov/parser-bank/internal/domain/model"
	"go.uber.org/zap"
)

// StorageFormatBanks .
type StorageFormatBanks interface {
	Open() error
	ReadFile(filePath string, logger *zap.Logger) ([]model.Payment, error)
}

// StoragePayments .
type StoragePayments interface {
	SaveToExcel(fileName string) error
}
