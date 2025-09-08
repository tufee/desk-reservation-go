package api

import (
	"net/http"

	"github.com/tufee/desk-reservation-go/internal/middleware"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("GET /home", Home)

	mux.HandleFunc("POST /user", CreateUserHandler)
	mux.HandleFunc("POST /reservation", middleware.AuthMiddleware(CreateReservationHandler))
	mux.HandleFunc("POST /login", LoginHandler)
	return mux
}
