package api

import (
	"fmt"
	"net/http"

	"github.com/nhudson/longbox/internal/api/handler"
	"github.com/sirupsen/logrus"
)

func NewServer(log logrus.FieldLogger, url string) http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(req.Method)
		if req.Method == "GET" {
			if _, err := w.Write([]byte(`{"status":"healthy"}`)); err != nil {
				fmt.Printf("Write error: +%v:", err)
			}
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	m.Handle("/search", handler.Search(log, url))

	return logger(log, m)
}
