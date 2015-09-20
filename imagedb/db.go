package imagedb

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var NoSuchImageErr = errors.New("the specified image does not exist")
var NotADirectoryErr = errors.New("not a directory")
var infoFileName = "info.json"

// The Db stores a history of images and manages their associated files.
type Db struct {
	lock sync.RWMutex

	infoPath string
	dirPath  string

	info struct {
		CurrentId int64   `json:"current_id"`
		Version   int64   `json:"version"`
		Images    []Image `json:"images"`
	}
}

// NewDb creates a new database or opens an existing one at a given directory.
//
// Only one process should open a database at any given time.
// The database is safe for concurrency in one Go process, but not across processes.
func NewDb(dirPath string) (*Db, error) {
	if info, err := os.Stat(dirPath); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return createDb(dirPath)
	} else if !info.IsDir() {
		return nil, NotADirectoryErr
	} else {
		return loadDb(dirPath)
	}
}

func createDb(dirPath string) (*Db, error) {
	var db Db
	db.dirPath = dirPath
	db.infoPath = filepath.Join(db.dirPath, infoFileName)

	if err := os.Mkdir(db.dirPath, 0700); err != nil {
		return nil, err
	} else if err := db.writeToFile(); err != nil {
		return nil, err
	} else {
		return &db, nil
	}
}

func loadDb(dirPath string) (*Db, error) {
	var db Db
	db.dirPath = dirPath
	db.infoPath = filepath.Join(db.dirPath, infoFileName)

	if data, err := ioutil.ReadFile(db.infoPath); err != nil {
		return nil, err
	} else if err := json.Unmarshal(data, &db.info); err != nil {
		return nil, err
	} else {
		return &db, nil
	}
}

// Add adds an image to the database.
//
// You must supply a temporary file with the image data and the MIME type of the file.
// This will automatically generate a thumbnail for the image and detect its dimensions.
//
// Whether or not this succeeds, the temporary file you provide will automatically be closed and
// deleted.
func (d *Db) Add(tempFile *os.File, mimeType string) (entry Image, version int64, err error) {
	entry, thumbnailFile, err := processImage(tempFile, mimeType)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return
	}

	tempFile.Close()
	if thumbnailFile != nil {
		thumbnailFile.Close()
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	entry.Id = d.info.CurrentId
	d.info.CurrentId++

	d.info.Version++
	version = d.info.Version

	if err = os.Rename(tempFile.Name(), d.pathForImage(entry.Id)); err != nil {
		d.info.Version--
		d.info.CurrentId--
		os.Remove(tempFile.Name())
		if thumbnailFile != nil {
			os.Remove(thumbnailFile.Name())
		}
		return
	}

	if thumbnailFile != nil {
		if err = os.Rename(thumbnailFile.Name(), d.pathForThumbnail(entry.Id)); err != nil {
			d.info.Version--
			d.info.CurrentId--
			os.Remove(d.pathForImage(entry.Id))
			os.Remove(thumbnailFile.Name())
			return
		}
	}

	oldImages := d.info.Images
	newImages := make([]Image, len(oldImages)+1)
	copy(newImages, oldImages)
	newImages[len(oldImages)] = entry
	d.info.Images = newImages

	if err = d.writeToFile(); err != nil {
		d.info.Version--
		d.info.CurrentId--
		d.info.Images = oldImages
		os.Remove(d.pathForImage(entry.Id))
		os.Remove(d.pathForThumbnail(entry.Id))
	}

	return
}

// Delete removes an image from the database.
func (d *Db) Delete(imageId int64) (version int64, err error) {
	d.lock.Lock()

	imageIndex := -1
	var deletingImage Image
	for i, image := range d.info.Images {
		if image.Id == imageId {
			deletingImage = image
			imageIndex = i
		}
	}
	if imageIndex < 0 {
		d.lock.Unlock()
		return 0, NoSuchImageErr
	}

	d.info.Version++
	version = d.info.Version

	oldImages := d.info.Images
	newImages := make([]Image, len(oldImages)-1)
	copy(newImages, oldImages[:imageIndex])
	copy(newImages[imageIndex:], oldImages[imageIndex+1:])
	d.info.Images = newImages

	if err = d.writeToFile(); err != nil {
		// If we fail to save the database, we can undo our changes without harm.
		d.info.Version--
		d.info.Images = oldImages
		d.lock.Unlock()
		return
	}

	d.lock.Unlock()

	os.Remove(d.pathForImage(imageId))
	if deletingImage.HasThumbnail {
		os.Remove(d.pathForThumbnail(imageId))
	}

	return
}

// Images returns a read-only list of images in the database at the current moment.
//
// This also returns the current version of the database.
// This is useful for detecting when a list of images is out of date.
// See Version() for more information.
func (d *Db) Images() (images []Image, version int64) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.info.Images, d.info.Version
}

// OpenImage creates an io.ReadCloser for a given image.
func (d *Db) OpenImage(imageId int64) (io.ReadCloser, error) {
	return os.Open(d.pathForImage(imageId))
}

// OpenThumbnail creates an io.ReadCloser for an image's thumbnail.
func (d *Db) OpenThumbnail(imageId int64) (io.ReadCloser, error) {
	return os.Open(d.pathForThumbnail(imageId))
}

// Version returns the number of changes which have been made to this database since its creation.
// This number is useful for tracking if client-side information is out of date.
func (d *Db) Version() int64 {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.info.Version
}

func (d *Db) pathForImage(imageId int64) string {
	return filepath.Join(d.dirPath, strconv.FormatInt(imageId, 10))
}

func (d *Db) pathForThumbnail(imageId int64) string {
	return d.pathForImage(imageId) + "_thumb"
}

func (d *Db) writeToFile() error {
	if data, err := json.Marshal(&d.info); err != nil {
		return err
	} else {
		return ioutil.WriteFile(d.infoPath, data, 0700)
	}
}
