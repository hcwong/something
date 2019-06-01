package main

import (
	"log"
	"os"

	"github.com/hcwong/nani/commands"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("nani", "Tools for me")

	man     = app.Command("man", "Handle the man pages.")
	manName = man.Arg("name", "name of the man page").Required().String()

	manGen    = man.Flag("generate", "generate a new man page").Short('g').Bool()
	manDelete = man.Flag("delete", "delete the given man page").Short('d').Bool()
	manEdit   = man.Flag("edit", "edit a preexisting man page").Short('e').Bool()

	ls = app.Command("ls", "List all the man pages you currently have")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case man.FullCommand():
		if *manDelete {
			if err := commands.DeletePage(*manName); err != nil {
				log.Println("Could not delete page")
			}
		} else if *manGen {
			if err := commands.GenPage(*manName); err != nil {
				log.Println("Could not create page")
			}
		} else if *manEdit {
			if err := commands.EditPage(*manName); err != nil {
				log.Println("Could not edit page")
			}
		} else {
			log.Println("Flag not recognized")
		}
	case ls.FullCommand():
		commands.Ls()
	default:
		log.Println("Subcommand not recognized.")
	}
}
