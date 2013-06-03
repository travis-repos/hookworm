package hookworm

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	addrFlag           = flag.String("a", ":9988", "Server address")
	emailFlag          = flag.String("e", "smtp://localhost:25", "Email server address")
	emailRcptsFlag     = flag.String("r", "", "Email recipients (comma-delimited)")
	stableBranchesFlag = flag.String("b", "", "Stable branches (comma-delimited)")
	useSyslogFlag      = flag.Bool("S", false, "Send all received events to syslog")
	printVersionFlag   = flag.Bool("v", false, "Print version and exit")

	logTimeFmt = "2/Jan/2006:15:04:05 -0700" // "%d/%b/%Y:%H:%M:%S %z"
)

type Server struct {
	pipeline Handler
}

func ServerMain() {
	flag.Parse()
	if *printVersionFlag {
		printVersion()
		return
	}

	cfg := &HandlerConfig{
		EmailUri:       *emailFlag,
		UseSyslog:      *useSyslogFlag,
		EmailRcpts:     commaSplit(*emailRcptsFlag),
		StableBranches: commaSplit(*stableBranchesFlag),
	}

	server := NewServer(cfg)
	if server == nil {
		log.Fatal("No server?  No worky!")
	}

	http.Handle("/", server)
	log.Printf("Listening on %v", *addrFlag)
	log.Fatal(http.ListenAndServe(*addrFlag, nil))
}

func NewServer(cfg *HandlerConfig) *Server {
	return &Server{pipeline: NewHandlerPipeline(cfg)}
}

func (me *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusBadRequest

	defer func() {
		w.WriteHeader(status)
		fmt.Fprintf(os.Stderr, "%s - - [%s] \"%s %s %s\" %d -\n",
			r.RemoteAddr, time.Now().UTC().Format(logTimeFmt),
			r.Method, r.URL.Path, r.Proto, status)
	}()

	rawPayload := r.FormValue("payload")
	if len(rawPayload) < 1 {
		log.Println("Empty payload!")
		return
	}

	payload := &Payload{}
	err := json.Unmarshal([]byte(rawPayload), payload)
	if err != nil {
		log.Println("Failed to unmarshal payload: ", err)
		log.Println("Raw payload: ", rawPayload)
		return
	}

	status = http.StatusNoContent

	if me.pipeline == nil {
		return
	}

	err = me.pipeline.HandlePayload(payload)
	if err != nil {
		status = http.StatusInternalServerError
		return
	}
}
