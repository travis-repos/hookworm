package hookworm

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"time"
)

var (
	addrFlag      = flag.String("a", ":9988", "Server address")
	emailFlag     = flag.String("e", "smtp://localhost:25", "Email server address")
	useSyslogFlag = flag.Bool("S", false, "Send all received events to syslog")

	logTimeFmt = "2/Jan/2006:15:04:05 -0700" // "%d/%b/%Y:%H:%M:%S %z"
)

type Server struct {
	emailer   *Emailer
	useSyslog bool
	sysLogger *log.Logger
}

func ServerMain() {
	flag.Parse()

	server := NewServer(*emailFlag, *useSyslogFlag)
	if server == nil {
		log.Fatal("No server?  No worky!")
	}
	http.Handle("/", server)
	log.Printf("hookworm-server listening on %v", *addrFlag)
	log.Fatal(http.ListenAndServe(*addrFlag, nil))
}

func NewServer(emailUri string, useSyslog bool) *Server {
	var err error
	server := &Server{
		emailer:   &Emailer{serverUri: emailUri},
		useSyslog: useSyslog,
	}
	if useSyslog {
		server.sysLogger, err = syslog.NewLogger(syslog.LOG_INFO, log.LstdFlags)
		if err != nil {
			log.Println("Failed to initialize syslogger!", err)
			return nil
		}
		log.Println("Initialized syslog logger")
	}

	return server
}

func (me *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusBadRequest

	defer func() {
		fmt.Fprintf(os.Stderr, "%s - - [%s] \"%s %s %s\" %d -\n",
			r.RemoteAddr, time.Now().UTC().Format(logTimeFmt),
			r.Method, r.URL.Path, r.Proto, status)
	}()

	rawPayload := r.FormValue("payload")
	if len(rawPayload) < 1 {
		log.Println("Empty payload!")
		w.WriteHeader(status)
		return
	}

	payload := &Payload{}
	err := json.Unmarshal([]byte(rawPayload), payload)
	if err != nil {
		log.Println("Failed to unmarshal payload: ", err)
		w.WriteHeader(status)
		return
	}

	me.logEvent(payload)

	status = http.StatusNoContent
	w.WriteHeader(status)
}

func (me *Server) logEvent(payload *Payload) {
	asJsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to re-serialize payload:", err)
		return
	}
	asJson := string(asJsonBytes)

	log.Printf("payload=%+v", payload)
	log.Printf("payload json=%v", asJson)

	if me.useSyslog && me.sysLogger != nil {
		me.sysLogger.Println(asJson)
	}
}
