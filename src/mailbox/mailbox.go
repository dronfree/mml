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
	"errors"
	"encoding/base64"
	"golang.org/x/text/encoding/charmap"
	"fmt"
	"time"
)

type JsonMail struct {
	Id string
	Date string
	From string
	Subject string
	BodyText string
	BodyHtml string
	Body string
}


func IsMultiPart(msg *mail.Message) bool {
	mediaType, _, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if  err != nil {
		return false
	}
	return strings.HasPrefix(mediaType, `multipart/`)
}

func GetBoundary(msg *mail.Message) (boundary string, err error) {
	_, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}
	if _, ok := params["boundary"]; !ok {
		return "", errors.New("Boundary not found")
	}
	return params["boundary"], nil
}

func ReadMultiPartMail(msg *mail.Message) (email JsonMail, err error) {
	boundary, err := GetBoundary(msg)
	if err != nil {
		return JsonMail{}, err
	}
	mr := multipart.NewReader(msg.Body, boundary)
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(`ERROR: getting next part of multipart mail`, err)
			continue
		}
		mediaType, params, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
		if err != nil {
			log.Println(`ERROR: parsing Content-Type of part`, err)
			continue
		}

		var bodyReader io.Reader = p;
		if charset, ok := params["charset"]; ok {
			if ("koi8-r" == charset) {
				bodyReader = charmap.KOI8R.NewDecoder().Reader(p)
			}
		}

		if mediaType == "text/html" {
			htmlBody, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				log.Println(`ERROR: reading html part`, err)
				continue
			}
			if "base64" == p.Header.Get("Content-Transfer-Encoding") {
				htmlBody, err = base64.StdEncoding.DecodeString(string(htmlBody))
				if err != nil {
					log.Println(`ERROR: decoding base64`, err)
					continue
				}
			}
			email.BodyHtml = strings.Trim(string(htmlBody), "\n")
		}
		if mediaType == "text/plain" {
			textBody, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				log.Println(`ERROR: reading text part`, err)
				continue
			}
			if "base64" == p.Header.Get("Content-Transfer-Encoding") {
				textBody, err = base64.StdEncoding.DecodeString(string(textBody))
				if err != nil {
					log.Println(`ERROR: decoding base64`, err)
					continue
				}
			}
			email.BodyText = strings.Trim(string(textBody), "\r\n")
		}
	}
	if len(email.BodyHtml) != 0 {
		email.Body = email.BodyHtml
	} else {
		email.Body = email.BodyText
	}
	email.From = msg.Header.Get("From")
	email.Subject = msg.Header.Get("Subject")
	decoder := new(mime.WordDecoder)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if "koi8-r" == charset {
			return charmap.KOI8R.NewDecoder().Reader(input), nil
		}
		return nil, fmt.Errorf("unhandled charset for subject %q", charset)
	}
	decodedSubject, err := decoder.Decode(email.Subject)
	if err == nil {
		email.Subject = decodedSubject
	}
	return email, nil
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

		boxFile = boxPath + `/` + file.Name()
		inFile, err = os.Open(boxFile)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		defer inFile.Close()

		reader := bufio.NewReader(inFile)
		msg, err := mail.ReadMessage(reader)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		header := msg.Header
		var email JsonMail
		if IsMultiPart(msg) {
			email, err = ReadMultiPartMail(msg)
			if err != nil {
				log.Println(err)
				continue
			}
			email.Date = file.ModTime().Format(time.UnixDate)
			email.Id = file.Name()
		} else {
			body, err := ioutil.ReadAll(msg.Body)
			if err != nil {
				log.Println(`ERROR: reading non multipart mail body`, err)
			}
			b := `<pre>` + string(body) + `</pre>`
			email = JsonMail{file.Name(), file.ModTime().Format(time.UnixDate), header.Get("From"), header.Get("Subject"), b, "", b}
		}
		mails = append(mails, email)
	}
	return
}
