package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type MediaStore struct {
	DB        *sql.DB
	MediaPath string
}

func (ms *MediaStore) Connect(url string) error {
	var err error

	ms.DB, err = sql.Open("postgres", url)

	return err

}

func (ms *MediaStore) GetMediaByUser(mxid string, offset, limit int) ([]*File, error) {

	var files []*File

	rows, err := ms.DB.Query("SELECT media_id, file_size_bytes, creation_ts, base64hash, user_id, media_origin FROM mediaapi_media_repository WHERE user_id = $1 ORDER BY creation_ts DESC OFFSET $2 LIMIT $3", mxid, offset, limit)
	if err != nil {
		log.Printf("failed to query for media from user %s: %v\n", mxid, err)
		return []*File{}, nil
	}

	for rows.Next() {
		f := File{
			Store: ms,
		}

		err = rows.Scan(&f.ID, &f.Size, &f.Created, &f.Hash, &f.UserID, &f.Origin)
		if err != nil {
			log.Printf("failed to scan row: %v\n", err)
			return []*File{}, err
		}

		files = append(files, &f)

	}

	return files, nil
}

func (ms *MediaStore) GetLocalMedia(offset, limit int) ([]*File, error) {

	var files []*File

	rows, err := ms.DB.Query("SELECT media_id, file_size_bytes, creation_ts, base64hash, user_id, media_origin FROM mediaapi_media_repository WHERE user_id != '' ORDER BY creation_ts DESC OFFSET $1 LIMIT $2", offset, limit)
	if err != nil {
		log.Printf("failed to query for local media: %v\n", err)
		return []*File{}, nil
	}

	for rows.Next() {
		f := File{
			Store: ms,
		}

		err = rows.Scan(&f.ID, &f.Size, &f.Created, &f.Hash, &f.UserID, &f.Origin)
		if err != nil {
			log.Printf("failed to scan row: %v\n", err)
			return []*File{}, err
		}

		files = append(files, &f)

	}

	return files, nil
}

func (ms *MediaStore) GetRemoteMedia(offset, limit int) ([]*File, error) {

	var files []*File

	rows, err := ms.DB.Query("SELECT media_id, file_size_bytes, creation_ts, base64hash, user_id, media_origin FROM mediaapi_media_repository WHERE user_id = '' ORDER BY creation_ts DESC OFFSET $1 LIMIT $2", offset, limit)
	if err != nil {
		log.Printf("failed to query for remote media: %v\n", err)
		return []*File{}, nil
	}

	for rows.Next() {
		f := File{
			Store: ms,
		}

		err = rows.Scan(&f.ID, &f.Size, &f.Created, &f.Hash, &f.UserID, &f.Origin)
		if err != nil {
			log.Printf("failed to scan row: %v\n", err)
			return []*File{}, err
		}

		files = append(files, &f)

	}

	return files, nil
}

func (ms *MediaStore) GetMediaFromOrigin(origin string, offset, limit int) ([]*File, error) {

	var files []*File

	rows, err := ms.DB.Query("SELECT media_id, file_size_bytes, creation_ts, base64hash, user_id, media_origin FROM mediaapi_media_repository WHERE media_origin = $1 ORDER BY creation_ts DESC OFFSET $2 LIMIT $3", origin, offset, limit)
	if err != nil {
		log.Printf("failed to query for media from origin %s: %v\n", origin, err)
		return []*File{}, nil
	}

	for rows.Next() {
		f := File{
			Store: ms,
		}

		err = rows.Scan(&f.ID, &f.Size, &f.Created, &f.Hash, &f.UserID, &f.Origin)
		if err != nil {
			log.Printf("failed to scan row: %v\n", err)
			return []*File{}, err
		}

		files = append(files, &f)

	}

	return files, nil
}

func (ms *MediaStore) DeleteFile(id string) error {

	// TODO: my homeserver doesn't have thumbnails enabled, but I understand this
	//		 will probably not be standard practice. Support for removing thumb-
	//		 nails would make sense.

	tx, err := ms.DB.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v\n", err)
		return err
	}

	stmt, err := tx.Prepare("DELETE FROM mediaapi_media_repository WHERE media_id = $1")
	if err != nil {
		log.Fatalf("failed to prepare statement: %v\n", err)
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatalf("failed to delete entry: %v\n", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit: %v\n", err)
		return err
	}

	return nil
}
