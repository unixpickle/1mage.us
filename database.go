package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/howeyc/gopass"
)

// An Image contains metadata about a single post. It does not store the image data itself.
type Image struct {
	MIME      string `json:"mimeType"`
	Seq       int    `json:"id"`
	Timestamp int64  `json:"timestamp"`

	HasSize bool `json:"has_size"`
	Width   int  `json:"width"`
	Height  int  `json:"height"`

	HasThumbnail    bool `json:"has_thumbnail"`
	ThumbnailWidth  int  `json:"thumb_width"`
	ThumbnailHeight int  `json:"thumb_height"`
}

// The Database stores a history of images posted to this site.
//
// Database extends sync.RWMutex and should be locked according to the requirements of the user. If
// you plan to modify Images or CurrentId, you should lock the database for writing. However, if you
// are simply accessing the fields for reading, you only need to lock it for reading. Certain
// methods will require you to lock the Database for reading or writing. Consult the documentation
// of these methods for more information.
type Database struct {
	sync.RWMutex `json:"-"`

	DbPath        string `json:"-"`
	DirectoryPath string `json:"-"`

	CurrentId    int     `json:"current_id"`
	Images       []Image `json:"images"`
	PasswordHash string  `json:"password_hash"`
}

// LoadDatabase loads a database from a given data directory. If the database was not yet
// configured, this will prompt the user to setup a new database.
func LoadDatabase(path string) (*Database, error) {
	dbPath := filepath.Join(path, "db.json")
	database := Database{DirectoryPath: path, DbPath: dbPath, Images: []Image{}}

	if contents, err := ioutil.ReadFile(dbPath); err == nil {
		if err := json.Unmarshal(contents, &database); err != nil {
			return nil, err
		}
	} else if err = database.Save(); err != nil {
		return nil, err
	}

	if database.PasswordHash == "" {
		fmt.Print("Setup new password: ")
		pass := gopass.GetPasswdMasked()
		database.PasswordHash = HashPassword(string(pass))
		if err := database.Save(); err != nil {
			return nil, err
		}
	}

	return &database, nil
}

// Reload updates the configuration parameters by reading them from the DB file.
// This allows the administrator to update the server's configuration while it is running.
// This will not re-load the images and current ID from the database, since administrators
// should have no reason to modify them by hand.
// The database must be locked for writing.
func (d *Database) ReloadConfig() error {
	var newDb Database
	if contents, err := ioutil.ReadFile(d.DbPath); err != nil {
		return err
	} else if err = json.Unmarshal(contents, &newDb); err != nil {
		return err
	}
	d.PasswordHash = newDb.PasswordHash
	return nil
}

// Save writes the database to the file from which it was loaded.
// This requires that the database was locked for writing.
func (d *Database) Save() error {
	if data, err := json.Marshal(d); err != nil {
		return err
	} else {
		return ioutil.WriteFile(d.DbPath, data, 0700)
	}
}

// HashPassword returns the SHA-256 hash of a string.
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}
