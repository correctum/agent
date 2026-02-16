package main

import (
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sys/windows/svc"
)

type service struct{}

func (m *service) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	var accepts = svc.AcceptStop | svc.AcceptPauseAndContinue | svc.AcceptPreShutdown | svc.AcceptShutdown
	tick := time.Tick(5 * time.Second)
	status <- svc.Status{State: svc.Running, Accepts: accepts}

loop:
	for {
		select {
		case <-tick:
			log.Print("Tick Handled...!")
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				log.Print("Shutting service...!")
				break loop
			case svc.Pause:
				status <- svc.Status{State: svc.Paused, Accepts: accepts}
			case svc.Continue:
				status <- svc.Status{State: svc.Running, Accepts: accepts}
			default:
				log.Printf("Unexpected service control request #%d", c)
			}
		}
	}

	status <- svc.Status{State: svc.StopPending}
	return false, 1
}
