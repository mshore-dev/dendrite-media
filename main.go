package main

import (
	"flag"
	"log"
)

var (
	dbUrl     = flag.String("db", "", "postgresql database connection url")
	mediaPath = flag.String("media-path", "", "base path of the Dendrite media store")
	dryRun    = flag.Bool("dry-run", false, "specificies wether or not to do a dry run")

	media MediaStore
)

func main() {

	log.Println("dendrite-media command line tool")

	// parse global flags
	flag.Parse()

	if *dbUrl == "" || *mediaPath == "" {
		log.Fatalf("you must specifcy the base media path and database connection string.")
	}

	media = MediaStore{
		MediaPath: *mediaPath,
	}

	err := media.Connect(*dbUrl)
	if err != nil {
		log.Fatalf("failed to connect to database: %v\n", err)
	}

	args := flag.Args()

	if len(args) == 0 {
		log.Fatalln("Please specify an action.")
	}

	cmd, args := args[0], args[1:]

	switch cmd {
	case "purge-local":
		cmdPurgeLocal(args)
	case "purge-remote":
		cmdPurgeRemote(args)
	case "purge-origin":
		cmdPurgeOrigin(args)
	case "purge-user":
		cmdPurgeUser(args)
	case "purge-mxid":
		cmdPurgeMxid(args)
	case "media-info":
		cmdMediaInfo(args)
	default:
		log.Fatalf("Invalid subcommand %s\n", cmd)

	}

}
