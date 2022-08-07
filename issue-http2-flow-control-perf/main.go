package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func runServer() {
	h2server := &http2.Server{}
	server := &http.Server{
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			// Simulate slow handler by sleeping for non-priority requests.
			if r.Header.Get("request-type") != "priority" {
				time.Sleep(time.Second * 5)
				return
			}

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			_, _ = w.Write(b)
		}), h2server),
		Addr: ":8081",
	}

	fmt.Printf("H2c Server starting on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func runClient() {
	client := http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}

	payload := generatePayload(1 << 21) // 2MB payload.
	var wg sync.WaitGroup

	// Low priority request
	wg.Add(1)
	go func() {
		req, err := http.NewRequest("POST", "http://localhost:8081", bytes.NewBuffer(payload))
		if err != nil {
			panic(err)
		}

		start := time.Now()
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println("low priority req response time:", time.Since(start))
		wg.Done()
	}()

	// High priority request must finish asap
	wg.Add(1)
	go func() {
		req, err := http.NewRequest("POST", "http://localhost:8081", bytes.NewBuffer(payload))
		if err != nil {
			panic(err)
		}
		req.Header.Add("request-type", "priority")

		start := time.Now()
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("high priority req response time:", time.Since(start))
		wg.Done()
	}()

	wg.Wait()
}

func generatePayload(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = 'A'
	}

	return b
}

func main() {
	go runServer()
	runClient()
}
