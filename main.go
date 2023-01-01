package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Doridian/fakerfs/ffs"
	"github.com/Doridian/fakerfs/sysfs"
	"github.com/Doridian/fakerfs/util"
	"github.com/sevlyar/go-daemon"
)

var srcStr = flag.String("src", "./src", "Mount source")
var targetStr = flag.String("target", "./target", "Mount target")
var configStr = flag.String("config", "./config.yml", "Config file")

var daemonizeBool = flag.Bool("daemonize", false, "Daemonize")
var logFileStr = flag.String("log-file", "", "Log file")
var pidFileStr = flag.String("pid-file", "", "PID file")

var fakeFS *ffs.FakerFS

func runFSMount() {
	cfg, err := util.LoadConfig(*configStr)
	if err != nil {
		panic(err)
	}

	fakeFS, err = ffs.NewFakerFS(*srcStr)
	if err != nil {
		panic(err)
	}

	flatFiles := cfg.Flatten("")

	for path, file := range flatFiles {
		var handler sysfs.FileHandler
		switch file.Type {
		case "fixed":
			handler = &sysfs.FixedHandler{}
		case "integer":
			handler = &sysfs.IntegerHandler{}
		case "choice":
			handler = &sysfs.ChoiceHandler{}
		}
		err = handler.LoadConfig(file.Config)
		if err != nil {
			panic(err)
		}
		fuseFile := sysfs.MakeFile(handler)
		fakeFS.AddHandler(path, fuseFile)
	}

	err = fakeFS.Mount(*targetStr, false)
	if err != nil {
		panic(err)
	}
	log.Printf("Mounted %v to %v with %d fakes", *srcStr, *targetStr, len(flatFiles))
}

func waitFSMount() {
	fakeFS.Wait()
}

func main() {
	flag.Parse()

	if !*daemonizeBool {
		runFSMount()
		waitFSMount()
		return
	}

	cntxt := &daemon.Context{
		PidFileName: *pidFileStr,
		PidFilePerm: 0644,
		LogFileName: *logFileStr,
		LogFilePerm: 0640,
	}

	child, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to daemonize: ", err)
	}

	if child != nil {
		errorChannel := make(chan error, 1)
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGUSR1)

		go func() {
			_, childErr := child.Wait()
			if childErr == nil {
				childErr = errors.New("child died unexpectedly")
			}
			errorChannel <- childErr
		}()

		select {
		case err := <-errorChannel:
			log.Fatalf("Child error: %v", err)
		case <-signals:
			log.Printf("Mounted %v to %v in the background", *srcStr, *targetStr)
		}
		return
	}
	defer cntxt.Release()

	runFSMount()
	syscall.Kill(os.Getppid(), syscall.SIGUSR1)
	waitFSMount()
}
