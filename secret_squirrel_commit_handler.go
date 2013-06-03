package hookworm

import (
	"fmt"
	"log"
	"sort"
)

type SecretSquirrelCommitHandler struct {
	emailer        *Emailer
	recipients     []string
	stableBranches []string
	nextHandler    Handler
}

func (me *SecretSquirrelCommitHandler) HandlePayload(payload *Payload) error {
	me.checkIfSecretSquirrelCommit(payload)
	return nil
}

func (me *SecretSquirrelCommitHandler) SetNextHandler(handler Handler) {
	me.nextHandler = handler
}

func (me *SecretSquirrelCommitHandler) NextHandler() Handler {
	return me.nextHandler
}

func (me *SecretSquirrelCommitHandler) checkIfSecretSquirrelCommit(payload *Payload) {
	if !me.isStableBranch(payload.Ref.String()) {
		return
	}

	if payload.IsPullRequestMerge() {
		return
	}

	me.alertSecretSquirrelCommit(payload)
}

func (me *SecretSquirrelCommitHandler) isStableBranch(ref string) bool {
	return sort.SearchStrings(me.stableBranches, ref) > -1
}

func (me *SecretSquirrelCommitHandler) alertSecretSquirrelCommit(payload *Payload) {
	// FIXME use the emailer thing here
	log.Printf("WARNING secret squirrel commit! %+v\n", payload)
	if len(me.recipients) == 0 {
		log.Println("No email recipients specified, so no emailing!")
	}

	me.emailer.Send("hookworm@localhost", me.recipients,
		[]byte(fmt.Sprintf("Secret squirrel commit!\n%+v\n", payload)))
}
