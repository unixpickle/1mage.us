package main

import (
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reader, err := r.MultipartReader()
	if err != nil {
		sendUploadError(w, r, "multipart error: "+err.Error())
		return
	}

	ids := []int{}
	for {
		part, err := reader.NextPart()
		if err != nil {
			if err != io.EOF {
				sendUploadError(w, r, "error reading part: "+err.Error())
				return
			}
			break
		}

		if imageId, err := uploadImage(part); err != nil {
			sendUploadError(w, r, "error uploading image: "+err.Error())
			return
		} else {
			ids = append(ids, imageId)
		}
	}

	responseMap := map[string]interface{}{"ids": ids}
	if len(ids) == 1 {
		responseMap["identifier"] = ids[0]
	}
	data, _ := json.Marshal(responseMap)
	w.Write(data)
}

func sendUploadError(w http.ResponseWriter, r *http.Request, errStr string) {
	log.Print("Error from " + r.RemoteAddr + ": " + errStr)
	packet := map[string]string{"error": errStr}
	data, _ := json.Marshal(&packet)
	w.Write(data)
}

func uploadImage(part *multipart.Part) (int, error) {
	// TODO: read the image, generate the thumbnail, add it to the DB, and return the new ID.
	return 0, nil
}
