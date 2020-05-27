package payments

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/mpuzanov/parser-bank/internal/domain/model"
	"github.com/stretchr/testify/assert"
)

const countPayments = 10000

func BenchmarkSaveToExcel(b *testing.B) {
	testPayments := prepareTestData()
	fileName := "file1.xlsx"
	defer os.Remove(fileName)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testPayments.SaveToExcel(fileName)
	}
}

func TestSaveToExcel(t *testing.T) {
	testPayments := prepareTestData()
	fileName := "file1.xlsx"
	defer os.Remove(fileName)
	err := testPayments.SaveToExcel(fileName)
	assert.Empty(t, err)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Errorf("%s does not exist", fileName)
	}
}

func BenchmarkSaveToExcel2(b *testing.B) {
	testPayments := prepareTestData()
	fileName := "file2.xlsx"
	defer os.Remove(fileName)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testPayments.SaveToExcel2(fileName)
	}
}

func TestSaveToExcel2(t *testing.T) {
	testPayments := prepareTestData()
	fileName := "file2.xlsx"
	defer os.Remove(fileName)
	err := testPayments.SaveToExcel2(fileName)
	assert.Empty(t, err)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Errorf("%s does not exist", fileName)
	}
}

func BenchmarkSaveToExcelStream(b *testing.B) {
	testPayments := prepareTestData()
	tmpfile, err := ioutil.TempFile(".", "fileStream*.xlsx")
	if err != nil {
		b.Errorf("error create tmp file")
	}
	fileName := tmpfile.Name()
	defer os.Remove(fileName)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testPayments.SaveToExcelStream(fileName)
	}
}

func TestSaveToExcelStream(t *testing.T) {
	testPayments := prepareTestData()
	tmpfile, err := ioutil.TempFile(".", "fileStream*.xlsx")
	if err != nil {
		t.Errorf("error create tmp file")
	}
	fileName := tmpfile.Name()
	defer os.Remove(fileName)
	err = testPayments.SaveToExcelStream(fileName)
	assert.Empty(t, err)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Errorf("%s does not exist", fileName)
	}
}

func BenchmarkPrepareTestData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		prepareTestData()
	}
}
func BenchmarkPrepareTestData2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		prepareTestData2()
	}
}

func prepareTestData() *ListPayments {
	// создаём тестовый слайс платежей
	tp := ListPayments{}
	tp.Db = make([]model.Payment, countPayments, countPayments)
	for i := 0; i < countPayments; i++ {
		tp.Db[i].Occ = rand.Intn(999999)
		tp.Db[i].Address = "Пушкинская, 240А, 50"
		tp.Db[i].Date = time.Date(2018, time.August, 30, 0, 0, 0, 0, time.UTC) //time.Now()
		tp.Db[i].Value = rand.Float64()
		tp.Db[i].Commission = rand.Float64()
		tp.Db[i].Fio = "Иванов Иван Иванович"
		tp.Db[i].PaymentAccount = "12345678901234567890"
	}
	return &tp
}

func prepareTestData2() *ListPayments {
	// создаём тестовый слайс платежей
	testPayments := ListPayments{}
	testPayments.Db = make([]model.Payment, 0)
	for i := 0; i < countPayments; i++ {
		p := model.Payment{}
		p.Occ = rand.Intn(999999)
		p.Address = "Пушкинская, 240А, 50"
		p.Date = time.Now()
		p.Value = rand.Float64()
		p.Commission = rand.Float64()
		p.Fio = "Иванов Иван Иванович"
		p.PaymentAccount = "12345678901234567890"
		testPayments.Db = append(testPayments.Db, p)
	}
	return &testPayments
}
