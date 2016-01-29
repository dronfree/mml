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
	log.Fatal(http.ListenAndServe(":8080", nil))
}
