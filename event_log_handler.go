package hookworm

import (
	"encoding/json"
	"log"
)

type EventLogHandler struct {
	sysLogger   *log.Logger
	nextHandler Handler
}

func (me *EventLogHandler) HandlePayload(payload *Payload) error {
	asJsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to re-serialize payload:", err)
		return err
	}
	asJson := string(asJsonBytes)

	log.Printf("pull request merge? %v", payload.IsPullRequestMerge())
	log.Printf("payload=%+v", payload)
	log.Printf("payload json=%v", asJson)

	if me.sysLogger != nil {
		me.sysLogger.Println(asJson)
	}

	if nxt := me.NextHandler(); nxt != nil {
		return nxt.HandlePayload(payload)
	}

	return nil
}

func (me *EventLogHandler) SetNextHandler(handler Handler) {
	me.nextHandler = handler
}

func (me *EventLogHandler) NextHandler() Handler {
	return me.nextHandler
}
