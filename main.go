package main

import (
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
	err := http.ListenAndServeTLS(":443", "server.crt", "../private/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
