package main
import (
	"fmt"
	"net/http"
	"log"
	"html"
	"flag"
	"os"
	"bufio"
	"io"
	"strconv"
	"os/exec"
	"strings"
)

type Params struct {
	port      int
	mailboxes string
	rentfor   int
}
var params Params

func init() {
	flag.IntVar(&params.port, "port", 8080, "Port to start app on")
	flag.StringVar(&params.mailboxes, "mailboxes", "vmailbox", "Postfix virtual map file")
	flag.IntVar(&params.rentfor, "rentfor", 3600, "Mailbox rent time in seconds")
}

func main() {
	var (
		err    error
		inFile *os.File
		uuid   []byte
		sessId string
	)
	flag.Parse()
	if inFile, err = os.Open(params.mailboxes); err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()
	reader := bufio.NewReader(inFile)
	eof := false
	for !eof {
		var line string
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			eof = true
			continue
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Print(line)
	}

	fmt.Println(params.mailboxes)
	fmt.Println(params.port)
	fmt.Println(params.rentfor)
	uuid, err = exec.Command("uuidgen").Output()
	sessId = strings.Trim(string(uuid), "\n")
	fmt.Println(sessId)


	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome from Mail boxes manager")
	})
	http.HandleFunc("/box", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "box: path %q, query %q", html.EscapeString(r.URL.Path), html.EscapeString(r.URL.RawQuery))
	})
	http.HandleFunc("/mails", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "mails: path %q, query %q", html.EscapeString(r.URL.Path), html.EscapeString(r.URL.RawQuery))
	})
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(params.port), nil))
}
