package main

import (
	"os"
	"bufio"
	"io/ioutil"
	"log"
	"flag"
)

type Params struct {
	mailboxBase string
	mailbox string
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
	var content []byte
	if content, err = ioutil.ReadAll(reader); err != nil {
		log.Fatal(err)
	}

	content = append(content, []byte("\n\nmy custom delimiter\n\n")...)
	box := params.mailboxBase + `/` + params.mailbox

	if outFile, err = os.OpenFile(box, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0666); err != nil {
		log.Fatal(err)
	}
	defer outFile.Close();
	var writer = bufio.NewWriter(outFile)
	if _, err = writer.Write(content); err != nil {
		log.Fatal(err)
	}
	writer.Flush()
}
