package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/kardianos/service"
)

//
// run myMain as a service
//

// wrap myMain as a service
type program struct {
}

// start service by SCM(service control manager)
func (p *program) Start(s service.Service) error {
	log.Printf("service start is called, pid=%d\n", os.Getpid())
	// must start myMain in backgroud, MUST NOT block here
	go func() {
		if err := myMain(!service.Interactive()); err != nil {
			log.Println("failed running as a service")
		} else {
			log.Println("exit running as a service")
		}
	}()
	return nil
}

// stop service by SCM
func (p *program) Stop(s service.Service) error {
	log.Printf("service stop is called, pid=%d\n", os.Getpid())
	select {
	case chStop <- struct{}{}:
		time.Sleep(1 * time.Second)
	default:
	}
	return nil
}

// new a service with my program
func newMyMainService() (service.Service, error) {
	svcConfig := &service.Config{
		Name:        "GoService",
		DisplayName: "GoServiceDisp",
		Description: "windows service from golang",
		Option: service.KeyValue{
			"StartType":              "automatic",
			"OnFailure":              "restart",
			"OnFailureDelayDuration": "10s",
			"OnFailureResetPeriod":   "30s",
		},
		Arguments: []string{
			"-greeting=hi",
			"-name=shang",
		},
	}

	return service.New(&program{}, svcConfig)
}

// start service event loop
func runMyMainAsService() error {
	// switch logout to file when running service
	f, err := os.Create("e:\\gowinservice.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	log.Printf("start program by service manager: %v, pid=%d\n", os.Args, os.Getpid())

	if svc, err := newMyMainService(); err != nil {
		log.Println("new service:", err)
		return err
	} else if err := svc.Run(); err != nil {
		log.Print("run service:", err)
		return err
	}
	return nil
}

//
// manage service:
//   install/uninstall/remove
//   start/stop/restart
//

type SvcCmdType string

func (scmd SvcCmdType) String() string {
	return string(scmd)
}

const (
	SvcCmd_Install   SvcCmdType = "install"
	SvcCmd_UnInstall SvcCmdType = "uninstall"
	SvcCmd_Remove    SvcCmdType = "remove"
	SvcCmd_Start     SvcCmdType = "start"
	SvcCmd_ReStart   SvcCmdType = "restart"
	SvcCmd_Stop      SvcCmdType = "stop"
)

func parseSvcCmd(cmd string) (svcCmd SvcCmdType, ok bool) {
	cmdMap := map[string]SvcCmdType{
		SvcCmd_Install.String():   SvcCmd_Install,
		SvcCmd_UnInstall.String(): SvcCmd_UnInstall,
		SvcCmd_Remove.String():    SvcCmd_Remove,
		SvcCmd_Start.String():     SvcCmd_Start,
		SvcCmd_ReStart.String():   SvcCmd_ReStart,
		SvcCmd_Stop.String():      SvcCmd_Stop,
	}

	svcCmd, ok = cmdMap[strings.ToLower(cmd)]
	return
}

// send command to SCM to manage service
func manageSvc(cmd SvcCmdType) error {

	log.Println("managing service: ", cmd)

	svc, err := newMyMainService()
	if err != nil {
		return err
	}

	log.Printf("%s service %v ...\n", cmd, svc)
	switch cmd {
	case SvcCmd_Install:
		err = svc.Install()
	case SvcCmd_UnInstall, SvcCmd_Remove:
		err = svc.Uninstall()
	case SvcCmd_Start:
		err = svc.Start()
	case SvcCmd_Stop:
		err = svc.Stop()
	case SvcCmd_ReStart:
		err = svc.Restart()
	}
	if err != nil {
		log.Printf("%s service %v failed!\n", cmd, svc)
	} else {
		log.Printf("%s service %v successfully!\n", cmd, svc)
	}

	return err
}
