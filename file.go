package main

import (
	"log"
	"os"
	"path"

	"github.com/dustin/go-humanize"
)

type File struct {
	ID      string
	Size    uint64
	Created int64
	Hash    string
	UserID  string
	Origin  string

	Store *MediaStore
}

func (f *File) Path() string {

	if f.Hash == "" {
		return ""
	}

	return path.Join(f.Store.MediaPath, f.Hash[0:1], f.Hash[1:2], f.Hash[2:])
}

func (f *File) Delete() error {

	fullPath := f.Path()

	log.Printf("[info] deleting file %s (%s)...\n", f.Path(), humanize.Bytes(uint64(f.Size)))

	f.Store.DeleteFile(f.ID)

	err := os.RemoveAll(fullPath)
	if err != nil {
		log.Fatalf("failed to remove file path: %v\n", err)
		return err
	}

	// exists, err := f.Exists()
	// if err != nil {
	// 	log.Fatalf("failed to confirm deletion: %v\n", err)
	// 	return err
	// }

	// if exists {
	// 	log.Printf("file %s was supposed to be gone. why is it here?\n", f.Path())
	// 	return fmt.Errorf("file wasn't gone'd")
	// }

	return nil
}

func (f *File) Exists() (bool, error) {

	_, err := os.Stat(path.Join(f.Path(), "file"))
	if err != nil {
		log.Printf("could not stat %s: %v\n", f.Path(), err)
		return false, err
	}

	return true, nil
}

func (f *File) HasThumbnail() (bool, error) {

	return false, nil
}
