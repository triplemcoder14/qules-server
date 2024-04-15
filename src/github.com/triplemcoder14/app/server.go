package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Is request, Coming through?")
	})

	// Add a health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	log.Println("Starting server...")
	l, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		log.Fatal(http.Serve(l, nil))
	}()

	log.Println("Waiting for server to start...")
	if err := waitForServerStart("http://localhost:8081/health"); err != nil {
		log.Fatal(err)
	}

	log.Println("Server is ready, sending request...")
	res, err := http.Get("http://localhost:8081/hello")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Reading response...")
	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		log.Fatal(err)
	}
} 	

func waitForServerStart(url string) error {
	for i := 0; i < 10; i++ {
		time.Sleep(500 * time.Millisecond)
		_, err := http.Get(url)
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("server did not start within expected time")
}

