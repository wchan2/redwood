package main

import (
	"fmt"
	"strings"
	"time"
)

type TrafficMonitor interface {
	Monitor(Event)
	Stop()
}

type TrafficStatistics struct {
	Section            string
	AveragePayloadSize float64
	TotalPayloadSize   int64
	Count              int64
}

func (t *TrafficStatistics) String() string {
	statisticsFormat := "Section: %s\nAverage Payload: %f\nTotal Payload: %d\nCount: %d\n"
	return fmt.Sprintf(statisticsFormat, t.Section, t.AveragePayloadSize, t.TotalPayloadSize, t.Count)
}

type SummaryStatsTrafficMonitor struct {
	duration     time.Duration
	notification Notification

	ticker *time.Ticker
	events chan Event

	statistics map[string]TrafficStatistics
}

func NewSummaryStatsTrafficMonitor(duration time.Duration, notification Notification) *SummaryStatsTrafficMonitor {
	monitor := &SummaryStatsTrafficMonitor{
		duration:     duration,
		notification: notification,
		events:       make(chan Event),

		statistics: map[string]TrafficStatistics{},
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
		s.statistics = map[string]TrafficStatistics{}
	}
}

func (s *SummaryStatsTrafficMonitor) consumeEvents() {
	for event := range s.events {
		if statistic, ok := s.statistics[event.Host]; ok {
			s.statistics[event.Host] = TrafficStatistics{
				Section:            event.Host,
				AveragePayloadSize: (statistic.AveragePayloadSize + float64(event.PayloadSize)) / float64((statistic.Count + 1)),
				TotalPayloadSize:   statistic.TotalPayloadSize + int64(event.PayloadSize),
				Count:              statistic.Count + 1,
			}
		} else {
			s.statistics[event.Host] = TrafficStatistics{
				Section:            event.Host,
				AveragePayloadSize: float64(event.PayloadSize),
				TotalPayloadSize:   int64(event.PayloadSize),
				Count:              1,
			}
		}
	}
}
