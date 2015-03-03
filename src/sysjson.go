package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	listen = flag.String("listen", ":5374", "Address to listen on")
	tls    = flag.Bool("tls", false, "Use TLS (requires -cert and -key)")
	cert   = flag.String("cert", "", "TLS cert file")
	key    = flag.String("key", "", "TLS key file")
	from	= flag.String("from", "", "Address to accept connections from")
)

func main() {
	flag.Parse()

	if len(*from) > 0 {
		ip := net.ParseIP(*from)
		
		if ip == nil {
			log.Fatal("Invalid textual representation of an IP address provided in -from")
		}
		
		log.Printf("[notice] sys.json will accept connections from %s", *from)	
	}

	log.Printf("[notice] sys.json listening on %s", *listen)

	mux := http.NewServeMux()
	mux.HandleFunc("/", statsHandler)
	
	if *tls {
		log.Printf("[notice] Using TLS")
		log.Fatal(http.ListenAndServeTLS(*listen, *cert, *key, mux))
	} else {
		log.Fatal(http.ListenAndServe(*listen, mux))
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	resp := j{}

	if len(*from) == 0 || (len(*from) > 0 && strings.HasPrefix(r.RemoteAddr, *from)){
		loadTime := time.Now()
		resp["current_time"] = j{
		"string": loadTime,
		"unix":   loadTime.Unix(),
		}
	
		hostname, _ := os.Hostname()
		resp["hostname"] = hostname
		
		loadModules(resp, r.URL.Query().Get("modules"))
	} else {
		resp["error"] = "Not authorised"
	}
		
	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("[error] Fatal! Could not construct JSON response: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}
