package main

import (
	"fmt"
	"log"
	"net/http"
)

const destinationUrl = "localhost:8080"

func StartDestinationServer() {
	go func() {
		router := http.NewServeMux()
		router.HandleFunc("/", Home)
		fmt.Println("listening to port: " + destinationUrl)
		log.Fatalln(http.ListenAndServe(destinationUrl, router))
	}()
}

func Home(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("YES, destination is available")
}
