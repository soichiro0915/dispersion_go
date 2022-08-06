package main

import (
	"log"
	"github.com/soichiro0915/dispersion_go/internal/server"
)

func main(){
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}