# redwood

An extensible http log collecting application that can monitor and send traffic alerts.

## Dependencies

The below dependencies are managed by the [Godeps](http://github.com/tools/godep) and will require godeps to be installed. Please see the [Godep installation](https://github.com/tools/godep#install) for more instructions.

- [ginkgo](https://github.com/onsi/ginkgo) for BDD style tests
- [gomega](github.com/onsi/gomega) for matchers used to create assertions in gingko

## Documentation

### Running the application

Running the application with the below command will require building it in [this section](#building).

```
sudo ./redwood
```

### Flags

Some flags that can be used to customize the application at runtime.

```
- file - File name of the file to monitor, collect, and/or alert on traffic logs
	- default: access.log
- monitor - Monitoring duration in seconds to which to send a summary
	- default: 10
- duration - Duration in seconds that
	- default: 30
- traffic - Traffic amount that should trigger an alert
	- default: 100
```

## Building

Run the below command in the

```
godep go build
```

## Running tests

Run the below command in the directory of the top most directory of the project.

```
godep go test ./...
```

## Design

Below are some of the extensible components, namely interfaces and what their responsibilities are. Under each component are a list of pre-existing components that implements the respective interface.

### Interfaces and Implementations

Components that can be extended or customized to be used in the application.

- `LogReader` reads logs and sends events through a channel
	- `FileLogReader` reads logs from file and sends parses the log into events to send through the channel
- `TrafficMonitor` monitors traffic
	- `SummaryStatsTrafficMonitor` generates statistical summaries for traffic received and sent
- `Alert` evaluates whether an event surpasses the threshold or reverts to normal
	- `TotalTrafficAlert` keeps track of the total number of events in a given time window
- `Notification` that determines when to alert
	- `ConsoleNotification` alerts to the console

## Domain Messages

Messages that are passed from one component to another.

- `Event` represents a network event within the http logs
- `TrafficStatistics` has fields for different traffic statistics such as average payload size and total payload size

## Application

Application that listens to network traffic and passes it through a filter, a monitor, a threshold, and eventually an alert if traffic surpasses the threshold.

- `Application` is composed of the different interfaces, namely the `LogReader`, `TrafficMonitor`, `Alert`, and `Notification` to allow custom components to read logs, monitor the filtered traffic, and alert when when the traffic surpasses some threshold

## License

redwood is released under the [MIT License](https://opensource.org/licenses/MIT).
