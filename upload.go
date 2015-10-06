package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/unixpickle/1mage.us/imagedb"
)

var RateLimitErr = errors.New("rate limit exceeded")

// An UploadResult is sent to clients whenever they use the upload API.
type UploadResult struct {
	GlobalError *string `json:"error"`

	Results []ImageOrError `json:"results"`

	// Identifier exists for legacy support and has no meaning for multi-file uploads.
	Identifier *int64 `json:"identifier"`
}

type ImageOrError struct {
	Image *imagedb.Image `json:"image"`
	Error *string        `json:"error"`
}

func ServeUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reader, err := r.MultipartReader()
	if err != nil {
		errorMessage := err.Error()
		data, _ := json.Marshal(UploadResult{GlobalError: &errorMessage})
		w.Write(data)
		return
	}

	var result UploadResult
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			errorMessage := err.Error()
			result.Results = append(result.Results, ImageOrError{nil, &errorMessage})
			continue
		}
		image, err := uploadImage(part, r)
		if err != nil {
			errorMessage := err.Error()
			result.Results = append(result.Results, ImageOrError{nil, &errorMessage})
		} else {
			result.Results = append(result.Results, ImageOrError{image, nil})
		}
	}

	// TODO: here, emit that a change was made to all connected clients.
	if len(result.Results) == 1 && result.Results[0].Image != nil {
		result.Identifier = &result.Results[0].Image.Id
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}

func uploadImage(p *multipart.Part, req *http.Request) (*imagedb.Image, error) {
	if RateLimitRequest(req) {
		return nil, RateLimitErr
	}

	// NOTE: we do not obtain the ShutdownLock here. That is because the uploader controls the speed
	// of the upload and we don't want them to hold us back from shutting down.
	reader := NewCappedReader(p)
	tempFile, err := ioutil.TempFile(TemporaryDirectory, "upload_image")
	if err != nil {
		return nil, err
	} else if _, err := io.Copy(tempFile, reader); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, err
	}

	mimeType := mimeTypeForPart(p)

	ShutdownLock.RLock()
	defer ShutdownLock.RUnlock()
	entry, _, err := GlobalDb.Add(tempFile, mimeType)

	return &entry, err
}

func mimeTypeForPart(part *multipart.Part) string {
	if contentType := part.Header.Get("Content-Type"); contentType != "" {
		if _, _, err := mime.ParseMediaType(contentType); err == nil {
			return contentType
		}
	}
	ext := path.Ext(part.FileName())
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "image/jpeg"
	}
	return mimeType
}
