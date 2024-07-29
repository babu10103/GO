package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/babu10103/mongo-golang/models"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mgo.Session
}

func NewUserController(session *mgo.Session) *UserController {
	return &UserController{session}
}

// GetUser retrieves a user by ID.
func (controller UserController) GetUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")

	if !bson.IsObjectIdHex(id) {
		http.NotFound(w, r)
		return
	}

	oid := bson.ObjectIdHex(id)

	user := models.User{}

	if err := controller.session.DB("mongo-golang").C("users").FindId(oid).One(&user); err != nil {
		http.NotFound(w, r)
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newUser models.User

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	newUser.Id = bson.NewObjectId()

	if err := uc.session.DB("mongo-golang").C("users").Insert(&newUser); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userJSON, err := json.Marshal(newUser)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userJSON)
}

// DeleteUser removes a user with the given ID.
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id) {
		http.NotFound(w, r)
		return
	}

	oid := bson.ObjectIdHex(id)
	err := uc.session.DB("mongo-golang").C("users").RemoveId(oid)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted user %s", id)
}
