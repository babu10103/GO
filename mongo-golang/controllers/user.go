package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/babu10103/mongo-golang/models"
	"github.com/babu10103/mongo-golang/utils"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	client *mongo.Client
}

func NewUserController(client *mongo.Client) *UserController {
	return &UserController{client}
}

func (uc UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	collection := uc.client.Database("mongo-golang").Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer cursor.Close(ctx)

	var users []models.User

	if err = cursor.All(ctx, &users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usersJson, err := json.Marshal(users)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usersJson)

}

// GetUser retrieves a user by ID.
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	oid, err := utils.GenerateObjectIdFromHex(id)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	collection := uc.client.Database("mongo-golang").Collection("users")

	user := models.User{}

	err = collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	userJson, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInsufficientStorage)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Id = utils.GenerateObjectId()

	collection := uc.client.Database("mongo-golang").Collection("users")

	_, err := collection.InsertOne(context.Background(), user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeleteUser removes a user with the given ID.
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	oid, err := utils.GenerateObjectIdFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	collection := uc.client.Database("mongo-golang").Collection("users")
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": oid})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted user %s", id)
}
