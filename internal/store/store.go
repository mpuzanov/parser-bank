package store

import (
	"encoding/json"
	"fmt"

	"github.com/mpuzanov/parser-bank/internal/domain"
)

// Store ...
type Store struct {
	FormatBanks domain.FormatBanks
}

// New .
func New() *Store {
	return &Store{}
}

// Open .
func (s *Store) Open() error {
	fb, err := LoadFormatBank()
	if err != nil {
		return err
	}
	s.FormatBanks = *fb
	return nil
}

// LoadFormatBank загрузить известные форматы реестров
func LoadFormatBank() (*domain.FormatBanks, error) {
	var formatBanks domain.FormatBanks

	// jsonFile, err := os.Open("./configs/format.json")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer jsonFile.Close()
	// byteValue, _ := ioutil.ReadAll(jsonFile)

	//formatData
	err := json.Unmarshal([]byte(formatData), &formatBanks)
	//json.Unmarshal(byteValue, &formatBanks)
	if err != nil || len(formatBanks.FormatBanks) == 0 {
		return nil, fmt.Errorf("варианты форматов реестров не загружены. %w", err)
	}
	//fmt.Println(formatBanks)

	return &formatBanks, nil
}
