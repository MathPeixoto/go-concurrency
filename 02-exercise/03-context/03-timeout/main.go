package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	// TODO: set a http client timeout

	req, err := http.NewRequest("GET", "https://andcloud.io", nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 500*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	// Close the response body on the return.
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("ERROR:", err)
		}
	}(resp.Body)

	// Write the response to stdout.
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return
	}
}
