package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
	"workout/internal/app"
	"workout/internal/routes"
)

func main() {
	app, err := app.NewApplication()
	var port int
	flag.IntVar(&port, "port", 8080, "go backend service port")
	flag.Parse()
	if err != nil {
		flag.Parse()
		panic(err)

	}
	defer app.DB.Close()
	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  60 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("We are running on port %d\n", port)
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}

}
