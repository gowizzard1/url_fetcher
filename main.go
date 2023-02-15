package main

import (
	"flag"
	"fmt"
	"time"
	client2 "url_fetcher/http/client"
	"url_fetcher/interfaces/client"
	"url_fetcher/utils"
)

var (
	numWorkers  = flag.Int("workers", 10, "Number of worker goroutines")
	timeout     = flag.Duration("timeout", 5*time.Second, "HTTP request timeout")
	maxIdle     = flag.Int("max-idle", 100, "Maximum number of idle connections")
	idleTimeout = flag.Duration("idle-timeout", 30*time.Second, "Idle connection timeout")
)

func main() {
	flag.Parse()
	urls := utils.ParseFlagsAndArgsToURL()

	client := client2.NewClient(numWorkers, maxIdle, timeout, idleTimeout)

	urlResponses := client.Request(urls)
	printResults(urlResponses)
}

func printResults(urlResponses []client.URLResponse) {
	for _, r := range urlResponses {
		fmt.Printf("%s %d\n", r.URL, r.BodySize)
	}
}
