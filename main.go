package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func registerTempWelcome(tpl *template.Template) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := tpl.ExecuteTemplate(w, "main.tpl.html", nil)
		if err != nil {
			w.Write([]byte("An unknown error has occurred."))
		}
	})
}

func main() {
	tpl := template.Must(template.ParseGlob("templates/*.tpl.html"))
	registerTempWelcome(tpl)

	insecure := flag.Bool("i", false, "i(nsecure) mode (no TLS)")
	flag.Parse()

	var err error
	if *insecure {
		fmt.Println("Listing on HTTP/8080")
		err = http.ListenAndServe(":8080", nil)
	} else {
		var cert tls.Certificate
		cert, err = tls.LoadX509KeyPair(
			"/etc/letsencrypt/live/ufd.world/fullchain.pem",
			"/etc/letsencrypt/live/ufd.world/privkey.pem",
		)
		if err != nil {
			panic("X509 error")
		}

		server := &http.Server{
			Addr: ":443",
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}
		fmt.Println("Listing on HTTPS/443")
		err = server.ListenAndServeTLS("", "")
	}
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
