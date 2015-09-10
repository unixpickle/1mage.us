package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

var ErrRateLimit error = errors.New("rate limit exceeded")

type UploadResult struct {
	Error       *string `json:"error"`
	RateLimited bool    `json:"rate_limited"`
	Id          *int    `json:"id"`
}

func NewUploadResultSuccess(id int) UploadResult {
	return UploadResult{Id: &id}
}

func NewUploadResultError(err error) UploadResult {
	str := err.Error()
	return UploadResult{Error: &str, RateLimited: err == ErrRateLimit}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reader, err := r.MultipartReader()
	if err != nil {
		packet := map[string]string{"error": "multipart error: " + err.Error()}
		data, _ := json.Marshal(&packet)
		w.Write(data)
		return
	}

	var results []UploadResult
	for {
		if part, err := reader.NextPart(); err == nil {
			if imageId, err := uploadImage(part, r); err != nil {
				results = append(results, NewUploadResultError(err))
			} else {
				results = append(results, NewUploadResultSuccess(imageId))
			}
		} else if err != io.EOF {
			packet := map[string]string{"error": "multipart error: " + err.Error()}
			data, _ := json.Marshal(&packet)
			w.Write(data)
			return
		} else {
			break
		}
	}

	responseMap := map[string]interface{}{"results": results}

	// For legacy purposes, set a global "error" or "identifier" field if applicable.
	if len(results) == 1 {
		result := results[0]
		if result.Error != nil {
			responseMap["error"] = *result.Error
		} else {
			responseMap["identifier"] = *result.Id
		}
	}

	data, _ := json.Marshal(responseMap)
	w.Write(data)
}

func mimeTypeForPart(part *multipart.Part) string {
	// TODO: check the multipart header to see if the browser provides a MIME type
	ext := path.Ext(part.FileName())
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "image/jpeg"
	}
	return mimeType
}

func sendUploadError(w http.ResponseWriter, r *http.Request, errStr string) {
	log.Print("Error from " + r.RemoteAddr + ": " + errStr)
	packet := map[string]string{"error": errStr}
	data, _ := json.Marshal(&packet)
	w.Write(data)
}

func uploadImage(part *multipart.Part, r *http.Request) (id int, err error) {
	if !RateLimitRequest(r) {
		err = ErrRateLimit
		return
	}

	tempFile, err := ioutil.TempFile("", "1mage")
	if err != nil {
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	cappedReader := NewCappedReader(part)
	if _, err := io.Copy(tempFile, cappedReader); err != nil {
		return 0, err
	}

	// TODO: read image file to generate thumbnail.
	// TODO: write thumbnail file.
	// TODO: generate entry in database.

	return 0, nil
}
