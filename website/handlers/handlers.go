package handlers

import (
	"net/http"
)

var UrlMap = map[string]func(http.ResponseWriter, *http.Request){
	"/":   Index,
	"/lt": ProcLtStatus,
}

func Index(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	if _, err := response.Write([]byte("OK!")); err != nil {
	}
}

func ProcLtStatus(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	if _, err := response.Write([]byte("OK status -->!")); err != nil {
	}
}
