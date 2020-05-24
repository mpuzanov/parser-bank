package model

// FormatBank структура файла с платежами из банка
type FormatBank struct {
	Name             string `json:"name"`
	Visible          bool   `json:"visible,omitempty"`
	Ext              string `json:"ext"`
	CodePage         string `json:"code_page"`
	ExtBank          string `json:"ext_bank"`
	CharZag          string `json:"char_zag"`
	CharRazd         string `json:"char_razd"`
	FilenameFilter   string `json:"filename_filter"`
	LicNo            int    `json:"lic_no"`
	LicSize          int    `json:"lic_size"`
	LicName          string `json:"lic_name"`
	DataPlatNo       int    `json:"data_plat_no"`
	DataPlatSize     int    `json:"data_plat_size"`
	Dateseparator    string `json:"dateseparator"`
	Decimalseparator string `json:"decimalseparator"`
	SummaNo          int    `json:"summa_no"`
	SummaSize        int    `json:"summa_size"`
	AddressNo        int    `json:"address_no"`
	AddressSize      int    `json:"address_size"`
	AddressName      string `json:"address_name,omitempty"`
	FioNo            int    `json:"fio_no,omitempty"`
	FioName          string `json:"fio_name,omitempty"`
	CommissNo        int    `json:"commis_no,omitempty"`
	RaschName        string `json:"rasch_name,omitempty"`
	RaschNo          int    `json:"rasch_no,omitempty"`
}

// FormatBanks .
type FormatBanks struct {
	Db []FormatBank `json:"dataset"`
}

/*
{
        "name": "Почта_D7L1A3S5",
        "visible": true,
        "ext": "TXT",
        "code_page": "ASCII",
        "ext_bank": "",
        "char_zag": "#",
        "char_razd": ";",
        "filename_filter": "*.*",
        "lic_no": 1,
        "lic_size": 9,
        "data_plat_no": 7,
        "data_plat_size": 10,
        "dateseparator": ".",
        "decimalseparator": ".",
        "summa_no": 5,
        "summa_size": 9,
        "address_no": 3,
        "address_size": 50,
        "fio_no": 2
    },
*/
