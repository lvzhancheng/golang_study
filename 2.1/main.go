package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// v3
func main() {
	// defer fmt.Println("panic退出前处理")
	mux := http.NewServeMux()
	mux.Handle("/", &myHandler{})
	mux.HandleFunc("/version", version)
	mux.HandleFunc("/healthZ", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "200")
		// log.Println(r.Host, r.Method, r.Response.StatusCode)
	})

	server := &http.Server{
		Addr:         ":80",
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	log.Println("Starting v3 httpserver")
	log.Fatal(server.ListenAndServe())
}

type myHandler struct {
}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is version 3"))
	log.Println(r.Host, r.URL, r.Method)
}
func version(w http.ResponseWriter, r *http.Request) {
	v, exists := os.LookupEnv("VERSION")
	if exists {
		w.Header().Add("version", v)
	} else {
		w.Header().Add("version", "NULL")
	}
	fmt.Fprintf(w, os.Getenv("GOOS"))
	log.Println(r.Host, r.URL, r.Method)
}
