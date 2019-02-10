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

// This flag keeps the maximum number of parallel HTTP requests
var parallel int

func init() {
	flag.IntVar(&parallel, "parallel", 10, "Maximum number of parallel requests")
}

func main() {

	if len(os.Args) == 1 {
		showHelp()
		return
	}

	if len(os.Args) == 2 && os.Args[1] == "--help" {
		showHelp()
		return
	}

	flag.Parse()
	if flag.NArg() == 0 {
		showHelp()
		return
	}

	args := flag.Args()

	// Extract urls from input
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

	// Listen to SIGINT (Ctrl + C) signals to cancel in-flight requests
	// and shutdown gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ctxCf := make(chan context.CancelFunc, len(urls))

	go func() {
		<-sig
		// Cancel all in-flight requests
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
			// Buffer all cancel funcs to have the possibility of terminating in-flight requests
			ctx, cancel := context.WithCancel(req.Context())
			req = req.WithContext(ctx)
			ctxCf <- cancel
			hash, err := httpsender.Make(c, req)
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

func showHelp() {
	help := `Myhttp is a tool that fetches the content of multiple urls and prints each one of them on a separate line 
along with the corresponding MD5 hash of the HTML content located at the current url.
			
Usage:
			
$> ./myhttp golang.org google.com
http://www.golang.org d1b40e2a2ba488a054186e4ed0733f9752f66949
http://google.com 9d8ec921bdd275fb2a605176582e08758eb60641

Use the -parallel flag to control the maximum number of parallel HTTP requests

$> ./myhttp -parallel 3 golang.org google.com
http://www.golang.org d1b40e2a2ba488a054186e4ed0733f9752f66949
http://google.com 9d8ec921bdd275fb2a605176582e08758eb60641`

	fmt.Println(help)
}
