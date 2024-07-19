package main

import (
	"log"

	server "github.com/alexPavlikov/gora_driver_location_service/cmd"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
