package main

import (
	"context"
	"log"
	"microservices/handlers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	l := log.New(os.Stdout, "microservices ", log.LstdFlags)

	h := handlers.NewHello(l)
	p := handlers.NewProducts(l)

	mux := http.NewServeMux()
	mux.Handle("/", h)
	mux.Handle("/products/", p)

	server := http.Server{
		Addr:         ":9090",           // configure the bind address
		Handler:      mux,               // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	go func() {
		l.Println("Starting the server")

		err := server.ListenAndServe()

		if err != nil {
			l.Printf("Error starting server: %s", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	sig := <-c
	l.Printf("Got signal: %s", sig.String())


	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
