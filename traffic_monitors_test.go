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
					Host:        `example.com`,
					Path:        `/pages/reports`,
					PayloadSize: 20,
					Time:        currentTime,
				},
				{
					Host:        `example.com`,
					Path:        `/pages/orders`,
					PayloadSize: 10,
					Time:        currentTime.Add(1 * time.Second),
				},
				{
					Host:        `example2.com`,
					Path:        `/pages/items`,
					PayloadSize: 30,
					Time:        currentTime.Add(2 * time.Second),
				},
			}
		})

		JustBeforeEach(func() {
			for _, event := range events {
				trafficMonitor.Monitor(event)
			}
		})

		It(`calculates the summary statistics`, func() {
			destination1TotalPayloadSize := events[0].PayloadSize + events[1].PayloadSize
			destination1AvgPayloadSize := float64(destination1TotalPayloadSize) / 2

			destination2TotalPayloadSize := events[2].PayloadSize
			destination2AvgPayloadSize := float64(destination2TotalPayloadSize) / 1

			statisticsFormat := "Section: %s\nAverage Payload: %f\nTotal Payload: %d\nCount: %d\n"

			Eventually(func() string { return notification.message }, 2*time.Second, 500*time.Millisecond).Should(And(
				ContainSubstring(fmt.Sprintf(statisticsFormat, "example.com", destination1AvgPayloadSize, destination1TotalPayloadSize, 2)),
				ContainSubstring(fmt.Sprintf(statisticsFormat, "example2.com", destination2AvgPayloadSize, destination2TotalPayloadSize, 1)),
			))
		})
	})
})
