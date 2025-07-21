package wire

import (
	"encoding/json"
	"net/http"

	"github.com/etzba/pggo/pkg/logger"
)

type Responder interface {
	SendOK(w http.ResponseWriter, obj ...interface{})
	SendError(w http.ResponseWriter, err error)
}

type Respond struct {
	Logger *logger.Log
}

func (r Respond) SendOK(w http.ResponseWriter, obj ...interface{}) {
	if obj != nil {
		body, err := json.Marshal(obj)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := []byte(`{"message":"could not marshal"}`)
			_, _ = w.Write(res)
		}
		w.Header().Set("Content-Type", "Application/json")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		_, _ = w.Write(body)
	} else {
		w.Write([]byte("OK!")) //nolint:errcheck
	}
}

func (r Respond) SendError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error())) //nolint:errcheck
}
