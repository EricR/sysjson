package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"strings"
)

var (
	listen   = flag.String("listen", ":5374", "Address to listen on")
	tls      = flag.Bool("tls", false, "Use TLS (requires -cert and -key)")
	cert     = flag.String("cert", "", "TLS cert file")
	key      = flag.String("key", "", "TLS key file")
	password = flag.String("password", "", "Enable basic authentication")
	whitelist	=	flag.String("whitelist", "", "CIDR whitelist for requests")
)

func main() {
	flag.Parse()
	
	if len(*whitelist) > 0 {
		_, _, err := net.ParseCIDR(*whitelist)
		if err != nil {
			log.Fatal("Unable to parse CIDR address")
		}
		log.Printf("[notice] sys.json request whitelist set to %s", *whitelist)
	}

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
	resp := map[string]interface{}{}
	
	requestIpPort := strings.Split(r.RemoteAddr, ":")
	requestIp := net.ParseIP(requestIpPort[0])
	_, whitelistedNet, _ := net.ParseCIDR(*whitelist)

	if len(*whitelist) == 0 || (len(*whitelist) > 0 && whitelistedNet.Contains(requestIp)){
		loadModules(resp, r.URL.Query().Get("modules"))
	} else {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
		
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
			http.Error(w, "authentication is required", http.StatusUnauthorized)
			return
		}

		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)

		if auth[0] != "Basic" || len(auth) != 2 {
			http.Error(w, "bad syntax", http.StatusBadRequest)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		parsed := string(payload)

		if strings.Contains(parsed, ":") {
			pair := strings.SplitN(string(payload), ":", 2)
			parsed = pair[1]
		}

		if !Validate(parsed) {
			http.Error(w, "authentication failed", http.StatusUnauthorized)
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
