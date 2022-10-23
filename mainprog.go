package main

import (
	"flag"
	"log"
	"os"
	"time"
)

// main program logic

// notify myMain to exit
var chStop chan struct{}

func init() {
	chStop = make(chan struct{}, 1)
}

func myMain(runInService bool) error {
	mode := "command line"
	if runInService {
		mode = "service"
	}
	log.Printf("myMian is running, mode is %v, pid=%d\n", mode, os.Getpid())

	// parse args
	greeting := flag.String("greeting", "xxx", "greeting string")
	userName := flag.String("name", "yyy", "user name")
	flag.Parse()

	fgStop := false
	for !fgStop {
		select {
		case <-chStop:
			fgStop = true
			log.Println("receive stop signal")
			continue
		case <-time.After(3 * time.Second):
		}

		log.Printf("myMain: %v %v\n", *greeting, *userName)
	}

	log.Println("myMain exit")
	return nil
}
