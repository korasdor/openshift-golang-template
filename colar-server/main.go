package main

import (
	"github.com/gorilla/mux"

	"fmt"
	"net/http"
)

func main() {


	//http.HandleFunc("/", hello)
	//
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	//
	//err := http.ListenAndServe(":8080", nil)
	//if err != nil {
	//	panic(err)
	//}

	r := mux.NewRouter()
	r.HandleFunc("/", hello)
	//r.HandleFunc("/static/media/{file}", handler.ServeStaticFiles)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "hello, world from %s", "s")
}
