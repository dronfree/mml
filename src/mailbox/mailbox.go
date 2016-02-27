package mailbox
import (
	"os"
	"io/ioutil"
	"log"
	"bufio"
	"net/mail"
	"mime"
	"strings"
	"mime/multipart"
	"io"
)

type JsonMail struct {
	Date string
	From string
	Subject string
	BodyText string
	BodyHtml string
	Body string
}


func Read(boxPath string) (mails []JsonMail, err error) {
	var (
		files []os.FileInfo
		file  os.FileInfo
		inFile *os.File
	)
	boxPath += `/new`
	if files, err = ioutil.ReadDir(boxPath); err != nil {
		log.Println(err)
		return
	}
	for _, file = range files {
		var msg *mail.Message
		var boxFile string
		var mediaType string
		var params map[string]string
		var htmlBody []byte
		var textBody []byte
		var defaultBody []byte

		boxFile = boxPath + `/` + file.Name()
		if inFile, err = os.Open(boxFile); err != nil {
			log.Println(err)
			return
		}
		defer inFile.Close()
		reader := bufio.NewReader(inFile)
		if msg, err = mail.ReadMessage(reader); err != nil {
			log.Println(err)
			return
		}
		header := msg.Header
		if mediaType, params, err = mime.ParseMediaType(header.Get("Content-Type")); err != nil {
			log.Println(`ERROR: parsing Content-Type of mail`, err)
			err = nil
		}
		if len(mediaType) != 0 && strings.HasPrefix(mediaType, "multipart/") {
			mr := multipart.NewReader(msg.Body, params["boundary"])
			for {
				var p *multipart.Part
				p, err = mr.NextPart()
				if err == io.EOF {
					log.Println(`NOTICE: End of multipart mail reached`)
					err = nil
					break
				}
				if err != nil {
					log.Println(`ERROR: getting next part of multipart mail`, err)
					return
				}
				if mediaType, params, err = mime.ParseMediaType(p.Header.Get("Content-Type")); err != nil {
					log.Println(`ERROR: parsing Content-Type of part`, err)
					return
				}
				if mediaType == "text/html" {
					htmlBody, err = ioutil.ReadAll(p)
					if err != nil {
						log.Println(`ERROR: reading html part`, err)
						return
					}
				}
				if mediaType == "text/plain" {
					textBody, err = ioutil.ReadAll(p)
					if err != nil {
						log.Println(`ERROR: reading text part`, err)
						return
					}
				}
			}
			if len(htmlBody) != 0 {
				defaultBody = htmlBody
			} else {
				defaultBody = textBody
			}
			mails = append(mails, JsonMail{file.ModTime().String(), header.Get("From"), header.Get("Subject"), string(textBody), string(htmlBody), string(defaultBody)})
		} else {
			body, err := ioutil.ReadAll(msg.Body)
			if err != nil {
				log.Println(`ERROR: reading non multipart mail body`, err)
			}
			b := `<pre>` + string(body) + `</pre>`
			mails = append(mails, JsonMail{file.ModTime().String(), header.Get("From"), header.Get("Subject"), b, "", b})
		}
	}
	return
}
