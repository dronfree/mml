package main
import (
	"fmt"
	"net/http"
	"log"
	"html"
)


func main() {
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome from Mail boxes manager")
	})
	http.HandleFunc("/box", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "box: path %q, query %q", html.EscapeString(r.URL.Path), html.EscapeString(r.URL.RawQuery))
	})
	http.HandleFunc("/mails", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "mails: path %q, query %q", html.EscapeString(r.URL.Path), html.EscapeString(r.URL.RawQuery))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
