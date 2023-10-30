package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func checkHTTPStatus(domain string, port string, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{
		Timeout: timeout,
	}
	url := "http://" + domain + ":" + port

	res, err := client.Get(url)
	if err != nil {
		// http not found, check with https
		url = "https://" + domain + ":" + port
		res, err = client.Get(url)
		if err != nil {
			return
		}
	}

	defer res.Body.Close()

	if res.StatusCode == 200 || res.StatusCode == 202 || res.StatusCode == 204 || res.StatusCode == 301 || res.StatusCode == 302 {
		fmt.Printf("%s\n", url)
	}
}

func main() {
	timeout := flag.Duration("t", 5*time.Second, "Timeout for HTTP requests. Default 5s")
	port := flag.String("p", "80,443", "Ports to check (comma-separated). Default is 80,443")
	flag.Parse()

	ports := strings.Split(*port, ",")
	scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
	var wg sync.WaitGroup

	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		for _, p := range ports {
			wg.Add(1)
			go checkHTTPStatus(domain, p, *timeout, &wg)
		}
	}

	wg.Wait()
}
