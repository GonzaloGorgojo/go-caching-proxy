package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {

	PORT := flag.Int("port", 0, "--port is the port on which the caching proxy server will run. It needs to be a int")
	ORIGIN := flag.String("origin", "", "--origin is the URL of the server to which the requests will be forwarded. It needs to be a valid URL")

	flag.Parse()

	requiredFlags := []string{"port", "origin"}
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range requiredFlags {
		if !seen[req] {
			log.Fatalf("Missing required -%s flag\n", req)
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello world")
	})

	srv := &http.Server{
		Handler:      mux,
		Addr:         ":" + strconv.Itoa(*PORT),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Server will start at http://localhost:%v\n", *PORT)
	fmt.Printf("Server will redirect requests to %v\n", *ORIGIN)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
