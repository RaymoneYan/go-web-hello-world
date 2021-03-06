package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><body bgcolor=yellow>Hello, %s %s %s<br>\n", r.Method, r.URL, r.Proto)
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header field %q, Value %q<br>\n", k, v)
		}
		fmt.Fprintf(w, "Host = %q<br>\n", r.Host)
		fmt.Fprintf(w, "RemoteAddr= %q<br>\n", r.RemoteAddr)
		fmt.Fprintf(w, "</body></html>", r.RemoteAddr)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
