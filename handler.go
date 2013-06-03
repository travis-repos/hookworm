package hookworm

import (
	"log"
	"log/syslog"
)

type HandlerConfig struct {
	EmailUri        string
	EmailFromAddr   string
	EmailRcpts      []string
	UseSyslog       bool
	PolicedBranches []string
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

	if len(cfg.PolicedBranches) > 0 {
		sscHandler := &SecretSquirrelCommitHandler{
			emailer:         NewEmailer(cfg.EmailUri),
			fromAddr:        cfg.EmailFromAddr,
			recipients:      cfg.EmailRcpts,
			policedBranches: cfg.PolicedBranches,
		}
		elHandler.SetNextHandler(sscHandler)
		log.Printf("Added secret squirrel handler for policed branches %+v\n",
			cfg.PolicedBranches)
	} else {
		log.Println("No secret squirrel handler added")
	}

	return (Handler)(elHandler)
}
