package main

import (
	"fmt"
	"net/http"
)

func main() {


	http.HandleFunc("/", hello)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	
    bind := fmt.Sprintf("%s:%s", "http://0.0.0.0", "8080")
	fmt.Printf("listening on %s...", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "hello, world from %s", "s")
}
