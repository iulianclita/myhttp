package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/iulianclita/myhttp/httpsender"
)

var parallel int

func init() {
	flag.IntVar(&parallel, "parallel", 10, "Maximum number of parallel requests")
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	args := flag.Args()

	var urls []string
	for _, arg := range args {
		var url = arg
		if !strings.HasPrefix(arg, "http://") {
			url = fmt.Sprintf("http://%s", arg)
		}
		urls = append(urls, url)
	}

	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))

	sem := make(chan struct{}, parallel)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ctxCf := make(chan context.CancelFunc, len(urls))

	go func() {
		<-sig
		for fn := range ctxCf {
			fn()
		}
	}()

	for _, url := range urls {
		go func(url string) {
			defer func() {
				<-sem
				wg.Done()
			}()
			sem <- struct{}{}
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				log.Fatal("Cannot create request for url", url)
			}

			ctx, cancel := context.WithCancel(req.Context())
			req = req.WithContext(ctx)
			ctxCf <- cancel
			// hash, err := fetch(c, req)
			hash, err := httpsender.SendRequest(c, req)
			if err != nil {
				fmt.Printf("Failed to fetch url %s: %v\n", url, err)
				return
			}

			fmt.Println(url, hash)
		}(url)
	}

	wg.Wait()
	close(ctxCf)
}
