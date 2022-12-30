package main

import (
	"flag"
	"log"

	"github.com/Doridian/fakerfs/dev"
	"github.com/Doridian/fakerfs/ffs"
	"github.com/Doridian/fakerfs/util"
)

func main() {
	srcStr := flag.String("src", "./src", "Mount source")
	targetStr := flag.String("target", "./target", "Mount target")
	configStr := flag.String("config", "./config.yml", "Config file")
	flag.Parse()

	cfg, err := util.LoadConfig(*configStr)
	if err != nil {
		panic(err)
	}

	fakeFS, err := ffs.NewFakerFS(*srcStr)
	if err != nil {
		panic(err)
	}

	flatFiles := cfg.Flatten("")

	for path, file := range flatFiles {
		var handler dev.FileHandler
		switch file.Type {
		case "fixed":
			handler = &dev.FixedHandler{}
		case "integer":
			handler = &dev.IntegerHandler{}
		case "choice":
			handler = &dev.ChoiceHandler{}
		}
		err = handler.LoadConfig(file.Config)
		if err != nil {
			panic(err)
		}
		fuseFile := dev.MakeFile(path, handler)
		fakeFS.AddHandler(fuseFile)
	}

	err = fakeFS.Mount(*targetStr)
	if err != nil {
		panic(err)
	}
	log.Printf("Mounted %v to %v with %d fakes", *srcStr, *targetStr, len(flatFiles))

	fakeFS.Wait()
}
