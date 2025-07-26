package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/plycedes/go-server/routes"
)

func main() {
	log.Fatal(http.ListenAndServe(":8000", routes.Router()))
	fmt.Println("Server running at port 8000")
}