package main

import (
	"context"
	"flag"
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
		log.Println("Listening on port:", port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln("ListenAndServe:", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(stop, os.Interrupt)

	// Block until we receive signal
	<-stop

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
		log.Println("shutting down...")
		err := srv.Shutdown(ctx)
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln("Shutdown:", err)
		}
		log.Println("shutdown complete")
	}()

	wg.Wait()
	os.Exit(0)
}
