package main

import (
	"database/sql"
	"log"
	"os"
	"path"

	"github.com/dustin/go-humanize"
	_ "github.com/lib/pq"
)

var (
	mediaBasePath = "/var/lib/dendrite/media_store/"

	dryRun = true

	db *sql.DB
)

type File struct {
	ID      string
	Size    int64
	Created int64
	Hash    string
	UserID  string
}

func (f *File) Path() string {

	if f.Hash == "" {
		return ""
	}

	return path.Join(mediaBasePath, f.Hash[0:1], f.Hash[1:2], f.Hash[2:])
}

func (f *File) Delete() error {

	if dryRun {
		log.Printf("[dry] deleting file %s (%s)...\n", f.Path(), humanize.Bytes(uint64(f.Size)))
		return nil
	}

	fullPath := f.Path()

	log.Printf("[info] deleting file %s (%s)...\n", f.Path(), humanize.Bytes(uint64(f.Size)))

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v\n", err)
		return err
	}

	stmt, err := tx.Prepare("DELETE FROM mediaapi_media_respository WHERE media_id = ?;")
	if err != nil {
		log.Fatalf("failed to prepare statement: %v\n", err)
		return err
	}

	_, err = stmt.Exec(f.ID)
	if err != nil {
		log.Fatalf("failed to delete entry: %v\n", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit: %v\n", err)
		return err
	}

	err = os.RemoveAll(fullPath)
	if err != nil {
		log.Fatalf("failed to remove file path: %v\n", err)
		return err
	}

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

func main() {
	pgUrl := os.Args[1]

	var err error

	db, err = sql.Open("postgres", pgUrl)
	if err != nil {
		panic(err)
	}

	var files []File

	rows, err := db.Query("SELECT media_id, file_size_bytes, creation_ts, base64hash, user_id FROM mediaapi_media_repository ORDER BY creation_ts ASC LIMIT 25;")
	if err != nil {
		log.Fatalf("failed to query (at most) 25 medias: %v\n", err)
	}

	for rows.Next() {

		var f File

		err = rows.Scan(&f.ID, &f.Size, &f.Created, &f.Hash, &f.UserID)
		if err != nil {
			log.Fatalf("failed to scan row: %v\n", err)
		}

		files = append(files, f)

	}

	for i := 0; i < len(files); i++ {

		exists, err := files[i].Exists()
		if err != nil {
			panic(err)
		}

		log.Printf("file exists? %d\n", exists)

		err = files[i].Delete()
		if err != nil {
			log.Fatalf("failed to delete %s: %v\n", files[i].ID, err)
		}
	}
}
