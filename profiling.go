package main

import (
	"log"
	"os"

	"github.com/felixge/fgprof"
)

// profile returns a function that can be deferred to write a performance profile
// to the given file.
func profile(file string) func() error {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	return fgprof.Start(f, fgprof.FormatPprof)
}
