package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/howeyc/gopass"
	"github.com/unixpickle/1mage.us/imagedb"
)

// These are the default configuration parameters.
const (
	DefaultMaxFileSize     = 5 << 20
	DefaultMaxCountPerHour = 30
)

// Config is a set of user-defined parameters for an instance of the application.
type Config struct {
	PasswordHash      string
	MaxFileSize       int64
	MaxUploadsPerHour int64
}

// Db encapsulates imagedb.Db for image storage and adds functionality for configuration storage.
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
		res.config.PasswordHash = HashPassword(string(pass))
		saveConfig = true
	}
	if res.config.MaxFileSize == 0 {
		res.config.MaxFileSize = DefaultMaxFileSize
		saveConfig = true
	}
	if res.config.MaxUploadsPerHour == 0 {
		res.config.MaxUploadsPerHour = DefaultMaxCountPerHour
		saveConfig = true
	}

	if saveConfig {
		if err := res.saveConfig(); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

// Config safely returns the current configuration.
func (d *Db) Config() Config {
	d.configLock.RLock()
	defer d.configLock.RUnlock()
	return d.config
}

// SetConfig safely updates the current configuration.
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
