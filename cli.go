package main

import (
	"flag"
	"log"
	"time"

	"github.com/dustin/go-humanize"
)

func cmdPurgeLocal(args []string) {

}

func cmdPurgeRemote(args []string) {

}

func cmdPurgeOrigin(args []string) {

	flag := flag.NewFlagSet("dendrite-media purge-local", flag.ExitOnError)
	origin := flag.String("origin", "", "Origin to purge all media from")
	flag.Parse(args)

	// var offset int
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

}

func cmdPurgeMxid(args []string) {

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
