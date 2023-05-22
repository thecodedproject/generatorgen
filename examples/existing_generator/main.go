package main

import (
	internal "github.com/thecodedproject/generatorgen/examples/existing_generator/internal"
	log "log"
)

func main() {

	// Some custom main func gubbins that shouldn't be touched by the generator...

	err := internal.Generate()
	if err != nil {
		log.Fatal(err.Error())
	}
}

