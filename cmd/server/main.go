package main

import (
	"Log"
	"github.com/soichiro0915/proglog/internal/server"
)

func main(){
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}