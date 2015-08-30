package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/unixpickle/1mage.us/thumbnail"
	"labix.org/v2/mgo"
)

type Image struct {
	MIME      string `json:"mimeType" bson:"mime"`
	Seq       int    `json:"id" bson:"sequence"`
	Timestamp int64  `json:"timestamp"`

	HasSize bool `json:"has_size"`
	Width   int  `json:"width"`
	Height  int  `json:"height"`

	HasThumbnail    bool `json:"has_thumbnail"`
	ThumbnailWidth  int  `json:"thumb_width"`
	ThumbnailHeight int  `json:"thumb_height"`
}

type ImageDb struct {
	Images    []Image `json:"images"`
	CurrentId int     `json:"current_id"`
}

func die(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 3 {
		die("Usage: migrate_db <images path> <output path>")
	}

	imagesPath := os.Args[1]
	outputPath := os.Args[2]

	if info, err := os.Stat(outputPath); err != nil {
		die(err)
	} else if !info.IsDir() {
		die("Given output path is not a directory.")
	}

	session, err := mgo.Dial("mongodb://127.0.0.1:27017/1mage")
	if err != nil {
		die(err)
	}
	collection := session.DB("1mage").C("images")
	var result []Image
	if err := collection.Find(nil).All(&result); err != nil {
		die(err)
	}

	var maxId int
	for i := range result {
		imageInfo := &result[i]

		fmt.Println("Processing", imageInfo.Seq)
		if imageInfo.Seq > maxId {
			maxId = imageInfo.Seq
		}

		imagePath := filepath.Join(imagesPath, strconv.Itoa(imageInfo.Seq))

		f, err := os.Open(imagePath)
		if err != nil {
			die(err)
		}
		defer f.Close()

		if stats, err := f.Stat(); err != nil {
			die(err)
		} else {
			imageInfo.Timestamp = stats.ModTime().Unix()
		}

		newImagePath := filepath.Join(outputPath, strconv.Itoa(imageInfo.Seq))
		thumbnailPath := filepath.Join(outputPath, strconv.Itoa(imageInfo.Seq)+
			"_thumb")

		if err := copyFile(imagePath, newImagePath); err != nil {
			die(err)
		}

		img, _, err := image.Decode(f)
		if err != nil {
			fmt.Println("Notice: could not decode image: " + imagePath)
			continue
		}
		imageInfo.HasSize = true
		imageInfo.Width = img.Bounds().Dx()
		imageInfo.Height = img.Bounds().Dy()

		thumb, width, height := thumbnail.MakeThumbnail(img)
		imageInfo.HasThumbnail = true
		imageInfo.Width = width
		imageInfo.Height = height
		if err := ioutil.WriteFile(thumbnailPath, thumb, 0700); err != nil {
			die(err)
		}
	}

	var db ImageDb
	db.Images = result
	db.CurrentId = maxId + 1
	dbPath := filepath.Join(outputPath, "db.json")
	if dbData, err := json.Marshal(db); err != nil {
		die(err)
	} else if err = ioutil.WriteFile(dbPath, dbData, 0700); err != nil {
		die(err)
	}
}

func copyFile(source, dest string) error {
	contents, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dest, contents, 0700)
}
