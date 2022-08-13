package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
)

var (
	httpPort = flag.String("httpPort", "8080", "http server port")
)

func main() {
	flag.Parse()

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"name\": \"TheGamerGuildBot\"}")
	})

	http.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"name\": \"TheGamerGuildBot\"}")
	})

	log.Println("Http server started on port " + *httpPort)
	http.ListenAndServe(":"+*httpPort, nil)
}
