package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: 1mage.us <port> <data dir>")
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			ServeAssetForRequest(w, r)
		} else {
			ServePage(w, "upload.html")
		}
	})

	// TODO: catch termination interrupt and shutdown gracefully.

	if err := http.ListenAndServe(":"+os.Args[1], nil); err != nil {
		fmt.Fprintln(os.Stderr, "Error listening and serving:", err)
	}
}
