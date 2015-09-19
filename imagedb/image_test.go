package imagedb

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func TestProcessImage(t *testing.T) {
	sizes := [][]int{{100, 150}, {400, 450}, {300, 300}, {1000, 100}}
	thumbnailSizes := [][]int{{100, 150}, {228, 256}, {256, 256}, {256, 26}}
	for i, size := range sizes {
		tempFile, err := ioutil.TempFile("", "test_process_image")
		if err != nil {
			t.Error("failed to create temp file")
			continue
		}
		defer tempFile.Close()
		defer os.Remove(tempFile.Name())

		_, pngData := makeTestImage(size[0], size[1])
		if _, err := tempFile.Write(pngData); err != nil {
			t.Error("failed to write temp file")
			continue
		}

		entry, thumbnailFile, err := processImage(tempFile, "image/png")
		if err != nil {
			t.Error("processImage failed:", err)
			continue
		}

		if thumbnailFile == nil {
			t.Error("no thumbnail file")
		} else {
			defer thumbnailFile.Close()
			defer os.Remove(thumbnailFile.Name())

			thumbnailFile.Seek(0, 0)
			decodedThumbnail, _, err := image.Decode(thumbnailFile)
			if err != nil {
				t.Error("failed to decode thumbnail")
			} else {
				if decodedThumbnail.Bounds().Dx() != thumbnailSizes[i][0] {
					t.Error("decoded width is wrong", decodedThumbnail.Bounds().Dx(), "for", i)
				}
				if decodedThumbnail.Bounds().Dy() != thumbnailSizes[i][1] {
					t.Error("decoded height is wrong", decodedThumbnail.Bounds().Dy(), "for", i)
				}
			}
		}

		if !entry.HasThumbnail {
			t.Error("no thumbnail")
		} else {
			if entry.ThumbnailWidth != thumbnailSizes[i][0] {
				t.Error("invalid width", entry.ThumbnailWidth, "for", i)
			}
			if entry.ThumbnailHeight != thumbnailSizes[i][1] {
				t.Error("invalid height", entry.ThumbnailHeight, "for", i)
			}
		}

		if !entry.HasSize {
			t.Error("size is missing")
		} else {
			if entry.Width != size[0] {
				t.Error("invalid width", entry.Width, "for", i)
			}
			if entry.Height != size[1] {
				t.Error("invalid height", entry.Height, "for", i)
			}
		}
	}

	// Run a test for an invalid image.
	bogusFile, err := ioutil.TempFile("", "test_process_image")
	if err != nil {
		t.Error("failed to create file")
		return
	}
	defer bogusFile.Close()
	defer os.Remove(bogusFile.Name())
	bogusFile.Write([]byte("hey there, this isn't a real image"))

	entry, thumbnailFile, err := processImage(bogusFile, "image/svg+xml")
	if err != nil {
		t.Error("failed to process bogus image")
	} else {
		if thumbnailFile != nil {
			t.Error("a thumbnail file exists")
			thumbnailFile.Close()
			os.Remove(thumbnailFile.Name())
		}
		if entry.HasSize {
			t.Error("claims to have size")
		}
		if entry.HasThumbnail {
			t.Error("claims to have thumbnail")
		}
	}
}

func makeTestImage(width, height int) (image.Image, []byte) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			color := color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)),
				uint8(rand.Intn(256)), uint8(rand.Intn(256))}
			img.SetRGBA(x, y, color)
		}
	}

	var buff bytes.Buffer
	png.Encode(&buff, img)
	return img, buff.Bytes()
}
