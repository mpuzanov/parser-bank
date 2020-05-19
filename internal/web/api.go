package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/mpuzanov/parser-bank/internal/domain/model"
	"github.com/mpuzanov/parser-bank/internal/parser"
	"github.com/mpuzanov/parser-bank/pkg/logger"
	"go.uber.org/zap"
)

var (
	fileTemplate = "./templates/index.html"
	pathUpload   = "./upload_files/"
)

func init() {
	if _, err := os.Stat(pathUpload); os.IsNotExist(err) {
		//надо попробовать создать каталог
		err := os.Mkdir(pathUpload, 0777)
		if err != nil {
			logger.LogSugar.Error("dir not exist", zap.String("dir", pathUpload), zap.Error(err))
			os.Exit(1)
		}
	}
}

func (s *myHandler) configRouter() {
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/parser-bank", s.UploadData)
	s.router.PathPrefix("/parser-bank").Handler(http.FileServer(http.Dir("./templates/")))
}

func (s *myHandler) loadStore() {
	if err := s.store.Open(); err != nil {
		s.logger.Error(err.Error())
	} else {
		s.logger.Info("Варианты форматов загружены")
	}
}

func (s *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// UploadData .
func (s *myHandler) UploadData(w http.ResponseWriter, req *http.Request) {

	switch req.Method {

	case "GET":
		s.logger.Debug("GET")
		t, err := template.ParseFiles(fileTemplate)
		if err != nil {
			s.logger.Error("Error:", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)

	case "POST":
		s.logger.Debug("POST")

		reader, err := req.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		strFiles := ""
		count := 0
		valuesTotal := []model.Payments{}
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			if part.FileName() == "" {
				continue
			}

			blobPath := pathUpload + part.FileName()
			dst, err := os.Create(blobPath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() {
				dst.Close()
				err = os.Remove(blobPath)
				if err != nil {
					s.logger.Error("Remove", zap.Error(err))
				} else {
					s.logger.Debug("File has been deleted successfully.", zap.String("fileName", blobPath))
				}
			}()

			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			values, err := parser.ReadFile(blobPath, &s.store.FormatBanks, s.logger)
			if err != nil {
				s.logger.Error("Error:", zap.Error(err))
				http.Error(w, err.Error(), 500)
				return
			}
			count++
			s.logger.Info("", zap.Int("Count values", len(values)))
			strFiles += fmt.Sprintf("%d. %s - кол-во платежей: %d<br>", count, part.FileName(), len(values))

			valuesTotal = append(valuesTotal, values...)
		}
		strFiles += fmt.Sprintf("Итого платежей: %d<br>", len(valuesTotal))

		parsedJSON, err := json.Marshal(valuesTotal)
		if err != nil {
			s.logger.Error("Error json.Marshal", zap.Error(err))
		}
		jsData, err := Prettyprint(parsedJSON)
		if err != nil {
			s.logger.Error("Error Prettyprint", zap.Error(err))
		}
		//fmt.Println(sting(jsData))

		t, err := template.ParseFiles(fileTemplate)
		if err != nil {
			s.logger.Error("Error:", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			UploadFiles template.HTML
			ContentJSON string
		}{UploadFiles: template.HTML(strFiles), ContentJSON: string(jsData)}
		t.Execute(w, data)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *myHandler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//logger:=s.logger.With(zap.String("remote_addr", r.RemoteAddr))
		s.logger.Info("Request",
			zap.String("Method", r.Method),
			zap.String("URI", r.RequestURI),
			zap.String("remoteaddr", r.RemoteAddr),
		)
		next.ServeHTTP(w, r)
	})
}

//Prettyprint Делаем красивый json с отступами
func Prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "    ")
	return out.Bytes(), err
}
