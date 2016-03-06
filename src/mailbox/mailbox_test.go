package mailbox
import (
	"testing"
	"net/mail"
	"os"
	"bufio"
	"io/ioutil"
	"strings"
)

func getMailFromFile(filePath string) (*mail.Message, error) {
	inFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	msg, err := mail.ReadMessage(reader)
	if err != nil {
		return nil, err
	}
	// reading mail body in memory before closing file
	// because it cases bad file descriptor error
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		return nil, err
	}
	// TODO: find a way to make io.Reader from []byte
	msg.Body = strings.NewReader(string(body))
	return msg, nil
}

func TestIsMultiPart01(t *testing.T) {
	eml := "./testdata/multipart.eml"
	mail, _ := getMailFromFile(eml)
	result := IsMultiPart(mail)
	if !result {
		t.Errorf("IsMultipart(%q) == %t, want %t", eml, result, true)
	}
}

func TestIsMultiPart02(t *testing.T) {
	eml := "./testdata/nonmultipart.eml"
	mail, _ := getMailFromFile(eml)
	result := IsMultiPart(mail)
	if result {
		t.Errorf("IsMultipart(%q) == %t, want %t", eml, result, false)
	}
}

func TestGetBoundary01(t *testing.T) {
	eml := "./testdata/multipart.eml"
	mail, _ := getMailFromFile(eml)
	boundary, err := GetBoundary(mail)
	if boundary == "" || err != nil {
		t.Errorf("GetBoundary(%q) == %q, want %q, err %v", eml, boundary, "<non empty string>", err)
	}
}

func TestGetBoundary02(t *testing.T) {
	eml := "./testdata/nonmultipart.eml"
	mail, _ := getMailFromFile(eml)
	boundary, err := GetBoundary(mail)
	if boundary != "" || err == nil {
		t.Errorf("GetBoundary(%q) == %q, want %q, err %v", eml, boundary, "<empty string>", err)
	}
}

func TestGetBoundary03(t *testing.T) {
	eml := "./testdata/multipart-no-boundary.eml"
	mail, _ := getMailFromFile(eml)
	boundary, err := GetBoundary(mail)
	if boundary != "" || err == nil {
		t.Errorf("GetBoundary(%q) == %q, want %q, err %v", eml, boundary, "<empty string>", err)
	}
}

func TestReadMultiPartMail01(t *testing.T) {

	eml := "./testdata/multipart.eml"
	mail, _ := getMailFromFile(eml)
	json, err := ReadMultiPartMail(mail)
	if err != nil {
		t.Errorf("ReadMultiPartMail(%q) == ERROR, err %v", eml, err)
	}
	master := JsonMail{
		"Thu, 25 Feb 2016 20:15:28 +0300",
		"User Name <user@gmail.com>",
	    "test subject",
`test mail body

regards,
alex`,
		`<div dir="ltr">test mail body<div><br></div><div>regards,</div><div>alex</div></div>`,
		`<div dir="ltr">test mail body<div><br></div><div>regards,</div><div>alex</div></div>`,
	}

	if json.Date != master.Date {
		t.Errorf("ReadMultiPartMail(%q) returned JsonMail.Date == %q, want %q", eml, json.Date, master.Date)
	}
	if json.From != master.From {
		t.Errorf("ReadMultiPartMail(%q) returned JsonMail.From == %q, want %q", eml, json.From, master.From)
	}
	if json.Subject != master.Subject {
		t.Errorf("ReadMultiPartMail(%q) returned JsonMail.Subject == %q, want %q", eml, json.Subject, master.Subject)
	}
	if json.BodyHtml != master.BodyHtml {
		t.Errorf("ReadMultiPartMail(%q) returned JsonMail.BodyHtml == %q, want %q", eml, json.BodyHtml, master.BodyHtml)
	}
	if json.BodyText != master.BodyText {
		t.Errorf("ReadMultiPartMail(%q) returned JsonMail.BodyText == %q, want %q", eml, json.BodyText, master.BodyText)
	}
	if json.Body != master.Body {
		t.Errorf("ReadMultiPartMail(%q) returned JsonMail.Body == %q, want %q", eml, json.Body, master.Body)
	}

}

func TestReadMultiPartMail02(t *testing.T) {

	eml := "./testdata/multipart-plain-text-only.eml"
	mail, _ := getMailFromFile(eml)
	json, err := ReadMultiPartMail(mail)
	if err != nil {
		t.Errorf("ReadMultiPartMail(%q) == ERROR, err %v", eml, err)
	}
	master := JsonMail{
		"Thu, 25 Feb 2016 20:15:28 +0300",
		"User Name <user@gmail.com>",
	    "test subject",
`test mail body

regards,
alex`,
		``,
`test mail body

regards,
alex`,
	}

	if json.Date != master.Date {
		t.Errorf("ReadMultiPartMail(%q) returned JsonMail.Date == %q, want %q", eml, json.Date, master.Date)
	}
	if json.From != master.From {
		t.Errorf("ReadMultiPartMail(%q) returned JsonMail.From == %q, want %q", eml, json.From, master.From)
	}
	if json.Subject != master.Subject {
		t.Errorf("ReadMultiPartMail(%q) returned JsonMail.Subject == %q, want %q", eml, json.Subject, master.Subject)
	}
	if json.BodyHtml != master.BodyHtml {
		t.Errorf("ReadMultiPartMail(%q) returned JsonMail.BodyHtml == %q, want %q", eml, json.BodyHtml, master.BodyHtml)
	}
	if json.BodyText != master.BodyText {
		t.Errorf("ReadMultiPartMail(%q) returned JsonMail.BodyText == %q, want %q", eml, json.BodyText, master.BodyText)
	}
	if json.Body != master.Body {
		t.Errorf("ReadMultiPartMail(%q) returned JsonMail.Body == %q, want %q", eml, json.Body, master.Body)
	}

}


