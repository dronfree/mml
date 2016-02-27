package mailbox
import (
	"testing"
	"net/mail"
	"os"
	"bufio"
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
