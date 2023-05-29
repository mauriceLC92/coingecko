package main

import (
	"fmt"
	"log"

	"github.com/mauriceLC92/coingecko"
)

func main() {
	client := coingecko.Client{
		Endpoint: coingecko.DefaultEndpoint,
	}

	ping, err := client.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", ping)
}
