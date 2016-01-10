package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

const SummaryStatisticsFormat = `
Section: %s
Average Payload: %f
Total Payload: %d
Successes: %d
Redirects: %d
Client Failures: %d
Server Failures: %d
Count: %d

`

type TrafficMonitor interface {
	Monitor(Event)
	Stop()
}

type TrafficStatistics struct {
	Section            string
	AveragePayloadSize float64
	TotalPayloadSize   int64
	Successes          int
	Redirects          int
	ClientFailures     int
	ServerFailures     int
	Count              int64
}

func (t *TrafficStatistics) String() string {
	return fmt.Sprintf(
		SummaryStatisticsFormat,
		t.Section,
		t.AveragePayloadSize,
		t.TotalPayloadSize,
		t.Successes,
		t.Redirects,
		t.ClientFailures,
		t.ServerFailures,
		t.Count,
	)
}

type SummaryStatsTrafficMonitor struct {
	duration     time.Duration
	notification Notification

	ticker *time.Ticker
	events chan Event

	statistics map[string]*TrafficStatistics
}

func NewSummaryStatsTrafficMonitor(duration time.Duration, notification Notification) *SummaryStatsTrafficMonitor {
	monitor := &SummaryStatsTrafficMonitor{
		duration:     duration,
		notification: notification,
		events:       make(chan Event),
		statistics:   map[string]*TrafficStatistics{},
	}
	go monitor.consumeEvents()
	go monitor.publishStatistics()
	return monitor
}

func (s *SummaryStatsTrafficMonitor) Monitor(event Event) {
	s.events <- event
}

func (s *SummaryStatsTrafficMonitor) Stop() {
	s.ticker.Stop()
	close(s.events)
}

func (s *SummaryStatsTrafficMonitor) summary() string {
	statistics := make([]string, len(s.statistics))
	i := 0
	for _, statistic := range s.statistics {
		statistics[i] = statistic.String()
		i++
	}
	return strings.Join(statistics, "\n")
}

func (s *SummaryStatsTrafficMonitor) publishStatistics() {
	s.ticker = time.NewTicker(s.duration)
	for _ = range s.ticker.C {
		if s.summary() != "" {
			s.notification.Send(s.summary())
		}
		s.statistics = map[string]*TrafficStatistics{}
	}
}

func (s *SummaryStatsTrafficMonitor) consumeEvents() {
	for event := range s.events {
		eventURL, err := url.Parse(event.Path)
		if err != nil {
			continue
		}
		pathSections := strings.Split(eventURL.Path, "/")
		if len(pathSections) >= 2 {
			section := strings.Join(pathSections[0:2], "/")
			s.updateStatistics(section, event)
		}
		s.updateTotalStatistics(event)
	}
}

func (s *SummaryStatsTrafficMonitor) updateStatistics(section string, event Event) {
	if _, ok := s.statistics[section]; !ok {
		s.statistics[section] = &TrafficStatistics{Section: section}
	}

	sectionStatistics := s.statistics[section]
	sectionStatistics.AveragePayloadSize = (float64(sectionStatistics.Count)*sectionStatistics.AveragePayloadSize + float64(event.PayloadSize)) / float64(sectionStatistics.Count+1)
	sectionStatistics.TotalPayloadSize += int64(event.PayloadSize)
	if event.StatusCode >= 200 && event.StatusCode < 300 {
		sectionStatistics.Successes += 1
	} else if event.StatusCode >= 300 && event.StatusCode < 400 {
		sectionStatistics.Redirects += 1
	} else if event.StatusCode >= 400 && event.StatusCode < 500 {
		sectionStatistics.ClientFailures += 1
	} else if event.StatusCode >= 500 {
		sectionStatistics.ServerFailures += 1
	}

	sectionStatistics.Count += 1
}

func (s *SummaryStatsTrafficMonitor) updateTotalStatistics(event Event) {
	s.updateStatistics("Total Traffic", event)
}
