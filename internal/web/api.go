package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/mpuzanov/parser-bank/internal/storage/payments"
	"go.uber.org/zap"
)

var (
	fileTemplate             = "./templates/index.html"
	pathDownload, pathUpload string
)

func (s *myHandler) configRouter() {
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/parser-bank", s.UploadData)
	s.router.HandleFunc("/parser-bank/upload", s.UploadData)
	s.router.HandleFunc("/parser-bank/download/{file}", s.DownloadFile)
	s.router.PathPrefix("/parser-bank").Handler(http.FileServer(http.Dir("./templates/")))

	pathDownload = path.Join(s.cfg.PathTmp, "out")
	pathUpload = path.Join(s.cfg.PathTmp, "in")
}

// loadStore загружаем возможные форматы банковских реестров
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

// DownloadFile .
func (s *myHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["file"]
	//fmt.Fprint(w, "file=", fileName)
	http.ServeFile(w, r, path.Join(pathDownload, fileName))

	// modtime := time.Now()
	// content := randomContent(modtime.UnixNano(), 1024)

	// // ServeContent uses the fileName for mime detection
	// //const fileName = "random.txt"

	// // tell the browser the returned content should be downloaded
	// w.Header().Add("Content-Disposition", "Attachment")

	// http.ServeContent(w, r, fileName, modtime, content)
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
		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "POST":
		s.logger.Debug("POST")

		reader, err := req.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		strFiles := ""
		count := 0
		valuesTotal := payments.ListPayments{}
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if part.FileName() == "" {
				continue
			}

			tmpfile, err := ioutil.TempFile(pathUpload, part.FileName())
			if err != nil {
				s.logger.Error("TempFile", zap.Error(err))
			}
			tmpFileName := tmpfile.Name()
			defer func() {
				tmpfile.Close()
				err = os.Remove(tmpFileName)
				if err != nil {
					s.logger.Error("Remove tmpfile", zap.Error(err))
				}
			}()
			if _, err := io.Copy(tmpfile, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			sError := ""
			values, err := s.store.ReadFile(tmpFileName, s.logger)
			if err != nil {
				s.logger.Error("Error:", zap.Error(err))
				sError = err.Error()
			}
			count++
			s.logger.Info("", zap.Int("Count values", len(values)))
			strFiles += fmt.Sprintf("%d. %s - кол-во платежей: %d %s<br>", count, part.FileName(), len(values), sError)

			valuesTotal.Db = append(valuesTotal.Db, values...)
		}
		strFiles += fmt.Sprintf("Итого платежей: %d<br>", len(valuesTotal.Db))

		tmpfile, err := ioutil.TempFile(pathDownload, "file*.xlsx")
		if err != nil {
			s.logger.Error("TempFile", zap.Error(err))
		}
		nameFile := tmpfile.Name()
		err = valuesTotal.SaveToExcelStream(nameFile)
		//err = valuesTotal.SaveToExcel2(nameFile)
		if err != nil {
			s.logger.Error("SaveToExcel", zap.Error(err))
		} else {
			nameFile := "parser-bank/download/" + filepath.Base(nameFile)
			url := fmt.Sprintf("<a href=\"%s\" target=\"_blank\">Скачать файл</a>", nameFile)
			strFiles += url
		}

		parsedJSON, err := json.Marshal(valuesTotal)
		if err != nil {
			s.logger.Error("Error json.Marshal", zap.Error(err))
		}
		jsData, err := Prettyprint(parsedJSON)
		if err != nil {
			s.logger.Error("Error Prettyprint", zap.Error(err))
		}

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
		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// newURL := "/parser-bank/download"
		// http.Redirect(w, req, newURL, http.StatusSeeOther)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// logRequest логируем доступ к сайту
func (s *myHandler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
