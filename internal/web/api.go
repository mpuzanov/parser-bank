package web

import (
	"bufio"
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
	"strings"

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
	w.Header().Add("Content-Disposition", "Attachment")
	http.ServeFile(w, r, path.Join(pathDownload, fileName))
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

		req.Body = http.MaxBytesReader(w, req.Body, 6<<20+1024)
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
			buf := bufio.NewReader(part)
			sniff, _ := buf.Peek(512)
			contentType := http.DetectContentType(sniff)
			//s.logger.Sugar().Infof("contentType: %s", contentType)
			if !strings.Contains(contentType, "text/plain") {
				http.Error(w, "file type not allowed", http.StatusBadRequest)
				return
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
			var maxSize int64 = 2 << 20 //2Mb
			lmt := io.MultiReader(buf, io.LimitReader(part, maxSize-511))
			written, err := io.Copy(tmpfile, lmt)
			if err != nil && err != io.EOF {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if written > maxSize {
				os.Remove(tmpfile.Name())
				http.Error(w, "file size over limit", http.StatusBadRequest)
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

		nameFileXls, err := valuesTotal.SaveToExcel(pathDownload, "file1*.xlsx")
		//nameFileXls, err := valuesTotal.SaveToExcel2(pathDownload, "file1*.xlsx")
		//nameFileXls, err := valuesTotal.SaveToExcelStream(pathDownload, "file1*.xlsx")
		if err != nil {
			s.logger.Error("SaveToExcel", zap.Error(err))
		}
		nameFileJSON, err := valuesTotal.SaveToJSON(pathDownload, "file1*.json")
		if err != nil {
			s.logger.Error("SaveToJSON", zap.Error(err))
		}
		nameFileXML, err := valuesTotal.SaveToXML(pathDownload, "file1*.xml")
		if err != nil {
			s.logger.Error("SaveToJSON", zap.Error(err))
		}
		nameFileXls = "parser-bank/download/" + filepath.Base(nameFileXls)
		nameFileJSON = "parser-bank/download/" + filepath.Base(nameFileJSON)
		nameFileXML = "parser-bank/download/" + filepath.Base(nameFileXML)
		s.logger.Sugar().Debugf("%s %s %s", nameFileXls, nameFileJSON, nameFileXML)
		url := fmt.Sprintf(`Скачать файл (
			<a href=%s target=\"_blank\">xlsx</a>,
			<a href=%s target=\"_blank\">json</a>,
			<a href=%s target=\"_blank\">xml</a>
			)<br>`, nameFileXls, nameFileJSON, nameFileXML)
		strFiles += url

		t, err := template.ParseFiles(fileTemplate)
		if err != nil {
			s.logger.Error("Error:", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			UploadFiles template.HTML
		}{UploadFiles: template.HTML(strFiles)}
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
