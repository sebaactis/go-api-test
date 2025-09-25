package main

import (
	"fmt"
	"log"
	"sebaactis/go-api-simple/db"
	"sebaactis/go-api-simple/models"
	"sebaactis/go-api-simple/users"

	"net/http"

	"github.com/gorilla/mux"
)

var DSN = "host=localhost port=5432 user=app password=secret dbname=app sslmode=disable"

func main() {
	// Conexion
	gbd, err := db.Connection(DSN)

	if err != nil {
		log.Fatal(err)
	}

	gbd.AutoMigrate(&models.User{})

	// Router
	r := mux.NewRouter()

	usersHandler := users.NewHandler(gbd)
	users.RegisterRoutes(r, usersHandler)

	fmt.Println("Escuchando en el puerto 8080")
	http.ListenAndServe(":8080", r)

}
