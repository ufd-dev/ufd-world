package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	insecure := flag.Bool("i", false, "i(nsecure) mode (no TLS)")
	flag.Parse()

	loadTemplates()
	r := configRoutes()

	server := &http.Server{
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	var err error
	if *insecure {
		server.Addr = ":8080"
		fmt.Println("Listing on HTTP/8080")
		err = server.ListenAndServe()
	} else {
		var cert tls.Certificate
		cert, err = tls.LoadX509KeyPair(
			"/etc/letsencrypt/live/ufd.world/fullchain.pem",
			"/etc/letsencrypt/live/ufd.world/privkey.pem",
		)
		if err != nil {
			panic("X509 error")
		}

		server.Addr = ":443"
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		fmt.Println("Listing on HTTPS/443")
		err = server.ListenAndServeTLS("", "")
	}
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
