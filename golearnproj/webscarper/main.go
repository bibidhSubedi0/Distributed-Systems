/* Things to learn
Goroutines	One scraper per site (non-blocking)
Channels	Collect titles/errors as they finish
HTTP	Built-in GET requests
HTML parsing	Parsing and traversing the DOM
Structs/slices	Managing data
Error handling	Idiomatic if err != nil patterns
*/

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// use capital letter for public, lowercase for private
type Results struct {
	URL   string
	Title string
	Error error
}

func projcore() {
	urls := []string{
		"https://golang.org",
		"https://github.com",
		"https://news.ycombinator.com",
		// "https://invalid.url",
	}

	// A channel is like a thread-safe pipe where we can send and receive value
	results := make(chan Results)

	/*
		A goroutine is a lightweight thread managed by Go’s runtime.
		A channel is a safe way for goroutines to communicate.
	*/

	for _, url := range urls {
		go func(u string) {
			result := Scrape(u)
			results <- result // send result to the channel
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		result := <-results // receive from the channel
		fmt.Printf("[%s] → \"%s\"\n", result.URL, result.Title)
	}

	println("We dont reach here until, we get all the damn titles huh")
}

func Scrape(url string) Results {
	res := Results{URL: url}

	rsp, err := http.Get(url)
	// fmt.Println(rsp.StatusCode)
	if err != nil {
		log.Println("Error fetching the URL:", err)
		return res
	}
	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		log.Println("Error reading the response body:", err)
		return res
	}

	bodyStr := string(body)
	titleStart := strings.Index(bodyStr, "<title>")
	titleEnd := strings.Index(bodyStr, "</title>")

	if titleStart != -1 && titleEnd != -1 {
		res.Title = bodyStr[titleStart+7 : titleEnd]
	} else {
		res.Title = "No title found"
	}

	return res
}
func main() {
	projcore()
}
