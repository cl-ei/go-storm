package website

import (
	"./handlers"
	"github.com/wonderivan/logger"
	"net/http"
)

type Server struct {
	addr string
}

type httpHandlerFunc func(w http.ResponseWriter, r *http.Request)

func defaultMiddleWare(f httpHandlerFunc) httpHandlerFunc {
	wrapped := func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Access: %s %s", r.Method, r.URL.EscapedPath())
		w.Header().Add("server", "madliar/1.18.9a6 (Darwin, based on golang)")
		f(w, r)
	}
	return wrapped
}

func (s *Server) Listen() {
	// load router.
	for urlPattern, handler := range handlers.UrlMap {
		http.HandleFunc(urlPattern, defaultMiddleWare(handler))
	}

	logger.Info("addr: ", s.addr)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		logger.Error("Error happened in http listener: ", err)
	}
}

func New(addr string) *Server {
	return &Server{addr}
}
