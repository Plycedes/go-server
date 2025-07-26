package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/plycedes/mongoAPI/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://blacksteel:napoleon@virgo0.ltm8rd6.mongodb.net"
const dbName = "goflixDB"
const colName = "watchlist"

var collection *mongo.Collection

func init(){
	clientOption := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOption)
	handleErr(err)

	fmt.Println("MongoDB connected")

	collection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection instance is ready")
}

func handleErr(err error){
	if err != nil {
		log.Fatal(err)
	}
}

func insertOneMovie(movie models.Netflix) *mongo.InsertOneResult {
	inserted, err := collection.InsertOne(context.Background(), movie)
	handleErr(err)

	fmt.Println("Inserted 1 movie in the db with id:", inserted.InsertedID)
	return inserted
}

func updateOneMovie(movieId string){
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	res, err := collection.UpdateOne(context.Background(), filter, update)
	handleErr(err)

	fmt.Println("Modified count:", res.ModifiedCount)
}

func deleteOneMovie(movieId string){
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	handleErr(err)

	fmt.Println("Movies deleted:", deleteCount)
}

func deleteAllMovies() int64 {
	res, err := collection.DeleteMany(context.Background(), bson.M{})
	handleErr(err)
	fmt.Println("Movies deleted", res.DeletedCount)
	return res.DeletedCount
}

func getAllMovies() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	handleErr(err)

	var movies []primitive.M

	for cur.Next(context.Background()){
		var movie bson.M
		err := cur.Decode(&movie)
		handleErr(err)
		movies = append(movies, movie)
	}

	defer cur.Close(context.Background())
	return movies
}

func getOneMovie(movieId string) primitive.M {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	cur, err := collection.Find(context.Background(), filter)
	handleErr(err)

	var movie primitive.M
	for cur.Next(context.Background()){
		err := cur.Decode(&movie)
		handleErr(err)
	}

	return movie
}

func GetAllMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func GetOneMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "x-www-form-urlencoded")

	params := mux.Vars(r)
	res := getOneMovie(params["id"])

	json.NewEncoder(w).Encode(res)
}

func CreateMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie models.Netflix
	err := json.NewDecoder(r.Body).Decode(&movie)
	handleErr(err)

	res := insertOneMovie(movie)
	json.NewEncoder(w).Encode(res)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT");

	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode("Movie marked as watched")
}

func DeleteOneMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-from-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])

	json.NewEncoder(w).Encode("Movie deleted successfully")
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	count := deleteAllMovies()
	json.NewEncoder(w).Encode("Deleted all movies successfully, count: "+ strconv.FormatInt(count, 10))
}

func Status(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	json.NewEncoder(w).Encode("Welcome to PlyAPI")
}