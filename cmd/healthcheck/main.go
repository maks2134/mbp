package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	url := "http://localhost:8000/api/posts"

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Успешный статус
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "Received non-2xx status code: %d\n", resp.StatusCode)
	os.Exit(1)
}
