package main

import (
	"log"
	"net/http"
)

func main() {
	//fs := http.FileServer(http.Dir("views"))
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	log.Println("Listening...")
	http.ListenAndServe("127.0.0.1:8080", nil)
}
