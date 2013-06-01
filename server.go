package hookworm

import (
	"flag"
	"log"
	"net/http"
)

var (
	addrFlag = flag.String("a", ":9988", "Server address")
)

type Server struct{}

func ServerMain() {
	server := &Server{}
	http.Handle("/", server)
	log.Printf("hookworm-server listening on %v", *addrFlag)
	log.Fatal(http.ListenAndServe(*addrFlag, nil))
}

func (me *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
}
