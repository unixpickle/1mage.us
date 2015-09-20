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
	"github.com/unixpickle/1mage.us/imagedb"
)

const DefaultMaxFileSize = 5 << 20
const DefaultMaxCountPerHour = 30

var GlobalDb *Db

type Config struct {
	PasswordHash    string
	MaxFileSize     int64
	MaxCountPerHour int64
}

type Db struct {
	*imagedb.Db

	configLock sync.RWMutex
	configPath string
	config     Config
}

// SetupDb creates a new database or opens an existing one at a given directory.
// This will prompt the user to set a password if the database is being created.
func SetupDb(path string) (*Db, error) {
	imageDb, err := imagedb.NewDb(path)
	if err != nil {
		return nil, err
	}

	res := Db{imageDb, sync.RWMutex{}, filepath.Join(path, "config.json"), Config{}}
	saveConfig := false

	if res.loadConfig() != nil {
		saveConfig = true
	}
	if res.config.PasswordHash == "" {
		fmt.Print("Setup new password: ")
		pass := gopass.GetPasswdMasked()
		res.config.PasswordHash = hashPassword(string(pass))
		saveConfig = true
	}
	if res.config.MaxFileSize == 0 {
		res.config.MaxFileSize = DefaultMaxFileSize
		saveConfig = true
	}
	if res.config.MaxCountPerHour == 0 {
		res.config.MaxCountPerHour = DefaultMaxCountPerHour
		saveConfig = true
	}

	if saveConfig {
		if err := res.saveConfig(); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *Db) Config() Config {
	d.configLock.RLock()
	defer d.configLock.RUnlock()
	return d.config
}

func (d *Db) SetConfig(c Config) (err error) {
	d.configLock.Lock()
	defer d.configLock.Unlock()

	oldConfig := d.config
	d.config = c

	if err = d.saveConfig(); err != nil {
		d.config = oldConfig
	}

	return
}

func (d *Db) loadConfig() error {
	if data, err := ioutil.ReadFile(d.configPath); err != nil {
		return err
	} else {
		return json.Unmarshal(data, &d.config)
	}
}

func (d *Db) saveConfig() error {
	if data, err := json.Marshal(d.config); err != nil {
		return err
	} else {
		return ioutil.WriteFile(d.configPath, data, 0700)
	}
}

// hashPassword returns the SHA-256 hash of a string.
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}
