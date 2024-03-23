package main

import (
	"flag"
	"log"
	"time"

	"github.com/dustin/go-humanize"
)

func cmdPurgeLocal(args []string) {
	log.Fatalln("purge-local is not implemented yet.")
}

func cmdPurgeRemote(args []string) {
	flag := flag.NewFlagSet("dendrite-media purge-remote", flag.ExitOnError)
	days := flag.Uint64("days", 30, "Purge remote media older than n days")
	flag.Parse(args)

	var count, storage uint64
	// var offset int
	var dryRunTag string

	if *dryRun {
		dryRunTag = "(pretend) "
	}

	for {

		files, err := media.GetRemoteMedia(0, 50)
		if err != nil {
			log.Fatalf("failed to query for remote media: %v\n", err)
		}

		for i := 0; i < len(files); i++ {

			// is the file old enough to be pruned?
			if !olderThan(files[i].Created, *days) {
				// log.Printf("file %s is not old enough to remove, breaking out", files[i].ID)
				continue
			}

			storage += files[i].Size
			count += 1

			log.Printf("%sdeleting %s (%s)...\n", dryRunTag, files[i].ID, humanize.Bytes(files[i].Size))

			if *dryRun {
				// don't actually do anything
				continue
			}

			err := files[i].Delete()
			if err != nil {
				log.Fatalf("failed to delete %s: %v\n", files[i].ID, err)
			}
		}

		if len(files) < 50 {
			break
		}

		// offset += 50

	}

	log.Printf("> Pruned %d files totalling %s. %s\n", count, humanize.Bytes(storage), dryRunTag)
}

func cmdPurgeOrigin(args []string) {

	flag := flag.NewFlagSet("dendrite-media purge-origin", flag.ExitOnError)
	origin := flag.String("origin", "", "Origin to purge all media from")
	flag.Parse(args)

	if *origin == "" {
		log.Fatalf("you must specify a origin to purge media from.")
	}

	var count, storage uint64
	var dryRunTag string

	if *dryRun {
		dryRunTag = "(pretend) "
	}

	for {

		// TODO: Start from the newest and work backwards?
		//		 It feels like doing back->front could cause problems.

		files, err := media.GetMediaFromOrigin(*origin, 0, 50)
		if err != nil {
			log.Fatalf("failed to query media from origin: %v\n", err)
		}

		for i := 0; i < len(files); i++ {
			storage += files[i].Size
			count += 1

			log.Printf("%sdeleting %s (%s)...\n", dryRunTag, files[i].ID, humanize.Bytes(files[i].Size))

			if *dryRun {
				// don't actually do anything
				continue
			}

			err := files[i].Delete()
			if err != nil {
				log.Fatalf("failed to delete %s: %v\n", files[i].ID, err)
			}
		}

		if len(files) < 50 {
			// no more files to process.
			break
		}

	}

	log.Printf("%sdeleted %d files totalling %s\n", dryRunTag, count, humanize.Bytes(storage))

}

func cmdPurgeUser(args []string) {

	flag := flag.NewFlagSet("dendrite-media purge-user", flag.ExitOnError)
	user := flag.String("user", "", "(Local) user to purge all media from")
	flag.Parse(args)

	if *user == "" {
		log.Fatalf("you must specify a user to purge media from.")
	}

	var count, storage uint64
	var dryRunTag string

	if *dryRun {
		dryRunTag = "(pretend) "
	}

	for {

		files, err := media.GetMediaByUser(*user, 0, 50)
		if err != nil {
			log.Fatalf("failed to query media from origin: %v\n", err)
		}

		for i := 0; i < len(files); i++ {
			storage += files[i].Size
			count += 1

			log.Printf("%sdeleting %s (%s)...\n", dryRunTag, files[i].ID, humanize.Bytes(files[i].Size))

			if *dryRun {
				// don't actually do anything
				continue
			}

			err := files[i].Delete()
			if err != nil {
				log.Fatalf("failed to delete %s: %v\n", files[i].ID, err)
			}
		}

		if len(files) < 50 {
			// no more files to process.
			break
		}
	}

	log.Printf("%sdeleted %d files totalling %s\n", dryRunTag, count, humanize.Bytes(storage))

}

func cmdPurgeMxid(args []string) {

	flag := flag.NewFlagSet("dendrite-media purge-mxid", flag.ExitOnError)
	mxid := flag.String("mxid", "", "mxid of file to purge")
	flag.Parse(args)

	if *mxid == "" {
		log.Fatalf("you must specify an mxid to purge.")
	}

	file, err := media.GetMediaByID(*mxid)
	if err != nil {
		log.Fatalf("failed to query media: %v\n", err)
	}

	if *dryRun {
		log.Printf("> Pretended to remove media %s (%s)\n", file.ID, humanize.Bytes(file.Size))
		return
	}

	err = file.Delete()
	if err != nil {
		log.Fatalf("failed to delete %s: %v\n", file.ID, err)
	}

	log.Printf("> Removed media %s (%s)\n", file.ID, humanize.Bytes(file.Size))

}

func cmdMediaInfo(args []string) {

	originStorage := make(map[string]uint64)
	originCount := make(map[string]uint64)

	var totalStorage, totalCount uint64

	offset := 0

	for {

		files, err := media.GetRemoteMedia(offset, 50)
		if err != nil {
			log.Fatalf("failed to get media info: %v\n", err)
		}

		for i := 0; i < len(files); i++ {
			originStorage[files[i].Origin] += files[i].Size
			originCount[files[i].Origin] += 1

			totalStorage += files[i].Size
			totalCount += 1
		}

		if len(files) < 50 {
			break
		}

		offset += 50
	}

	// reset offset
	offset = 0
	var localStorage, localCount uint64

	for {

		files, err := media.GetLocalMedia(offset, 50)
		if err != nil {
			log.Fatalf("failed to get media info: %v\n", err)
		}

		for i := 0; i < len(files); i++ {
			localStorage += files[i].Size
			localCount += 1

			totalStorage += files[i].Size
			totalCount += 1
		}

		if len(files) < 50 {
			break
		}

		offset += 50
	}

	log.Printf("> Your instance has used %s with %d file(s).", humanize.Bytes(localStorage), localCount)

	for origin, usage := range originStorage {
		log.Printf("> Origin %s has used %s with %d file(s).\n", origin, humanize.Bytes(usage), originCount[origin])
	}

	log.Printf("> Total usage is %s with %d file(s).\n", humanize.Bytes(totalStorage), totalCount)

}

// kinda yucky feeling, but w/e
func olderThan(ts int64, days uint64) bool {

	return uint64(time.Now().Sub(time.UnixMilli(ts))) > (uint64(time.Hour)*24)*days
}
