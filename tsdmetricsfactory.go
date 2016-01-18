/*
Copyright 2016 Ville Koskela

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package gometricsclient

import (
	"os"
	"github.com/Sirupsen/logrus"
)

var (
	_ MetricsFactory = (*TsdMetricsFactory)(nil)
	defaultHost string
)

func init() {
	var value string
	var err error

	// Hostname
	value, err = os.Hostname()
	if (err == nil) {
		defaultHost = value
	}
}

/*
MetricsFactory implementation for creating Metrics instances to publish time series data (TSD).

For more information about the semantics of this class and its methods please refer to the MetricsFactory interface
documentation.

Create an instance of this class with the newDefaultTsdMetricsFactory method. This method will use default settings
where possible.

	var mf *gometricsclient.MetricsFactory = gometricsclient.NewDefaultTsdMetricsFactory(
		 "MyService",
		 "MyService-US-Prod",
		 "/usr/local/var/my-app/logs")

To customize the factory instance use the nested newTsdMetricsFactory method:

	var mf *gometricsclient.MetricsFactory = gometricsclient.NewTsdMetricsFactory(
		"MyService",
		"MyService-US-Prod",
		gometricsclient.TsdMetricsFactorySinks([]gometricsclient.Sink{gometricsclient.NewStenoLogSink()}))

The above will write metrics to the current working directory in query.log.  It is strongly recommended that at least a
path be set:

	var mf *gometricsclient.MetricsFactory = gometricsclient.NewTsdMetricsFactory(
		"MyService",
		"MyService-US-Prod",
		gometricsclient.TsdMetricsFactorySinks([]gometricsclient.Sink{gometricsclient.NewStenoLogSink(
			gometricsclient.StenoLogSinkDirectory("/usr/local/var/my-app/logs"))}))

The above will write metrics to /usr/local/var/my-app/logs in query.log. Additionally, you can customize the base file
name and extension for your application. However, if you are using TSDAggregator remember to configure it to match.

	var mf *gometricsclient.MetricsFactory = gometricsclient.NewTsdMetricsFactory(
		"MyService",
		"MyService-US-Prod",
		gometricsclient.TsdMetricsFactorySinks([]gometricsclient.Sink{gometricsclient.NewStenoLogSink(
			gometricsclient.StenoLogSinkDirectory("/usr/local/var/my-app/logs"),
			gometricsclient.StenoLogSinkName("tsd"),
			gometricsclient.StenoLogSinkExtension(".txt"))}))

The above will write metrics to /usr/local/var/my-app/logs in tsd.txt. The extension is configured separately as the
files are rolled over every hour inserting a date-time between the name and extension like:

query-log.YYYY-MM-DD-HH.log

This class is thread safe.
*/
type TsdMetricsFactory struct {
	sinks []Sink
	service string
	host string
	cluster string
	logger logrus.StdLogger
}

func NewDefaultTsdMetricsFactory(service string, cluster string, directory string) *TsdMetricsFactory {
	return NewTsdMetricsFactory(
		service,
		cluster,
		TsdMetricsFactorySinks(NewStenoLogSink(StenoLogSinkDirectory(directory))))
}

func NewTsdMetricsFactory(service string, cluster string, options ...func(*TsdMetricsFactory)) *TsdMetricsFactory {
	var mf *TsdMetricsFactory = new(TsdMetricsFactory)
	var failures []string = make([]string, 0, 10)

	// Defaults
	mf.logger = logrus.StandardLogger()
	mf.host = defaultHost
	mf.sinks = []Sink{NewStenoLogSink(StenoLogSinkLogger(mf.logger))}

	// Apply parameters
	mf.service = service
	mf.cluster = cluster

	// Apply options
	for _, option := range options {
		option(mf)
	}

	// Post-apply defaults
	// NOTE: These either rely on other parameters or are more expensive to create
	if mf.sinks == nil {
		mf.sinks = []Sink{NewStenoLogSink(StenoLogSinkLogger(mf.logger))}
	}

	// Validate
	if mf.service == "" {
		mf.service = "<SERVICE_NAME>"
		failures = appendFailure(failures, "Service cannot be empty")
	}
	if mf.cluster == "" {
		mf.cluster = "<CLUSTER_NAME>"
		failures = appendFailure(failures, "Cluster cannot be empty")
	}
	if mf.host == "" {
		mf.host = "<HOST_NAME>"
		failures = appendFailure(failures, "Host cannot be empty")
	}

	// Apply fallback
	if len(failures) > 0 {
		mf.logger.Printf("Unable to construct TsdMetricsFactory, metrics disabled; failures=%s", failures)
		mf.sinks = []Sink{newWarningSink(failures, mf.logger)}
	}

	return mf
}

// Option to set the sinks. Optional; default is an instance of StenoLogSink.
func TsdMetricsFactorySinks(sinks ...Sink) func(*TsdMetricsFactory) {
	return func(mf *TsdMetricsFactory) {
		mf.sinks = sinks
	}
}

// Option to set the host name. Optional; cannot be empty. Defaults to os.Hostname().
func TsdMetricsFactoryHost(host string) func(*TsdMetricsFactory) {
	return func(mf *TsdMetricsFactory) {
		mf.host = host
	}
}

// Option to set the logger. Optional; cannot be empty. Defaults to logrus.StdLogger().
func TsdMetricsFactoryLogger(logger logrus.StdLogger) func(*TsdMetricsFactory) {
	return func(mf *TsdMetricsFactory) {
		mf.logger = logger
	}
}

func (mf *TsdMetricsFactory) Create() Metrics {
	return newTsdMetrics(
		mf.sinks,
		mf.service,
		mf.cluster,
		mf.host,
		mf.logger)
}

func appendFailure(fs []string, f string) []string {
	n := len(fs)
	if n == cap(fs) {
		newFs := make([]string, len(fs), 2 * len(fs) + 1)
		copy(newFs, fs)
		fs = newFs
	}
	fs = fs[0 : n + 1]
	fs[n] = f
	return fs
}