package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func Greet(writer io.Writer, name string) {
	_,err := fmt.Fprintf(writer, "Hello, %s", name)
	if err != nil {
		fmt.Println(err)
	}
}

func MyGreeterHandler(w http.ResponseWriter, _ *http.Request) {
	Greet(w, "world")
}

func main() {
	log.Fatal(http.ListenAndServe(":5000", http.HandlerFunc(MyGreeterHandler)))
}
