package main
import (
	"fmt"
	"net/http"
	"log"
)


func main() {
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome from Mail boxes manager")
	})
	http.HandleFunc("/box", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "box")
	})
	http.HandleFunc("/mails", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "mails")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
