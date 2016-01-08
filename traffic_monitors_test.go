package main_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/wchan2/redwood"
)

var _ = Describe(`SummaryStatsTrafficMonitor`, func() {
	Describe(`#Monitor`, func() {
		var (
			trafficMonitor TrafficMonitor
			notification   *notificationMock

			events []Event
		)

		BeforeEach(func() {
			notification = new(notificationMock)
			trafficMonitor = NewSummaryStatsTrafficMonitor(1*time.Second, notification)
			currentTime := time.Now()
			events = []Event{
				{
					Path:        `/section1/page1`,
					PayloadSize: 20,
					Time:        currentTime,
					StatusCode:  200,
				},
				{
					Path:        `/section1/page2`,
					PayloadSize: 10,
					Time:        currentTime.Add(1 * time.Second),
					StatusCode:  301,
				},
				{
					Path:        `/section1/page1`,
					PayloadSize: 10,
					Time:        currentTime.Add(1 * time.Second),
					StatusCode:  404,
				},
				{
					Path:        `/section2/page1`,
					PayloadSize: 30,
					Time:        currentTime.Add(2 * time.Second),
					StatusCode:  201,
				},
			}
		})

		JustBeforeEach(func() {
			for _, event := range events {
				trafficMonitor.Monitor(event)
			}
		})

		It(`calculates the summary statistics`, func() {
			// calculate section 1 statistics
			section1TotalPayloadSize := events[0].PayloadSize + events[1].PayloadSize + events[2].PayloadSize
			section1AvgPayloadSize := float64(section1TotalPayloadSize) / 3.0
			section1Successes := 1
			section1Redirects := 1
			section1ClientFailures := 1
			section1ServerFailures := 0

			// calculate section 2 statistics
			section2TotalPayloadSize := events[3].PayloadSize
			section2AvgPayloadSize := float64(section2TotalPayloadSize) / 1.0
			section2Successes := 1
			section2Redirects := 0
			section2ClientFailures := 0
			section2ServerFailures := 0

			// calculate total traffic statistics
			totalPayloadSize := events[0].PayloadSize + events[1].PayloadSize + events[2].PayloadSize + events[3].PayloadSize
			avgTotalPayloadSize := float64(totalPayloadSize) / 4.0
			totalSuccesses := 2
			totalClientFailures := 1
			totalServerFailures := 0
			totalRedirects := 1

			Eventually(func() string { return notification.message }, 2*time.Second, 500*time.Millisecond).Should(And(
				ContainSubstring(fmt.Sprintf(
					SummaryStatisticsFormat, "/section1", section1AvgPayloadSize, section1TotalPayloadSize, section1Successes, section1Redirects, section1ClientFailures, section1ServerFailures, 3,
				)),
				ContainSubstring(fmt.Sprintf(
					SummaryStatisticsFormat, "/section2", section2AvgPayloadSize, section2TotalPayloadSize, section2Successes, section2Redirects, section2ClientFailures, section2ServerFailures, 1),
				),
				ContainSubstring(
					fmt.Sprintf(SummaryStatisticsFormat, "Total Traffic", avgTotalPayloadSize, totalPayloadSize, totalSuccesses, totalRedirects, totalClientFailures, totalServerFailures, 4),
				),
			))
		})
	})
})
