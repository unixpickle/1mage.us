package main

import (
	"mime"
	"net/http"
	"path"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func ServeAssetForRequest(w http.ResponseWriter, r *http.Request) {
	urlPath := path.Clean(r.URL.Path)
	if data, err := Asset("assets/" + urlPath[1:]); err != nil {
		NotFound(w, r)
	} else {
		mimeType := mime.TypeByExtension(path.Ext(urlPath))
		if mimeType == "" {
			mimeType = "text/plain"
		}
		w.Header().Set("Content-Type", mimeType)
		w.Write(data)
	}
}

func ServePage(w http.ResponseWriter, assetName string) {
	if data, err := Asset("assets/" + assetName); err != nil {
		panic("page not found: " + assetName)
	} else {
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	}
}
