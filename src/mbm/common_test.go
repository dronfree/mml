package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
	"net"
	"log"
	"fmt"
	"time"
)

var _ = Describe("Common", func() {
	var page *agouti.Page

	BeforeEach(func() {
		var err error
		page, err = agoutiDriver.NewPage()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(page.Destroy()).To(Succeed())
	})

	It("should receive email", func() {

		By("accessing homepage", func() {
			Expect(page.Navigate("http://localhost/")).To(Succeed())
			Expect(page).To(HaveURL("http://localhost/"))
		})

		By("emulating sending email from remote mail server", func() {

			emailSelection := page.Find("#email")
			email, _ := emailSelection.Text()

			conn, err := net.Dial("tcp", "localhost:25")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Fprintf(conn, "MAIL FROM: user@local\r\n")
			fmt.Fprintf(conn, "RCPT TO: " + email + "\r\n")
			fmt.Fprintf(conn, "DATA\r\n")
			fmt.Fprintf(conn, "Subject: test subject\r\n")
			fmt.Fprintf(conn, "From: user@local\r\n")
			fmt.Fprintf(conn, "Mail body line 1\r\n")
			fmt.Fprintf(conn, ".\r\n")

			time.Sleep(time.Second * 5)
		})

		By("asserting email is displayed on the page", func() {
			Eventually(page.Find(".mailbox-read-message pre")).Should(HaveText("Mail body line 1"))
			Eventually(page.Find(".mailbox-read-info h3")).Should(HaveText("test subject"))
			Eventually(page.Find("#inboxAmount")).Should(HaveText("1"))
		})
	})
})
