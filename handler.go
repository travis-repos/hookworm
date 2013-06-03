package hookworm

import (
	"log"
	"log/syslog"
)

type HandlerConfig struct {
	EmailUri       string
	EmailRcpts     []string
	UseSyslog      bool
	StableBranches []string
}

type Handler interface {
	HandlePayload(*Payload) error
	SetNextHandler(Handler)
	NextHandler() Handler
}

func NewHandlerPipeline(cfg *HandlerConfig) Handler {
	var err error
	elHandler := &EventLogHandler{}

	if cfg.UseSyslog {
		elHandler.sysLogger, err = syslog.NewLogger(syslog.LOG_INFO, log.LstdFlags)
		if err != nil {
			log.Panicln("Failed to initialize syslogger!", err)
		}
		log.Println("Added syslog logger to event handler")
	} else {
		log.Println("No syslog logger added to event handler")
	}

	if len(cfg.StableBranches) > 0 {
		sscHandler := &SecretSquirrelCommitHandler{
			emailer:        &Emailer{serverUri: cfg.EmailUri},
			recipients:     cfg.EmailRcpts,
			stableBranches: cfg.StableBranches,
		}
		elHandler.SetNextHandler(sscHandler)
		log.Printf("Added secret squirrel handler for stable branches %+v",
			cfg.StableBranches)
	} else {
		log.Println("No secret squirrel handler added")
	}

	return (Handler)(elHandler)
}
