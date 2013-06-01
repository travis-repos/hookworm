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
	addrFlag  = flag.String("a", ":9988", "Server address")
	emailFlag = flag.String("e", "smtp://localhost:25", "Email server address")

	logTimeFmt = "2/Jan/2006:15:04:05 -0700" // "%d/%b/%Y:%H:%M:%S %z"
)

type Server struct {
	emailer *Emailer
}

func ServerMain() {
	flag.Parse()

	server := &Server{emailer: &Emailer{serverUri: *emailFlag}}
	http.Handle("/", server)
	log.Printf("hookworm-server listening on %v", *addrFlag)
	log.Fatal(http.ListenAndServe(*addrFlag, nil))
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

	log.Printf("payload=%+v", payload)

	status = http.StatusNoContent
	w.WriteHeader(status)
}
