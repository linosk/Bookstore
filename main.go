package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	uri = "mongodb://localhost:27017"
)

type Book struct {
	Title  string
	Author string
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
	fmt.Println("Hompage endpoint")
}

func list(w http.ResponseWriter, r *http.Request) {

}

func handleRequsts(ctx context.Context, collection *mongo.Collection) {

	r := mux.NewRouter()
	r.HandleFunc("/", homepage)
	r.HandleFunc("/list", list)

	http.ListenAndServe(":8080", r)
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("bookstore")
	collection := database.Collection("books")

	handleRequsts(ctx, collection)
	//POST
	//GET
	//GET
	//PUT
	//PATCH
	//DELETE
}
