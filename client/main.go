package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	destinationEndpoint := os.Getenv("DST_ENDPOINT")

	var wg sync.WaitGroup
	concurrencyCtrlCh := make(chan struct{}, 1000)

	ticker := time.NewTicker(10 * time.Second)

	hc := http.Client{
		Timeout: 3 * time.Second,
	}

loop:
	for {
		select {
		case <-ticker.C:
			break loop
		case concurrencyCtrlCh <- struct{}{}:
			wg.Add(1)
		}

		go func() {
			defer func() {
				wg.Done()
				<-concurrencyCtrlCh
			}()

			resp, err := hc.Get(destinationEndpoint)
			if err != nil {
				fmt.Printf("err: %s\n", err)
				return
			}
			defer func() {
				_ = resp.Body.Close()
			}()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("err: %s\n", err)
				return
			}
			fmt.Printf("%s\n", body)
		}()
	}

	wg.Wait()
}
