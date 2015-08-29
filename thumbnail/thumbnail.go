package thumbnail

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
)

// MaxThumbnailDimension is the maximum width and height for thumbnail images.
const MaxThumbnailDimension = 256

// MakeThumbnail generates a small thumbnail from an image. It returns the thumbnail's dimensions
// and PNG representation.
func MakeThumbnail(orig image.Image) (data []byte, width, height int) {
	origWidth := orig.Bounds().Dx()
	origHeight := orig.Bounds().Dy()

	if origWidth <= MaxThumbnailDimension && origHeight <= MaxThumbnailDimension {
		width = origWidth
		height = origHeight
	} else if origWidth > origHeight {
		width = MaxThumbnailDimension
		height = int(float64(origHeight) / float64(origWidth) * MaxThumbnailDimension)
	} else {
		height = MaxThumbnailDimension
		width = int(float64(origWidth) / float64(origHeight) * MaxThumbnailDimension)
	}

	bounds := image.Rectangle{image.Pt(0, 0), image.Pt(width, height)}
	thumbnailImage := image.NewRGBA(bounds)
	draw.Draw(thumbnailImage, bounds, orig, image.Pt(0, 0), draw.Over)

	var buff bytes.Buffer
	if err := png.Encode(&buff, thumbnailImage); err != nil {
		panic("png.Encode() should not have encountered I/O failure")
	}

	data = buff.Bytes()
	return
}
