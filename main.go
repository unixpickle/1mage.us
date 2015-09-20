package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// GlobalChangeLock should be locked for reading whenever a change is being made to the database.
// When the app is shutting down, this will be locked for writing so that no more changes may occur.
var GlobalChangeLock sync.RWMutex

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: 1mage.us <port> <data dir>")
		os.Exit(1)
	}

	var err error
	GlobalDb, err = SetupDb(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to setup database:", err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			ServeAssetForRequest(w, r)
		} else {
			ServePage(w, "upload.html")
		}
	})

	go func() {
		if err := http.ListenAndServe(":"+os.Args[1], nil); err != nil {
			fmt.Fprintln(os.Stderr, "Error listening and serving:", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Print("1mage.us shutting down...")
	GlobalChangeLock.Lock()
	os.Exit(0)
}
