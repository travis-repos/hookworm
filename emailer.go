package hookworm

import (
	"net/smtp"
	"net/url"
)

type Emailer struct {
	serverUri       string
	parsedServerUri *url.URL
}

func (me *Emailer) Send(from string, to []string, msg []byte) error {
	return smtp.SendMail(me.addr(), me.auth(), from, to, msg)
}

func (me *Emailer) addr() string {
	return ""
}

func (me *Emailer) auth() smtp.Auth {
	return nil
}
