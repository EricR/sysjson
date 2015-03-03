package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	listen   = flag.String("listen", ":5374", "Address to listen on")
	tls      = flag.Bool("tls", false, "Use TLS (requires -cert and -key)")
	cert     = flag.String("cert", "", "TLS cert file")
	key      = flag.String("key", "", "TLS key file")
	password = flag.String("password", "", "Enable basic authentication")
)

func main() {
	flag.Parse()
	log.Printf("[notice] sys.json listening on %s", *listen)

	mux := http.NewServeMux()

	if len(*password) > 0 {
		mux.HandleFunc("/", BasicAuth(statsHandler))
	} else {
		mux.HandleFunc("/", statsHandler)
	}

	if *tls {
		log.Printf("[notice] Using TLS")
		log.Fatal(http.ListenAndServeTLS(*listen, *cert, *key, mux))
	} else {
		log.Fatal(http.ListenAndServe(*listen, mux))
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	resp := j{}
	loadTime := time.Now()
	resp["current_time"] = j{
		"string": loadTime,
		"unix":   loadTime.Unix(),
	}

	hostname, _ := os.Hostname()
	resp["hostname"] = hostname

	loadModules(resp, r.URL.Query().Get("modules"))

	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("[error] Fatal! Could not construct JSON response: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func BasicAuth(pass http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if len(r.Header.Get("Authorization")) <= 0 {
			http.Error(w, "authorization is required", http.StatusUnauthorized)
			return
		}

		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "bad syntax", http.StatusBadRequest)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if !Validate(pair[1]) {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		pass(w, r)
	}
}

func Validate(pass string) bool {
	if pass == *password {
		return true
	}
	return false
}
