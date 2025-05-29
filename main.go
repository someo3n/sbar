package main

import (
	"log"
)

func main() {
	bar, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	bar.Loop()
}
