package client

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"sync"
	"time"
	client2 "url_fetcher/interfaces/client"
)

type cli struct {
	numWorkers  int
	timeout     *time.Duration
	idleTimeout *time.Duration
	maxIdle     *int
}

func NewClient(numWorkers, maxIdle *int, timeout, idleTimeout *time.Duration) client2.IClient {

	return &cli{
		numWorkers:  *numWorkers,
		timeout:     timeout,
		idleTimeout: idleTimeout,
		maxIdle:     maxIdle,
	}
}

func (c cli) Request(urls []string) []client2.URLResponse {
	var wg sync.WaitGroup
	responses := make(chan client2.URLResponse, len(urls))
	urlsChan := make(chan string, len(urls))
	bodyPool := sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   *c.timeout,
			KeepAlive: *c.idleTimeout,
		}).DialContext,
		MaxIdleConns:        *c.maxIdle,
		IdleConnTimeout:     *c.idleTimeout,
		DisableCompression:  true,
		DisableKeepAlives:   true,
		TLSHandshakeTimeout: *c.timeout,
	}
	client := &http.Client{
		Timeout:   *c.timeout,
		Transport: transport,
	}
	for i := 0; i < c.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urlsChan {
				resp, err := client.Get(url)
				if err != nil {
					fmt.Printf("Error fetching %s: %s\n", url, err)
					continue
				}
				body := bodyPool.Get().([]byte)
				n, err := io.ReadFull(resp.Body, body)
				if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
					fmt.Printf("Error reading response body from %s: %s\n", url, err)
					resp.Body.Close()
					bodyPool.Put(body)
					continue
				}
				responses <- client2.URLResponse{url, n}
				resp.Body.Close()
				bodyPool.Put(body)
			}
		}()
	}
	for _, url := range urls {
		urlsChan <- url
	}
	close(urlsChan)
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
