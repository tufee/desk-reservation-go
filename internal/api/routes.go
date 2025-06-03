package api

import "net/http"

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /user", CreateUserHandler)
	mux.HandleFunc("POST /reservation", CreateReservationHandler)
	return mux
}
