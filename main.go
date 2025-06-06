package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	registerRoutes()
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
