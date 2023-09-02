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

//type Book struct {
//	Title  string `json:"title"`
//	Author string `json:"author"`
//}

type Book struct {
	Title  string
	Author string
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	//func getBooks() http.HandlerFunc {
	//return func(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	//var books []Book
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

	//if err := cursor.All(ctx, &books); err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, books := range books {
	//	fmt.Fprintf(w, "Title: %s\nAuthor: %s\n", books.Title, books.Author)
	//}

	//var book Book
	//
	//filter := bson.D{{"title", "A"}}
	//
	//err := collection.FindOne(ctx, filter).Decode(&book)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Fprintf(w, "Title: %s\nAuthor: %s\n", book.Title, book.Author)
}

//func getBooks(w http.ResponseWriter, r *http.Request) {
//	//func getBooks() http.HandlerFunc {
//	//return func(w http.ResponseWriter, r *http.Request) {
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	var book Book
//
//	filter := bson.D{{"title", "A"}}
//
//	err := collection.FindOne(ctx, filter).Decode(&book)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	//fmt.Printf("Title: %s\nAuthor: %s\n", book.Title, book.Author)
//	fmt.Fprintf(w, "Title: %s\nAuthor: %s\n", book.Title, book.Author)
//	//}
//}

func useRoute(router *mux.Router) {
	router.HandleFunc("/", homepage)
	router.HandleFunc("/books", getBooks).Methods("GET")
}

//func homepage(w http.ResponseWriter, r *http.Client) {
//	fmt.Fprintf(w, "Welcome")
//	fmt.Println("Hompeage")
//}
//
//func getBooks(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(books)
//}
//
//type Book struct {
//	Title  string `json:"title"`
//	Author string `json:"author"`
//}

//type UserResponce struct {
//	Status  int                    `json:"status"`
//	Message string                 `json:"message"`
//	Data    map[string]interface{} `json:"data"`
//}
//
//type User struct {
//	Id       primitive.ObjectID `json:"id,omitempty"`
//	Name     string             `json:"name,omitempty" validate:"required"`
//	Location string             `json:"location,omitempty" validate:"required"`
//	Title    string             `json:"title,omitempty" validate:"required"`
//}

func main() {

	router := mux.NewRouter()

	useRoute(router)

	log.Fatal(http.ListenAndServe(":8080", router))

}
