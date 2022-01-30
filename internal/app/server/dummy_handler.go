package server

import (
	"net/http"
	"time"
)

func (s *GophermartService) DummyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := time.After(time.Second * 10)
		<-c
		_, err := w.Write([]byte(s.Config.ServerAddress))
		if err != nil {
			http.Error(w, "error", 500)
			return
		}
	}
}
