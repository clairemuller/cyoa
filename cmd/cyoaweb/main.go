package main

import (
	"flag"
	"fmt"
	"gophercises/cyoa"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the CYOA web app on")
	filename := flag.String("filename", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()

	fmt.Printf("Using the story in %s.\n", *filename)

	// open the file for reading
	file, err := os.Open(*filename)
	check(err)

	// use the cyoa package to create a new Story map
	story, err := cyoa.JSONStory(file)
	check(err)

	handler := cyoa.NewHandler(story)

	fmt.Printf("Starting the server on port: %d\n", *port)

	// ListenAndServe starts an HTTP server with a given address and handler
	// The handler is usually nil, which means to use DefaultServeMux
	// Here we've created an instance of the Handler interface to which
	// all requests should be sent
	addr := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(addr, handler))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
