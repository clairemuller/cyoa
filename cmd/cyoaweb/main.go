package main

import (
	"flag"
	"fmt"
	"gophercises/cyoa"
	"os"
)

func main() {
	filename := flag.String("filename", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	f, err := os.Open(*filename)
	check(err)

	story, err := cyoa.JsonStory(f)
	check(err)

	fmt.Printf("%+v\n", story)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// NOTES

// json.NewDecoder()
// pass in an io.Reader, whereas with json.Marshal you pass in a byte slice
