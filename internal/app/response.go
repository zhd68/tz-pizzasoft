package app

import "net/http"

func (a *App) newErrorResponse(w http.ResponseWriter, err error, statusCode int, msg []byte) {
	a.logger.Error(err)
	w.WriteHeader(statusCode)
	w.Write(msg)
}

func (a *App) newResponse(w http.ResponseWriter, statusCode int, msg []byte) {
	w.WriteHeader(statusCode)
	w.Write(msg)
}
