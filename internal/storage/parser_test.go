package storage

import (
	"testing"
	"time"

	"github.com/mpuzanov/parser-bank/internal/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestGetPaymentsVal(t *testing.T) {
	testCases := []struct {
		desc string
		line string
		fb   model.FormatBank
		want model.Payment
		err  error
	}{
		{
			desc: "Тест Почта_D7L1A3S5C6F2",
			line: "6149829;;Пушкинская, 240А, 50;;491;6.38;30.08.2018;",
			fb:   FormatDataMap["D7L1A3S5C6F2"],
			want: model.Payment{Occ: 6149829, Address: "Пушкинская, 240А, 50", Value: 491, Commission: 6.38,
				Date: time.Date(2018, time.August, 30, 0, 0, 0, 0, time.UTC), Fio: "", PaymentAccount: ""},
			err: nil,
		},
		{
			desc: "Тест Почта_D7L1A3S5F2",
			line: "347043;ЗАКИРОВА ФИРАЯ ЯВДАТОВНА;Ижевск, Инкубаторный, д. 6, кв. 6;;389.65;3791;16.01.2020;426032;4260322;",
			fb:   FormatDataMap["D7L1A3S5F2"],
			want: model.Payment{Occ: 347043, Address: "Ижевск, Инкубаторный, д. 6, кв. 6", Value: 389.65, Commission: 0,
				Date: time.Date(2020, time.January, 16, 0, 0, 0, 0, time.UTC), Fio: "ЗАКИРОВА ФИРАЯ ЯВДАТОВНА", PaymentAccount: ""},
			err: nil,
		},
		{
			desc: "Тест Почта_D8L1A3S5С6F2",
			line: "6149829;;Пушкинская, 240А, 50;;491;6.38;1731;30.08.2018;426008;42600805;",
			fb:   FormatDataMap["D8L1A3S5С6F2"],
			want: model.Payment{Occ: 6149829, Address: "Пушкинская, 240А, 50", Value: 491, Commission: 6.38,
				Date: time.Date(2018, time.August, 30, 0, 0, 0, 0, time.UTC), Fio: "", PaymentAccount: ""},
			err: nil,
		},
		{
			desc: "Тест Сбербанк_D1L6A7S8C10",
			line: "20-04-2020;01-28-59;8618;8618999V;401087864321;700154937;ИЖЕВСК, ПАСТУХОВА, Д. 57, КВ. 75;683,46;676,63;6,83",
			fb:   FormatDataMap["D1L6A7S8C10"],
			want: model.Payment{Occ: 700154937, Address: "ИЖЕВСК, ПАСТУХОВА, Д. 57, КВ. 75", Value: 683.46, Commission: 6.83,
				Date: time.Date(2020, time.April, 20, 0, 0, 0, 0, time.UTC), Fio: "", PaymentAccount: ""},
			err: nil,
		},
		{
			desc: "Тест Сбербанк_D3L7A8S5C6",
			line: "8618;8618999V;27/03/2017;9483893;2000.00;20.00;ЛИЦЕВОЙ СЧЕТ: 910000419; АДРЕС: Т.БАРАМЗИНОЙ 7А КВ 124;",
			fb:   FormatDataMap["D3L7A8S5C6"],
			want: model.Payment{Occ: 910000419, Address: "Т.БАРАМЗИНОЙ 7А КВ 124", Value: 2000.00, Commission: 20.00,
				Date: time.Date(2017, time.March, 27, 0, 0, 0, 0, time.UTC), Fio: "", PaymentAccount: ""},
			err: nil,
		},
		{
			desc: "Тест Сбербанк_D1L6A8S10C12F7",
			line: "25-09-2019;07-37-41;8618;8618999V;300863541515;910000667;КУМАЧЕВА ВЕРА АЛЕКСАНДРОВНА;ИЖЕВСК, Т.БАРАМЗИНОЙ, Д. 7А, КВ. 90;0819;3989,76;3945,87;43,89;5",
			fb:   FormatDataMap["D1L6A8S10C12F7"],
			want: model.Payment{Occ: 910000667, Address: "ИЖЕВСК, Т.БАРАМЗИНОЙ, Д. 7А, КВ. 90", Value: 3989.76, Commission: 43.89,
				Date: time.Date(2019, time.September, 25, 0, 0, 0, 0, time.UTC), Fio: "КУМАЧЕВА ВЕРА АЛЕКСАНДРОВНА", PaymentAccount: ""},
			err: nil,
		},
		{
			desc: "Тест Ижкомбанк_D2L4S3A5",
			line: "20190925001512897292;25/09/2019;1000.00; ЛИЦ.СЧЕТ: 700191647; АДРЕС: Ижевск, Кирова, д. 131, кв. 60; ФИО: ;;",
			fb:   FormatDataMap["D2L4S3A5"],
			want: model.Payment{Occ: 700191647, Address: "Ижевск, Кирова, д. 131, кв. 60", Value: 1000.00,
				Date: time.Date(2019, time.September, 25, 0, 0, 0, 0, time.UTC), Fio: "", PaymentAccount: ""},
			err: nil,
		},
		{
			desc: "Тест Сбербанк_D1L6A8S9C11F7",
			line: "01-10-2019;17-35-18;8618;8618999V;550319691514;20040341;КОРОБОВА СВЕТЛАНА ВЛАДИМИРОВНА;ИЖЕВСК, ВОСТОЧНАЯ, Д. 4, КВ. 34;10943,21;10735,29;207,92;",
			fb:   FormatDataMap["D1L6A8S9C11F7"],
			want: model.Payment{Occ: 20040341, Address: "ИЖЕВСК, ВОСТОЧНАЯ, Д. 4, КВ. 34", Value: 10943.21, Commission: 207.92,
				Date: time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC), Fio: "КОРОБОВА СВЕТЛАНА ВЛАДИМИРОВНА", PaymentAccount: ""},
			err: nil,
		},
		{
			desc: "Тест Сбербанк_D3L6A7S5",
			line: "8618;8618999V;02/04/2015;9866889;213.06;ЛИЦ. СЧЕТ: 776150849; АДРЕС: г.Ижевск, Красноармейская ул. д.168 кв.17;",
			fb:   FormatDataMap["D3L6A7S5"],
			want: model.Payment{Occ: 776150849, Address: "г.Ижевск, Красноармейская ул. д.168 кв.17", Value: 213.06, Commission: 0,
				Date: time.Date(2015, time.April, 2, 0, 0, 0, 0, time.UTC), Fio: "", PaymentAccount: ""},
			err: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := fbStore.getPaymentsVal(tC.line, &tC.fb)
			assert.Equal(t, tC.err, err)
			if err == nil {
				assert.NotEmpty(t, got)
				assert.Equal(t, tC.want.Occ, got.Occ)
				assert.Equal(t, tC.want.Address, got.Address)
				assert.Equal(t, tC.want.Value, got.Value)
				assert.Equal(t, tC.want.Commission, got.Commission)
				assert.Equal(t, tC.want.Date, got.Date)
				assert.Equal(t, tC.want.Fio, got.Fio)
				assert.Equal(t, tC.want.PaymentAccount, got.PaymentAccount)
			}
		})
	}
}
