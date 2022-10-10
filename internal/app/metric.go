package app

import (
	"net/http"
)

func (a *App) heartbeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
