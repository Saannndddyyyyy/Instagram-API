package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

var client1 *mongo.Client

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	name string             `json:"name,omitempty" bson:"name,omitempty"`
	pwd  string             `json:"pwd,omitempty" bson:"pwd,omitempty"`
}
func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var user User
	_ = json.NewDecoder(request.Body).Decode(&user)
	collection := client1.Database("test").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)


}

func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User
	collection := client1.Database("test").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func main() {
	fmt.Println("Starting the application...")

	clientOptions1 := options.Client().ApplyURI("mongodb://localhost:27017")
	fmt.Println("Hello, World! 1")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client1, _ = mongo.Connect(ctx, clientOptions1)
	fmt.Println("Hello, World! 2")
	router := mux.NewRouter()
	fmt.Println("Hello, World! 3")
	router.HandleFunc("/user", CreateUserEndpoint).Methods("POST")
	fmt.Println("Hello, World! 4")
	router.HandleFunc("/person/{id}", GetUserEndpoint).Methods("GET")
	fmt.Println("Hello, World! 5")
}