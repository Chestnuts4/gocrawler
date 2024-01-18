package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	rss1, err := os.ReadFile("rss/citrix-adc.rss")
	if err != nil {
		log.Fatalf("Failed to read RSS file: %v", err)
	}

	rss2, err := os.ReadFile("rss/citrix-adc2.rss")
	if err != nil {
		log.Fatalf("Failed to read RSS file: %v", err)
	}

	requestCount := 0
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if requestCount == 0 {
			fmt.Fprintf(w, "%s", rss1)
		} else if requestCount == 1 {
			fmt.Fprintf(w, "%s", rss2)
		} else {
			fmt.Fprintf(w, "%s", rss2)
		}
		requestCount++
	})

	fmt.Println("Server listening on port 58080")
	http.ListenAndServe(":58080", nil)
}
