package thumbnail

import (
	"bytes"
	"image"
	"image/png"

	"github.com/nfnt/resize"
)

// MaxThumbnailDimension is the maximum width and height for thumbnail images.
const MaxThumbnailDimension = 256

// MakeThumbnail generates a small thumbnail from an image. It returns the thumbnail's dimensions
// and PNG representation.
func MakeThumbnail(orig image.Image) (data []byte, width, height int) {
	origWidth := orig.Bounds().Dx()
	origHeight := orig.Bounds().Dy()

	var newImage image.Image
	if origWidth <= MaxThumbnailDimension && origHeight <= MaxThumbnailDimension {
		newImage = orig
	} else if origWidth > origHeight {
		newImage = resize.Resize(MaxThumbnailDimension, 0, orig, resize.Lanczos3)
	} else {
		newImage = resize.Resize(0, MaxThumbnailDimension, orig, resize.Lanczos3)
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
