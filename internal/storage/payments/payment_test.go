package payments

import (
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileName, _ := testPayments.SaveToExcel(".", "file1*.xlsx")
		defer os.Remove(fileName)
	}
}

func TestSaveToExcel(t *testing.T) {
	testPayments := prepareTestData()
	fileName, err := testPayments.SaveToExcel(".", "file1*.xlsx")
	assert.Empty(t, err)
	defer os.Remove(fileName)
	assert.FileExists(t, fileName)
}

func BenchmarkSaveToExcel2(b *testing.B) {
	testPayments := prepareTestData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileName, _ := testPayments.SaveToExcel2(".", "file2*.xlsx")
		defer os.Remove(fileName)
	}
}

func TestSaveToExcel2(t *testing.T) {
	testPayments := prepareTestData()
	fileName, err := testPayments.SaveToExcel2(".", "file2*.xlsx")
	assert.Empty(t, err)
	defer os.Remove(fileName)
	assert.FileExists(t, fileName)
}

func BenchmarkSaveToExcelStream(b *testing.B) {
	testPayments := prepareTestData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileName, _ := testPayments.SaveToExcelStream(".", "fileStream*.xlsx")
		defer os.Remove(fileName)
	}
}

func TestSaveToExcelStream(t *testing.T) {
	testPayments := prepareTestData()
	fileName, err := testPayments.SaveToExcelStream(".", "fileStream*.xlsx")
	assert.Empty(t, err)
	defer os.Remove(fileName)
	assert.FileExists(t, fileName)
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
	tp.Db = make([]model.Payment, countPayments)
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

func TestSaveToJSON(t *testing.T) {
	testPayments := prepareTestData()
	fileName, err := testPayments.SaveToJSON(".", "file*.json")
	assert.Empty(t, err)
	defer os.Remove(fileName)
	assert.FileExists(t, fileName)
}

func TestSaveToXML(t *testing.T) {
	testPayments := prepareTestData()
	fileName, err := testPayments.SaveToXML(".", "file*.xml")
	assert.Empty(t, err)
	defer os.Remove(fileName)
	assert.FileExists(t, fileName)
}
