package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gonzalogorgojo/go-caching-proxy/internal/cache"
	"github.com/gonzalogorgojo/go-caching-proxy/internal/db"
	"github.com/gonzalogorgojo/go-caching-proxy/internal/proxy"
)

func main() {
	port := flag.Int("port", 0, "--port is the port on which the caching proxy server will run. It needs to be a int")
	target := flag.String("target", "", "--target is the URL of the server to which the requests will be forwarded. It needs to be a valid URL")
	cleanInterval := flag.Int64("clean", 0, "--clean is the number in minutes that the cleanup service will run interval. It needs to be a int")
	clear := flag.Bool("clear-cache", false, "--clear-cache tells the program to clear the cache inmediatly.")

	flag.Parse()

	if *clear {
		clearCacheCommand()
		return
	}

	if *port == 0 || *target == "" {
		fmt.Println("Usage: caching-proxy --port <number> --target <url>")
		os.Exit(1)
	}

	startServer(*port, *target, *cleanInterval)
}

func startServer(port int, target string, cleanInterval int64) {
	database := db.InitDB()
	defer database.Close()

	c := &cache.Cache{
		DB: database,
	}

	if cleanInterval > 0 {
		fmt.Printf("Starting cache cleaning service that will run every %v minutes\n", cleanInterval)
		c.CleanUpService(time.Duration(cleanInterval) * time.Minute)
	}

	mux := http.NewServeMux()

	proxyHandler, err := proxy.ProxyHandler(target, c)
	if err != nil {
		log.Fatalf("Error creating proxy handler: %v", err)
	}

	mux.HandleFunc("/", proxyHandler)

	srv := &http.Server{
		Handler:      mux,
		Addr:         ":" + strconv.Itoa(port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Server started at http://localhost:%v and will forward request to: %v\n", port, target)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func clearCacheCommand() {
	database := db.InitDB()
	defer database.Close()

	c := &cache.Cache{DB: database}

	err := c.ClearCache()
	if err != nil {
		log.Fatalf("Failed clear cache: %v", err)
	}
	fmt.Println("Cache cleared successfully.")
}
