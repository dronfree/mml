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
	flag.StringVar(&params.boxpath, "boxpath", "/var/www/boxes", "Path to directory with stored boxes")
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

			availableBox := AvailableBox{box, sessId}
			if js, err = json.Marshal(availableBox); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if err = os.Remove(boxFile(box)); err != nil {
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
		var inFile *os.File
		var content []byte

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
			log.Println("sessid of box not match")
			return
		}
		if inFile, err = os.Open(boxFile(box)); err != nil {
			log.Println(err)
			return
		}
		defer inFile.Close()
		if content, err = ioutil.ReadAll(bufio.NewReader(inFile)); err != nil {
			log.Println(err)
			return
		}
		w.Write(content);
	})
	go expireBox()
	go makeFreeAvailable()
	http.ListenAndServe(":"+strconv.Itoa(params.port), nil)
}

func expireBox() {
	for {
		for i, k := range busy {
			if time.Now().After(k.expireAt)  {
				delete(busy, i)
				free = append(free, k.box)
				log.Println(k.box + " expired")
			}
		}
		time.Sleep(params.checkexpire)
	}
}

func makeFreeAvailable() {
	for {
		if len(free) > params.freecapacity {
			available = append(available, free[0:params.freecapacity-1]...)
			free = free[params.freecapacity:]
		}
		time.Sleep(params.makefreeavailable)
	}
}
