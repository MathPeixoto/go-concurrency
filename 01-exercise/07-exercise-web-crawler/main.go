package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/net/html"
)

var fetched map[string]bool

type results struct {
	url   string
	urls  []string
	err   error
	depth int
}

// Crawl uses findLinks to recursively crawl
// pages starting with url, to a maximum of depth.
//func Crawl(url string, depth int) {
//	// TODO: Fetch URLs in parallel.
//
//	if depth < 0 {
//		return
//	}
//	urls, err := findLinks(url)
//	if err != nil {
//		// fmt.Println(err)
//		return
//	}
//	fmt.Printf("found: %s\n", url)
//	fetched[url] = true
//	for _, u := range urls {
//		if !fetched[u] {
//			Crawl(u, depth-1)
//		}
//	}
//	return
//}

// Crawl uses findLinks to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int) {
	ch := make(chan *results)

	//Fetch URLs in parallel.
	fetch := func(url string, depth int) {
		urls, err := findLinks(url)
		ch <- &results{url, urls, err, depth}
	}

	go fetch(url, depth)
	fmt.Printf("found: %s\n", url)

	fetched[url] = true

	for fetching := 1; fetching > 0; fetching-- {
		result := <-ch

		if result.err != nil {
			fmt.Println(result.err)
			continue
		}
		fmt.Printf("found: %s\n", result.url)

		if result.depth > 0 {
			for _, u := range result.urls {
				if !fetched[u] {
					fetching++
					go fetch(u, result.depth-1)
					fetched[u] = true
				}
			}
		}
	}
	close(ch)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fetched = make(map[string]bool)
	now := time.Now()
	Crawl("http://andcloud.io", 2)
	fmt.Println("time taken:", time.Since(now))
}

func findLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err := resp.Body.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	return visit(nil, doc), nil
}

// visit appends to links each link found in n, and returns the result.
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
