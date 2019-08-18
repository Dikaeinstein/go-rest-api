package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/dikaeinstein/go-rest-api/model"
	"github.com/dikaeinstein/go-rest-api/route"
)

var wait = flag.Duration("graceful-timeout", time.Second*3, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")

func main() {
	flag.Parse()
	d := model.GetDB()
	defer d.Close() // Close db connection

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	r := mux.NewRouter()
	route := route.New(r)

	addr := ":" + port
	srv := http.Server{
		Addr:    addr,
		Handler: route.SetupRoutes(),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	var wg sync.WaitGroup
	wg.Add(2)
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		defer wg.Done()
		fmt.Println("Listening on port:", port)
		err := srv.ListenAndServe()
		if err != nil {
			log.Println("ListenAndServe:", err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive signal
	<-c

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), *wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	go func() {
		defer wg.Done()
		<-ctx.Done()
		fmt.Println()
		log.Println("shutting down...")
		srv.Shutdown(ctx)
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done) // Signal done channel
	}()
	<-done
	os.Exit(0)
}
