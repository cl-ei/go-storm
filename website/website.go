package website

import (
	"github.com/wonderivan/logger"
	"net/http"
)

type Server struct {
	addr string
}

func Index(w http.ResponseWriter, r *http.Request) {
	logger.Info("url: ", r.URL.EscapedPath(), ", method: ", r.Method)
	responseDat := "OK!"

	_, err := w.Write([]byte(responseDat))
	if err != nil {
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) Listen() {
	http.HandleFunc("/", Index)
	logger.Info("addr: ", s.addr)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		logger.Error("Error happened in http listener: ", err)
	}
}

func New(addr string) *Server {
	return &Server{addr}
}
