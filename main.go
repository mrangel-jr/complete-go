package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/mrangel-jr/complete-go/internals/app"
	"github.com/mrangel-jr/complete-go/internals/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Go backend server port")
	flag.Parse()
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	app.Logger.Printf("We are running out app on port %d", port)

	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
