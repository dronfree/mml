package main
import (
	"fmt"
	"net/http"
	"log"
	"flag"
	"os"
	"bufio"
	"io"
	"strconv"
	"os/exec"
	"strings"
	"time"
	"encoding/json"
	"io/ioutil"
	"net/mail"
	"mime"
	"mime/multipart"
)

type Params struct {
	port      int
	mailboxes string
	rentfor   int64
	boxpath   string
	checkexpire time.Duration
	freecapacity int
	makefreeavailable time.Duration
}
type BusyBox struct {
	box    string
	expireAt time.Time
}
type AvailableBox struct {
	Box    string
	Sessid string
}
type AllBoxBusy struct {
	Error string
}
type JsonMail struct {
	Date string
	From string
	Subject string
	BodyText string
	BodyHtml string
}
var params Params
var (
	available []string
	free      []string
	busy      = make(map[string]BusyBox)
)

func init() {
	flag.IntVar(&params.port, "port", 8080, "Port to start app on")
	flag.DurationVar(&params.checkexpire, "checkexpire", 5*time.Second, "How often to perform check expire boxes in seconds")
	flag.IntVar(&params.freecapacity, "freecapacity", 5, "Max number of expired boxes to return to queue")
	flag.DurationVar(&params.makefreeavailable, "makefreeavailable", 5*time.Second, "How often to perform makefreeavailable")
	flag.StringVar(&params.mailboxes, "mailboxes", "vmailbox", "Postfix virtual map file")
	flag.StringVar(&params.boxpath, "boxpath", "boxes", "Path to directory with stored boxes")
	flag.Int64Var(&params.rentfor, "rentfor", 300, "Mailbox rent time in seconds")
}

func SessId() string {
	var (
		uuid []byte
		err error
	)
	if uuid, err = exec.Command("uuidgen").Output(); err != nil {
		log.Fatal(err)
	}
	return strings.Trim(string(uuid), "\n")
}

func boxFile(box string) string {
	return params.boxpath + "/" + strings.Split(box, "@")[0]
}

func main() {
	var (
		err    error
		inFile *os.File
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
		lineArr := strings.Split(line, " ")
		available = append(available, lineArr[0])
	}

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome from Mail boxes manager")
	})
	http.HandleFunc("/box", func (w http.ResponseWriter, r *http.Request) {
		var box string
		var err error
		var js []byte
		if len(available) == 0 {
			if js, err = json.Marshal(AllBoxBusy{"allbusy"}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			box, available = available[0], available[1:len(available)]
			sessId := SessId()
			busy[sessId] = BusyBox{box, time.Now().Add(time.Duration(1e9 * params.rentfor))}

			availableBox := AvailableBox{box, sessId}
			if js, err = json.Marshal(availableBox); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if err = os.RemoveAll(boxFile(box)); err != nil {
				log.Println(err)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js);
	})
	http.HandleFunc("/mails", func (w http.ResponseWriter, r *http.Request) {
		var boxArr, sessidArr []string
		var box, sessid string
		var busyRow BusyBox
		var ok bool
		var err    error
		var content []JsonMail
		var js []byte

		if boxArr, ok = r.URL.Query()["box"]; !ok {
			log.Println("box paramether not set")
			return
		}
		if sessidArr, ok = r.URL.Query()["sessid"]; !ok {
			log.Println("sessid paramether not set")
			return
		}
		box = boxArr[0]
		sessid = sessidArr[0]

		if busyRow, ok = busy[sessid]; !ok || busyRow.box != box {
			log.Println("sessid or box not match")
			return
		}
		if content, err = readBoxContent(boxFile(box)); err != nil {
			log.Println(err)
			return
		}
		if js, err = json.Marshal(content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js);
	})
	go expireBox()
	go makeFreeAvailable()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(params.port), nil))
}

func readBoxContent(boxPath string) (mails []JsonMail, err error) {
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
			log.Println(err)
			return
		}
		if strings.HasPrefix(mediaType, "multipart/") {
			mr := multipart.NewReader(msg.Body, params["boundary"])
			for {
				var p *multipart.Part
				p, err = mr.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Println(err)
					return
				}
				if mediaType, params, err = mime.ParseMediaType(p.Header.Get("Content-Type")); err != nil {
					log.Println(err)
					return
				}
				if mediaType == "text/html" {
					htmlBody, err = ioutil.ReadAll(p)
					if err != nil {
						log.Println(err)
						return
					}
				}
				if mediaType == "text/plain" {
					textBody, err = ioutil.ReadAll(p)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
			mails = append(mails, JsonMail{file.ModTime().String(), header.Get("From"), header.Get("Subject"), string(textBody), string(htmlBody)})
		} else {
			body, err := ioutil.ReadAll(msg.Body)
			if err != nil {
				log.Println(err)
			}
			mails = append(mails, JsonMail{file.ModTime().String(), header.Get("From"), header.Get("Subject"), string(body), ""})
		}
	}
	return
}

func expireBox() {
	for {
		for i, k := range busy {
			if time.Now().After(k.expireAt)  {
				delete(busy, i)
				free = append(free, k.box)
				log.Println("expired: " + k.box)
			}
		}
		time.Sleep(params.checkexpire)
	}
}

func makeFreeAvailable() {
	for {
		if len(free) > params.freecapacity {
			becomeAvailable := free[0:params.freecapacity-1];
			available = append(available, becomeAvailable...)
			free = free[params.freecapacity:]
			log.Println("become available: ", becomeAvailable)
		}
		time.Sleep(params.makefreeavailable)
	}
}
