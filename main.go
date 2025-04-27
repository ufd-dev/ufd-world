package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func TempWelcome(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
<!DOCTYPE html>
<html>
  <head>
    <title>UFD.World</title>
    <meta charset="UTF-8">
	<style type="text/css">
	  body {
	    margin: 0;
		padding: 0;
	    background-color: black;
		color: white;
		font-family: sans-serif;
		font-size: 1.5rem;
	  }
      a:link, a:visited, a:hover, a:active {
        color: #007bff; /* Your desired color */
        text-decoration: none;
      }
	  .container {
	    padding: 16px;
	  }
	</style>
  </head>
  <body>
    <div class="container">
      <h1>UFD.World</h1>
      <p>Coming soon: an unofficial community site for Unicorn Fart Dust</p>
      <p>
        The official site is
        <a href="https://unicornfartdust.com/">https://unicornfartdust.com/</a>.
      </p>
    </div>
  </body>
</html>
`))
}

func main() {
	http.HandleFunc("/", TempWelcome)

	insecure := flag.Bool("i", false, "i(nsecure) mode (no TLS)")
	flag.Parse()

	var err error
	if *insecure {
		fmt.Println("Listing on HTTP/8080")
		err = http.ListenAndServe(":8080", nil)
	} else {
		cert, err := tls.LoadX509KeyPair(
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
