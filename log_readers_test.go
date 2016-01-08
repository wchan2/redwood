package main_test

import (
	"log"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/wchan2/redwood"
)

var _ = Describe(`FileLogReader`, func() {
	var (
		fileLogReader LogReader
		file          *os.File
		err           error

		testFile = "sample.test"
	)

	BeforeEach(func() {
		file, err = os.Create(testFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		fileLogReader, err = NewLogFileReader(testFile)
		if err != nil {
			log.Fatal(err.Error())
		}
	})

	AfterEach(func() {
		fileLogReader.Close()
		file.Close()
		os.Remove(file.Name())
	})

	Describe(`#Read`, func() {
		Context(`when there are events in the file`, func() {
			BeforeEach(func() {
				file.WriteString(`209.160.2.63 - - [23/Dec/2015:18:22:21 -0700] "POST /cart.do?action=purchase&itemId=EST-21&JSESSIONID=SD0SL6FF7ADFF4953 HTTP/1.1" 200 486 "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.46 Safari/536.5"`)
			})

			It(`reads those events from the log reader`, func() {
				Eventually(fileLogReader.Read()).Should(Receive(Equal(Event{
					Client:      "209.160.2.63",
					User:        "",
					Identifier:  "",
					Time:        time.Date(2015, 12, 23, 18, 22, 21, 0, time.FixedZone("", -7*60*60)),
					Method:      "POST",
					Path:        "/cart.do?action=purchase&itemId=EST-21&JSESSIONID=SD0SL6FF7ADFF4953",
					Protocol:    "HTTP/1.1",
					StatusCode:  200,
					PayloadSize: 486,
					UserAgent:   "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.46 Safari/536.5",
				})))
			})
		})

		Context(`when events are simultaneously written to a file`, func() {
			BeforeEach(func() {
				go func() {
					time.Sleep(200 * time.Millisecond)
					file.WriteString(`209.160.24.64 - - [23/Dec/2015:12:00:00] "POST /order.do?action=purchase&itemId=EST-21&JSESSIONID=SD0SL6FF7ADFF4953 HTTP/1.1" 201 396 "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.46 Safari/536.5"`)
				}()
			})
			It(`reads those events from the log reader`, func() {
				Eventually(fileLogReader.Read()).Should(Receive(Equal(Event{
					Client:      "209.160.24.64",
					User:        "",
					Identifier:  "",
					Time:        time.Date(2015, 12, 23, 12, 0, 0, 0, time.UTC),
					Method:      "POST",
					Path:        "/order.do?action=purchase&itemId=EST-21&JSESSIONID=SD0SL6FF7ADFF4953",
					Protocol:    "HTTP/1.1",
					StatusCode:  201,
					PayloadSize: 396,
					UserAgent:   "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.46 Safari/536.5",
				})))
			})
		})
	})
})
