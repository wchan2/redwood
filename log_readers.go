package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	clientMatcher        = regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`)
	timeMatcher          = regexp.MustCompile(`\[[a-zA-Z0-9\/\:\s\-]*\]`)
	userAgentMatcher     = regexp.MustCompile(`\"([^\"]*)\"$`)
	statusCodeAndPayload = regexp.MustCompile(`\d+ \d+`)

	httpRequestMatcher  = regexp.MustCompile(`\"([A-Z]+)[^"]*\"`)
	httpMethodMatcher   = regexp.MustCompile(`GET|POST|PUT|PATCH|POST|DELETE|HEAD|OPTIONS`)
	httpProtocolMatcher = regexp.MustCompile(`HTTP[\/|\s][0-9]\.[0-9]`)
	pathMatcher         = regexp.MustCompile(`\/[a-zA-Z\-\/\?\&\=\.0-9]*`)

	logDateFormatWithTimezone    = "2/Jan/2006:15:04:05 -0700"
	logDateFormatWithoutTimezone = "2/Jan/2006:15:04:05"
)

type LogReader interface {
	Read() <-chan Event
	Close()
}

type LogFileReader struct {
	logs chan Event
	file *os.File
}

func NewLogFileReader(filename string) (*LogFileReader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	LogFileReader := &LogFileReader{
		logs: make(chan Event),
		file: file,
	}
	go LogFileReader.consumeFromFile()
	return LogFileReader, nil
}

func (f *LogFileReader) Read() <-chan Event {
	return f.logs
}

func (f *LogFileReader) Close() {
	close(f.logs)
}

func (f *LogFileReader) consumeFromFile() {
	bufferedReader := bufio.NewReader(f.file)
	for true {
		if _, err := bufferedReader.Peek(1); err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		line, err := bufferedReader.ReadString('\n')
		if err != nil && err != io.EOF {
			continue
		}
		event, err := f.parseLog(line)
		if err != nil {
			continue
		}
		f.logs <- event
	}
}

func (f *LogFileReader) parseLog(line string) (Event, error) {
	var logDate time.Time
	if dateWithTimezone, err := time.Parse(logDateFormatWithTimezone, strings.Trim(timeMatcher.FindString(line), "[]")); err == nil {
		logDate = dateWithTimezone
	} else if dateWithoutTimezone, err := time.Parse(logDateFormatWithoutTimezone, strings.Trim(timeMatcher.FindString(line), "[]")); err == nil {
		logDate = dateWithoutTimezone
	} else {
		return Event{}, err
	}

	httpRequest := strings.Trim(httpRequestMatcher.FindString(line), `"`)
	httpMethod := httpMethodMatcher.FindString(httpRequest)
	httpProtocol := httpProtocolMatcher.FindString(httpRequest)
	path := pathMatcher.FindString(httpRequest)

	client := clientMatcher.FindString(line)
	statusCodeAndPayload := statusCodeAndPayload.FindString(line)

	httpResponseDetails := strings.Split(statusCodeAndPayload, " ")
	if len(httpResponseDetails) != 2 {
		return Event{}, errors.New("status code and payload not present")
	}
	statusCode, err := strconv.Atoi(httpResponseDetails[0])
	if err != nil {
		return Event{}, fmt.Errorf("could not parse status code: %s to integer", httpResponseDetails[0])
	}

	payloadSize, err := strconv.Atoi(httpResponseDetails[1])
	if err != nil {
		return Event{}, fmt.Errorf("could not parse payload size: %s to integer", httpResponseDetails[1])
	}

	return Event{
		Client:      client,
		Time:        logDate,
		Method:      httpMethod,
		Path:        path,
		Protocol:    httpProtocol,
		StatusCode:  statusCode,
		PayloadSize: payloadSize,
		UserAgent:   strings.Trim(userAgentMatcher.FindString(line), `"`),
	}, nil
}
