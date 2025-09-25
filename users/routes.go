package users

import "github.com/gorilla/mux"

func RegisterRoutes(r *mux.Router, h *Handler) {

	s := r.PathPrefix("/v1/users").Subrouter()

	s.HandleFunc("", h.Create).Methods("POST")
	s.HandleFunc("", h.List).Methods("GET")
	s.HandleFunc("/{id:[0-9]+}", h.GetByID).Methods("GET")
}
