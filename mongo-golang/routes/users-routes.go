package routes

import (
	"github.com/babu10103/mongo-golang/controllers"
	"github.com/gorilla/mux"
)

var RegisterUserRoutes = func(router *mux.Router, uc controllers.UserController) {
	router.HandleFunc("/users", uc.GetAllUsers).Methods("GET")
	router.HandleFunc("/user", uc.CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", uc.GetUser).Methods("GET")
	router.HandleFunc("/user/{id}", uc.DeleteUser).Methods("DELETE")
}
