package main

import (
	"log"
	"os"

	"github.com/kardianos/service"
)

// common entry to run this program as a command line or a service
//   start from command line:
//     if argv[1]==install/uninstall/remove/start/restart/stop: then
//       manage service
//     else
//       run myMain
//   start from SCM(service control manager):
//     run myMain
func main() {

	// start program by service manager
	if !service.Interactive() {
		log.Printf("start program by service manager: %v, pid=%v\n", os.Args, os.Getpid())
		if err := runMyMainAsService(); err != nil {
			log.Fatal("service:", err)
		}
	}

	// run service managing command in command line
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		if svcCmd, ok := parseSvcCmd(cmd); ok {
			log.Printf("run SCM command: %v, pid=%d\n", svcCmd, os.Getpid())
			if err := manageSvc(svcCmd); err != nil {
				log.Fatal("SCM command:", err)
			}
			return
		}
	}

	// start program by command line
	log.Printf("start program by command line: %v, pid=%d\n", os.Args, os.Getpid())
	if err := myMain(false); err != nil {
		log.Fatal("command line:", err)
	}
	return
}
