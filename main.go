package main

import (
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("nani", "Tools for me")
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
