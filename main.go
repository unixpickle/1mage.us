package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: 1mage.us <port> <data dir>")
		os.Exit(1)
	}

	var err error
	TemporaryDirectory, err = ioutil.TempDir("", "1mage_temp")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to setup temporary directory")
		os.Exit(1)
	}
	defer os.RemoveAll(TemporaryDirectory)

	GlobalDb, err = SetupDb(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to setup database:", err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("request for: ", r.URL.Path)
		switch r.URL.Path {
		case "/":
			ServePage(w, "upload.html")
		case "/upload":
			ServeUpload(w, r)
		case "/auth":
			ServeAuth(w, r)
		default:
			ServeAssetForRequest(w, r)
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
	ShutdownLock.Lock()
	return
}
