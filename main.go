package main

import (
	"fmt"
	"log"
	"net/http"
)

func print(w http.ResponseWriter, r *http.Request) {
	fmt.Println("")
}

func main() {

	http.HandleFunc("/print", print)
	log.Fatal(http.ListenAndServe(":8080", nil))
}