package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	uri string = "mongodb://localhost:27017"
)

func connectDatabase() *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal(err)
	}

	//_ used when we have to assign something, but not intend on using
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected")

	return client
}

func getCollection(client *mongo.Client, databaseName string, collectionName string) *mongo.Collection {

	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}

var (
	client     = connectDatabase()
	collection = getCollection(client, "bookstore", "books")
)

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome")
	fmt.Println("Hompeage")
}

type Book struct {
	Title  string
	Author string
}

func getBooks(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	findOptions := options.Find()
	cursor, err := collection.Find(ctx, bson.D{{}}, findOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var book Book
		if err := cursor.Decode(&book); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Title: %s\nAuthor: %s\n", book.Title, book.Author)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Books")

}

func getBook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	bookTitle := vars["title"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var book Book

	filter := bson.D{{"title", bookTitle}}

	err := collection.FindOne(ctx, filter).Decode(&book)
	if err != nil {
		fmt.Fprintf(w, "No book titled '%s'", bookTitle)
		//log.Fatal(err)
	} else {
		fmt.Printf("Title: %s\nAuthor: %s\n", book.Title, book.Author)
		fmt.Fprintf(w, "Title: %s\nAuthor: %s\n", book.Title, book.Author)
	}

}

func deleteBook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	bookTitle := vars["title"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var book Book := bson.D{{}}

}

func useRoute(router *mux.Router) {
	router.HandleFunc("/", homepage)
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{title}", getBook).Methods("GET")
	router.HandleFunc("/books/{title}", deleteBook).Methods("DELETE")
}

func main() {

	router := mux.NewRouter()

	useRoute(router)

	log.Fatal(http.ListenAndServe(":8080", router))

}
