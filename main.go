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
	manView   = man.Flag("view", "view your existing man page").Short('v').Bool()

	ls      = app.Command("ls", "List all the pages you currently have")
	lsMan   = ls.Flag("man", "see all man pages").Short('m').Bool()
	lsNotes = ls.Flag("notes", "see all notes").Short('n').Bool()

	notes     = app.Command("notes", "generate notes")
	notesName = notes.Arg("name", "name of the notes page").Required().String()

	notesGen    = notes.Flag("generate", "generate a new notes page").Short('g').Bool()
	notesDelete = notes.Flag("delete", "delete the given notes page").Short('d').Bool()
	notesEdit   = notes.Flag("edit", "edit a preexisting notes page").Short('e').Bool()

	deploy = app.Command("deploy", "Deploy your notes onto netlify")

	// link = app.Command("link", "Sync your local notes before deploying")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case man.FullCommand():
		handleMan()
	case ls.FullCommand():
		handleLs()
	case notes.FullCommand():
		handleNotes()
	// case link.FullCommand():
	// 	commands.Link()
	case deploy.FullCommand():
		commands.Deploy()
	default:
		log.Println("Subcommand not recognized.")
	}
}

func handleLs() {
	if *lsMan {
		commands.Ls("man")
	} else if *lsNotes {
		commands.Ls("notes")
	} else {
		log.Println("Flag not recongized")
	}
}

func handleMan() {
	if *manDelete {
		if err := commands.DeletePage(*manName, "man"); err != nil {
			log.Println("Could not delete page")
		}
	} else if *manGen {
		if err := commands.GenPage(*manName, "man"); err != nil {
			log.Println("Could not create page")
		}
	} else if *manEdit {
		if err := commands.EditPage(*manName, "man"); err != nil {
			log.Println("Could not edit page")
		}
	} else if *manView {
		commands.ViewMan(*manName)
	} else {
		log.Println("Flag not recognized")
	}
}

func handleNotes() {
	if *notesDelete {
		if err := commands.DeletePage(*notesName, "notes"); err != nil {
			log.Println("Could not delete page")
		}
	} else if *notesGen {
		if err := commands.GenPage(*notesName, "notes"); err != nil {
			log.Println("Could not create page")
		}
	} else if *notesEdit {
		if err := commands.EditPage(*notesName, "notes"); err != nil {
			log.Println("Could not edit page")
		}
	} else {
		log.Println("Flag not recognized")
	}
}
