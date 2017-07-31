package main

import (
	"flag"
	"os"
)

func init() {
	//specify ecset file
	flag.StringVar(&isfile, "f", "ecset", "The settings file")

	//ask for help
	var h bool
	flag.BoolVar(&h, "h", false, "Request Usage")

	//Finally parse
	flag.Parse()

	if h {
		flag.Usage()
		os.Exit(0)
	}
}
