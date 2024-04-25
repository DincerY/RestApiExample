package main

import (
	"RestApiExample/app"
	"RestApiExample/controllers"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()

	router.Use(app.JwtAuthentication)

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/user/get/{userId}", controllers.GetUser).Methods("GET")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("Running on port " + port)

	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		fmt.Println(err)
	}

}
