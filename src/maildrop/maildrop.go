package main

import (
	"os"
	"bufio"
	"log"
	"flag"
	"net/mail"
	"io/ioutil"
	"encoding/json"
)

type Params struct {
	mailboxBase string
	mailbox string
}
type JsonMail struct {
	Date string
	From string
	Subject string
	Body string
}
var params Params

func init() {
	flag.StringVar(&params.mailboxBase, "mailboxBase", "/var/www/boxes", "Mailbox base path")
	flag.StringVar(&params.mailbox, "mailbox", "general", "Mailbox (i.e. box01)")
}

func main() {
	flag.Parse()
	var err error
	var outFile = os.Stdout
	var inFile = os.Stdin
	var reader = bufio.NewReader(inFile)
	var msg *mail.Message
	var jsMail JsonMail

	if msg, err = mail.ReadMessage(reader); err != nil {
		log.Fatal(err)
	}

	header := msg.Header
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		log.Fatal(err)
	}

	jsMail = JsonMail{header.Get("Date"), header.Get("From"), header.Get("Subject"), string(body)}
	js, err := json.Marshal(jsMail)
	if err != nil {
		log.Fatal(err)
	}
	js = append(js, '\n')

	box := params.mailboxBase + `/` + params.mailbox

	if outFile, err = os.OpenFile(box, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0666); err != nil {
		log.Fatal(err)
	}
	defer outFile.Close();
	var writer = bufio.NewWriter(outFile)
	if _, err = writer.Write(js); err != nil {
		log.Fatal(err)
	}
	writer.Flush()
}
