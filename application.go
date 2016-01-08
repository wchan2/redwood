package main

type Application struct {
	logReader      LogReader
	trafficMonitor TrafficMonitor
	alert          Alert
}

func NewApplication(logReader LogReader, monitor TrafficMonitor, alert Alert) *Application {
	return &Application{
		logReader:      logReader,
		trafficMonitor: monitor,
		alert:          alert,
	}
}

func (a *Application) Run() {
	for event := range a.logReader.Read() {
		a.trafficMonitor.Monitor(event)
		a.alert.Check(event)
	}
}
