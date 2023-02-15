package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"sync"
	"time"
	client2 "url_fetcher/interfaces/client"
)

type cli struct {
	numWorkers int
	client     http.Client
}

func NewClient(numWorkers, maxIdle *int, timeout, idleTimeout *time.Duration) client2.IClient {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   *timeout,
				KeepAlive: *idleTimeout,
			}
			return dialer.DialContext(ctx, network, addr)
		},
		MaxIdleConns:        *maxIdle,
		MaxIdleConnsPerHost: *maxIdle,
		IdleConnTimeout:     *idleTimeout,
	}
	client := &http.Client{
		Timeout:   *timeout,
		Transport: transport,
	}
	return &cli{
		client:     *client,
		numWorkers: *numWorkers,
	}
}

func (c cli) Request(urls []string) []client2.URLResponse {
	var wg sync.WaitGroup
	responses := make(chan client2.URLResponse, len(urls))
	urlsChan := make(chan string, len(urls))
	for i := 0; i < c.numWorkers; i++ {
		go func() {
			for url := range urlsChan {
				resp, err := c.client.Get(url)
				if err != nil {
					fmt.Printf("Error fetching %s: %s\n", url, err)
					continue
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("Error reading response body from %s: %s\n", url, err)
					continue
				}
				responses <- client2.URLResponse{url, len(body)}
			}
			wg.Done()
		}()
	}
	for _, url := range urls {
		urlsChan <- url
	}
	close(urlsChan)
	wg.Add(c.numWorkers)

	wg.Wait()
	close(responses)

	var urlResponses []client2.URLResponse
	for r := range responses {
		urlResponses = append(urlResponses, r)
	}
	sort.Slice(urlResponses, func(i, j int) bool {
		return urlResponses[i].BodySize < urlResponses[j].BodySize
	})

	return urlResponses
}

/*  go run main.go --timeout=10s https://api.publicapis.org/entries https://catfact.ninja/fact https://api.coindesk.com/v1/bpi/currentprice.json
 */
