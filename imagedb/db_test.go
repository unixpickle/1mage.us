package imagedb

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewDb(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_new_db")
	if err != nil {
		t.Fatal("failed to create workspace directory")
	}
	defer os.RemoveAll(tempDir)

	_, err = NewDb(filepath.Join(tempDir, "db1"))
	if err != nil {
		t.Error(err)
	} else {
		_, err := NewDb(filepath.Join(tempDir, "db1"))
		if err != nil {
			t.Error(err)
		}
	}

	infoPath := filepath.Join(tempDir, "db1", "info.json")
	if err := os.Remove(infoPath); err != nil {
		t.Error("failed to remove info.json:", err)
	} else {
		_, err := NewDb(filepath.Join(tempDir, "db1"))
		if err == nil {
			t.Error("did not get error for missing info.json")
		}
	}

	if err := ioutil.WriteFile(infoPath, []byte("yo yo!"), 0700); err != nil {
		t.Error("failed to write bogus info.json:", err)
	} else {
		_, err := NewDb(filepath.Join(tempDir, "db1"))
		if err == nil {
			t.Error("bogus info.json should have broken everything.")
		}
	}
}

func TestAdd(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_db_add")
	if err != nil {
		t.Fatal("failed to create workspace directory")
	}
	defer os.RemoveAll(tempDir)

	db, err := NewDb(filepath.Join(tempDir, "db"))

	// Test regular additions where everything works normally.
	imageDatas := [][]byte{[]byte("this isn't a real image")}
	_, pngData := makeTestImage(1920, 1080)
	imageDatas = append(imageDatas, pngData)
	for i, imageData := range imageDatas {
		imageFile, err := ioutil.TempFile(tempDir, "test_db_add")
		if err != nil {
			t.Fatal("failed to create temp file:", err)
		}
		if _, err := imageFile.Write(imageData); err != nil {
			imageFile.Close()
			t.Fatal("failed to write image:", err)
		}

		entry, version, err := db.Add(imageFile, "image/png")
		if err != nil {
			t.Error("failed to add image to DB:", err)
		} else {
			if version != int64(i+1) {
				t.Error("invalid version:", version)
			}
			if entry.Id != int64(i) {
				t.Error("invalid ID:", entry.Id)
			}
		}
		if imageFile.Close() == nil {
			t.Error("the imageFile should already have been closed")
		}
		if _, err := os.Stat(imageFile.Name()); err == nil || !os.IsNotExist(err) {
			t.Error("the imageFile still exists")
		}
	}

	if images, _ := db.Images(); len(images) != len(imageDatas) {
		t.Error("invalid number of images in db:", len(images))
	}

	os.RemoveAll(filepath.Join(tempDir, "db"))
	tempFile, err := ioutil.TempFile(tempDir, "test_db_add")
	if err != nil {
		t.Fatal("failed to create temp file:", err)
	}
	if _, _, err := db.Add(tempFile, "image/png"); err == nil {
		t.Error("adding another image should have failed")
	}
	if tempFile.Close() == nil {
		t.Error("the temporary file should already have been closed")
	}
	if _, err := os.Stat(tempFile.Name()); err == nil || !os.IsNotExist(err) {
		t.Error("the temporary file still exists")
	}
}

