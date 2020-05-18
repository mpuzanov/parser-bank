package web

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/mpuzanov/parser-bank/internal/parser"
	"go.uber.org/zap"
)

var (
	fileTemplate = "./templates/index.html"
	pathUpload   = "./data/"
)

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
	if req.Method == "GET" {
		s.logger.Debug("GET")
		t, err := template.ParseFiles(fileTemplate)
		if err != nil {
			s.logger.Error("Error:", zap.Error(err))
			http.Error(w, err.Error(), 500)
			return
		}
		t.Execute(w, nil)

	} else if req.Method == "POST" {
		s.logger.Debug("POST")
		file, handler, err := req.FormFile("uploadfile")
		if err != nil {
			s.logger.Error("Error:", zap.Error(err))
			http.Error(w, err.Error(), 500)
			return
		}

		defer file.Close()
		if err != nil {
			s.logger.Error("Error while Posting data")
			t, _ := template.ParseFiles(fileTemplate)
			t.Execute(w, nil)
		} else {
			s.logger.Info("OpenFile", zap.String("Filename:", handler.Filename))
			if _, err := os.Stat(pathUpload); os.IsNotExist(err) {
				s.logger.Error("dir not exist", zap.String("dir", pathUpload), zap.Error(err))
				return
			}
			blobPath := pathUpload + handler.Filename
			f, err := os.OpenFile(blobPath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				s.logger.Error("OpenFile", zap.Error(err))
				t, _ := template.ParseFiles(fileTemplate)
				t.Execute(w, nil)
			}
			defer func() {
				f.Close()
				err = os.Remove(blobPath)
				if err != nil {
					s.logger.Error("Remove", zap.Error(err))
				} else {
					s.logger.Debug("File has been deleted successfully.", zap.String("fileName", blobPath))
				}
			}()
			io.Copy(f, file)
			values, err := parser.ReadFile(blobPath, &s.store.FormatBanks, s.logger)
			if err != nil {
				s.logger.Error("Error:", zap.Error(err))
				http.Error(w, err.Error(), 500)
				return
			}

			s.logger.Info("", zap.Int("Count values", len(values)))
			parsedJSON, err := json.Marshal(values)
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
				http.Error(w, err.Error(), 500)
				return
			}
			t.Execute(w, string(jsData))
		}
	} else {
		s.logger.Error("Error while Posting data")
		t, err := template.ParseFiles(fileTemplate)
		if err != nil {
			s.logger.Error("Error:", zap.Error(err))
			http.Error(w, err.Error(), 500)
			return
		}
		t.Execute(w, nil)
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
