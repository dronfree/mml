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
	"mailbox"
	"net/url"
	"errors"
	"io/ioutil"
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
	ExpiresIn int64
}
type AllBoxBusy struct {
	Error string
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
	flag.Int64Var(&params.rentfor, "rentfor", 3600, "Mailbox rent time in seconds")
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

			availableBox := AvailableBox{box, sessId, params.rentfor}
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
		var box string
		var err    error
		var content []mailbox.JsonMail
		var js []byte

		box, _, _, err = validateRequest(r.URL.Query(), busy)
		if err != nil {
			log.Println(err)
			return
		}

		if content, err = mailbox.Read(boxFile(box)); err != nil {
			log.Println(`ERROR: reading box content`, err)
			return
		}
		if js, err = json.Marshal(content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js);
	})
	http.HandleFunc("/mail", func (w http.ResponseWriter, r *http.Request) {
		var box, id string
		var err error
		var content []byte

		box, _, id, err = validateRequest(r.URL.Query(), busy)
		if err != nil {
			return
		}
		content, err = ioutil.ReadFile(boxFile(box) + "/new" + "/" + id)
		if err != nil {
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(content)
	})
	go expireBox()
	go makeFreeAvailable()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(params.port), nil))
}

func validateRequest(values url.Values, busyBoxes map[string]BusyBox) (box string, sessid string, id string, err error) {
	var (
		boxArr, sessidArr, idArr []string
	    ok bool
	    busyRow BusyBox
	)

	if boxArr, ok = values["box"]; !ok {
		return "", "", "", errors.New("box paramether not set")
	}
	if sessidArr, ok = values["sessid"]; !ok {
		return "", "", "", errors.New("sessid paramether not set")
	}
	box = boxArr[0]
	sessid = sessidArr[0]
	if busyRow, ok = busyBoxes[sessid]; !ok || busyRow.box != box {
		return "", "", "", errors.New("sessid or box not match")
	}
	idArr, ok = values["id"]
	if !ok {
		return box, sessid, "", nil
	}
	id = idArr[0]

	return box, sessid, id, nil
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
