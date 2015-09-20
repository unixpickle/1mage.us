package imagedb

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"time"

	"github.com/nfnt/resize"
)

const maxThumbnailDimension = 256

// An Image contains metadata about a single post. It does not store the image data itself.
type Image struct {
	MIME      string `json:"mimeType"`
	Id        int64  `json:"id"`
	Timestamp int64  `json:"timestamp"`

	HasSize bool `json:"has_size"`
	Width   int  `json:"width"`
	Height  int  `json:"height"`

	HasThumbnail    bool `json:"has_thumbnail"`
	ThumbnailWidth  int  `json:"thumb_width"`
	ThumbnailHeight int  `json:"thumb_height"`
}

// processImage reads an image from a temporary file and gleans all the information it can from it.
func processImage(tempFile *os.File, mimeType string) (entry Image, thumbnail *os.File, err error) {
	entry.MIME = mimeType
	entry.Timestamp = time.Now().Unix()

	if _, err = tempFile.Seek(0, 0); err != nil {
		return
	}

	parsedImage, _, err := image.Decode(tempFile)
	if err != nil {
		err = nil
		return
	}

	entry.HasSize = true
	entry.Width = parsedImage.Bounds().Dx()
	entry.Height = parsedImage.Bounds().Dy()
	entry.HasThumbnail = true

	thumbnail, err = ioutil.TempFile("", "1mage_thumbnail")
	if err != nil {
		return
	}

	var thumbnailData []byte
	thumbnailData, entry.ThumbnailWidth, entry.ThumbnailHeight = makeThumbnail(parsedImage)
	if _, err = thumbnail.Write(thumbnailData); err != nil {
		thumbnail.Close()
		os.Remove(thumbnail.Name())
	}

	return
}

// makeThumbnail generates a small thumbnail from an image.
// It returns the thumbnail's dimensions and PNG representation.
func makeThumbnail(orig image.Image) (data []byte, width, height int) {
	origWidth := orig.Bounds().Dx()
	origHeight := orig.Bounds().Dy()

	var newImage image.Image
	if origWidth <= maxThumbnailDimension && origHeight <= maxThumbnailDimension {
		newImage = orig
	} else if origWidth > origHeight {
		newImage = resize.Resize(maxThumbnailDimension, 0, orig, resize.Lanczos3)
	} else {
		newImage = resize.Resize(0, maxThumbnailDimension, orig, resize.Lanczos3)
	}

	width = newImage.Bounds().Dx()
	height = newImage.Bounds().Dy()

	var buff bytes.Buffer
	if err := png.Encode(&buff, newImage); err != nil {
		panic("png.Encode() should not have encountered I/O failure")
	}

	data = buff.Bytes()
	return
}
