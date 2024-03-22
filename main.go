package main

import (
	"log"
	"os"
	"time"

	"github.com/dustin/go-humanize"
)

var (
	mediaBasePath = "/var/lib/dendrite/media_store/"

	dryRun = true
)

func main() {

	media := MediaStore{}

	err := media.Connect(os.Args[1])
	if err != nil {
		panic(err)
	}

	files, err := media.GetAllMediaFromSource(os.Args[2], 0, 20)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(files); i++ {

		log.Printf("%s (%s, %s)\n", files[i].ID, humanize.Bytes(uint64(files[i].Size)), humanize.Time(time.UnixMilli(files[i].Created)))

		// err = files[i].Delete()
		// if err != nil {
		// 	panic(err)
		// }
	}

}
