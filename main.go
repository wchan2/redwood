package main

import (
	"flag"
	"log"
	"time"
)

var (
	file string

	monitor  int
	duration int
	traffic  int
)

func init() {
	flag.StringVar(&file, "file", "access.log", "File name of the file to monitor, collect, and/or alert on traffic logs")
	flag.IntVar(&monitor, "monitor", 10, "Monitoring duration in seconds to which to send a summary")
	flag.IntVar(&duration, "duration", 120, "Duration in seconds that")
	flag.IntVar(&traffic, "traffic", 1000, "Traffic amount that should trigger an alert")
	flag.Parse()
}

func main() {
	flag.Parse()
	totalTrafficAlert := NewTotalTrafficAlert(traffic, time.Duration(duration)*time.Second, ConsoleNotification)
	trafficMonitor := NewSummaryStatsTrafficMonitor(time.Duration(monitor)*time.Second, ConsoleNotification)
	fileLogReader, err := NewLogFileReader(file)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer fileLogReader.Close()

	app := NewApplication(fileLogReader, trafficMonitor, totalTrafficAlert)
	app.Run()
}
