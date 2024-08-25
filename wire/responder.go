package wire

import (
	"net/http"

	"github.com/etzba/pggo/pkg/logger"
)

type Responder interface {
	SendOK(w http.ResponseWriter)
}

type Respond struct {
	Logger *logger.Log
}

func (r Respond) SendOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK!"))
}
