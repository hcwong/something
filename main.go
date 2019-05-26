package main

import (
	"log"
	"os"

	"github.com/hcwong/nani/commands"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("nani", "Tools for me")

	gen     = app.Command("gen", "generate a new man page.")
	genName = gen.Arg("name", "name of the man page").Required().String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case gen.FullCommand():
		err := commands.GenPage(*genName)
		if err != nil {
			log.Println("Could not create page")
		}
	default:
		log.Println("Subcommand not recognized.")
	}
}
