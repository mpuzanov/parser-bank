package storage

import (
	"github.com/mpuzanov/parser-bank/internal/domain/model"
)

// ListFormats список определяемых форматов
var ListFormats = map[string]string{
	"D7L1A3S5C6F2":   "Почта_D7L1A3S5C6F2",      //1
	"D8L1A3S5С6F2":   "Почта_D8L1A3S5С6F2",      //2
	"D1L6A7S8C10":    "Сбербанк_D1L6A7S8C10",    //3
	"D1L6A8S10C12F7": "Сбербанк_D1L6A8S10C12F7", //4
	"D1L6A8S9C11F7":  "Сбербанк_D1L6A8S9C11F7",  //5
	"D2L4A5S3":       "Ижкомбанк_D2L4S3A5",      //6
	"D3L7A8S5C6":     "Сбербанк_D3L7A8S5C6",     //7
	//"":"",
	//"":"",
}

// FormatDataMap список форматов для анализа платежей
var FormatDataMap = map[string]model.FormatBank{
	"D7L1A3S5С6F2": {Name: "Почта_D7L1A3S5С6F2", CharZag: "#", CharRazd: ";", Dateseparator: ".", Decimalseparator: ".",
		DataPlatNo: 7, LicNo: 1, AddressNo: 3, SummaNo: 5, CommissNo: 6, FioNo: 2},
	"D8L1A3S5С6F2": {Name: "Почта_D8L1A3S5С6F2", CharZag: "#", CharRazd: ";", Dateseparator: ".", Decimalseparator: ".",
		DataPlatNo: 8, LicNo: 1, AddressNo: 3, SummaNo: 5, CommissNo: 6, FioNo: 2},
	"D1L6A7S8C10": {Name: "Сбербанк_D1L6A7S8C10", CharZag: "=", CharRazd: ";", Dateseparator: "-", Decimalseparator: ",",
		DataPlatNo: 1, LicNo: 6, AddressNo: 7, SummaNo: 8, CommissNo: 10},
	"D3L7A8S5C6": {Name: "Сбербанк_D3L7A8S5C6", CharZag: "~", CharRazd: ";", Dateseparator: "/", Decimalseparator: ",",
		DataPlatNo: 3, LicNo: 7, LicName: "ЛИЦЕВОЙ СЧЕТ", AddressNo: 8, AddressName: "АДРЕС", SummaNo: 5, CommissNo: 6},
	"D1L6A8S10C12F7": {Name: "Сбербанк_D1L6A8S10C12F7", CharZag: "=", CharRazd: ";", Dateseparator: "-", Decimalseparator: ",",
		DataPlatNo: 1, LicNo: 6, AddressNo: 8, SummaNo: 10, CommissNo: 12, FioNo: 7},
	"D2L4S3A5": {Name: "Ижкомбанк_D2L4S3A5", CharZag: "~", CharRazd: ";", Dateseparator: "/", Decimalseparator: ".",
		DataPlatNo: 2, LicNo: 4, LicName: "ЛИЦ.СЧЕТ", AddressNo: 5, AddressName: "АДРЕС", SummaNo: 3},
	"D1L6A8S9C11F7": {Name: "Сбербанк_D1L6A8S9C11F7", CharZag: "=", CharRazd: ";", Dateseparator: "-", Decimalseparator: ",",
		DataPlatNo: 1, LicNo: 6, AddressNo: 8, SummaNo: 9, CommissNo: 11, FioNo: 7},
}

// FormatData список форматов для анализа платежей в JSON
var formatData = `
{ "dataset": [
    {
        "name": "Почта_D7L1A3S5C6F2",
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
        "commis_no": 6,
        "fio_no": 2
    },
    {
        "name": "Почта_D8L1A3S5С6F2",
        "visible": true,
        "ext": "TXT",
        "code_page": "ASCII",
        "ext_bank": "",
        "char_zag": "#",
        "char_razd": ";",
        "filename_filter": "*.*",
        "lic_no": 1,
        "lic_size": 9,
        "data_plat_no": 8,
        "data_plat_size": 10,
        "dateseparator": ".",
        "decimalseparator": ".",
        "summa_no": 5,
        "summa_size": 9,
        "address_no": 3,
        "address_size": 50,
        "commis_no": 6
    },
    {
        "name": "Сбербанк_D1L6A7S8C10",
        "visible": true,
        "ext": "TXT",
        "code_page": "ASCII",
        "ext_bank": "",
        "char_zag": "=",
        "char_razd": ";",
        "filename_filter": "*.*",
        "lic_no": 6,
        "lic_size": 9,
        "data_plat_no": 1,
        "data_plat_size": 10,
        "dateseparator": "-",
        "decimalseparator": ",",
        "summa_no": 8,
        "summa_size": 9,
        "address_no": 7,
        "address_size": 50,
        "lic_name": "",
        "address_name": "",
        "commis_no": 10
    },
    {
        "name": "Сбербанк_D1L6A8S10C12F7",
        "visible": true,
        "ext": "TXT",
        "code_page": "ASCII",
        "ext_bank": "",
        "char_zag": "=",
        "char_razd": ";",
        "filename_filter": "*.*",
        "lic_no": 6,
        "lic_size": 9,
        "data_plat_no": 1,
        "data_plat_size": 10,
        "dateseparator": "-",
        "decimalseparator": ",",
        "summa_no": 10,
        "summa_size": 9,
        "address_no": 8,
        "address_size": 50,
        "lic_name": "",
        "address_name": "",
        "commis_no": 12
    },
	{
        "name": "Сбербанк_D1L6A8S9C11F7",
        "visible": true,
        "ext": "TXT",
        "code_page": "ASCII",
        "ext_bank": "",
        "char_zag": "=",
        "char_razd": ";",
        "filename_filter": "*.*",
        "lic_no": 6,
        "lic_size": 9,
        "data_plat_no": 1,
        "data_plat_size": 10,
        "dateseparator": "-",
        "decimalseparator": ",",
        "summa_no": 9,
        "summa_size": 9,
        "address_no": 8,
        "address_size": 50,
        "lic_name": "",
        "commis_no": 11,
        "fio_no": 7
    },
    {
        "name": "Ижкомбанк_D2L4S3A5",
        "visible": true,
        "ext": "TXT",
        "code_page": "ASCII",
        "ext_bank": "",
        "char_zag": "~",
        "char_razd": ";",
        "filename_filter": "*.*",
        "lic_no": 4,
        "lic_size": 9,
        "lic_name": "ЛИЦ.СЧЕТ",		
        "data_plat_no": 2,
        "data_plat_size": 10,
        "dateseparator": "/",
        "decimalseparator": ".",
        "summa_no": 3,
        "summa_size": 9,
        "address_no": 5,
        "address_size": 50,
        "address_name": "АДРЕС",
        "commis_no": null,
        "rasch_name": null,
        "rasch_no": null,
        "fio_no": 6,
        "fio_name": "ФИО"
    },
    {
        "name": "Сбербанк_D3L7A8S5C6",
        "visible": true,
        "ext": "TXT",
        "code_page": "ASCII",
        "ext_bank": "",
        "char_zag": "~",
        "char_razd": ";",
        "filename_filter": "*.*",
        "lic_no": 7,
        "lic_size": 9,
		"lic_name": "ЛИЦЕВОЙ СЧЕТ",
        "data_plat_no": 3,
        "data_plat_size": 10,
        "dateseparator": "/",
        "decimalseparator": ",",
        "summa_no": 5,
        "summa_size": 9,
        "address_no": 8,
        "address_size": 50,        
        "address_name": "АДРЕС",
        "commis_no": 6,
        "rasch_name": null,
        "rasch_no": null
    }
]}
`
