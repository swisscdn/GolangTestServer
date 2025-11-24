package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Get the TEST environment variable
	testEnv := os.Getenv("TEST")

	// If not set, return a default message
	if testEnv == "" {
		testEnv = "TEST environment variable is not set"
		log.Println("TEST environment variable is not set")
	} else {
		log.Println("TEST variable is:", testEnv)
	}

	// Set response headers
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Write the response
	fmt.Fprintf(w, testEnv)
}
