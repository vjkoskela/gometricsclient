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
	"math"
	"strings"
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/Sirupsen/logrus"
	"github.com/vjkoskela/gosteno"
)

var (
	_ Sink = (*TsdLogSink)(nil)
)

// Sink implementation which emits a json log line per event to a specified file.
type TsdLogSink struct {
	directory string
	name string
	extension string
	maxHistory int
	logger logrus.StdLogger
	output *logrus.Logger
}

func NewTsdLogSink(options ...func(*TsdLogSink)) *TsdLogSink {
	var tls *TsdLogSink = new(TsdLogSink)

	// Defaults
	tls.logger = logrus.StandardLogger()
	tls.maxHistory = 24
	tls.extension = ".log"
	tls.name = "query"
	tls.directory = "./"

	// Apply options
	for _, option := range options {
		option(tls)
	}

	// Ensure directory ends with slash
	if !strings.HasSuffix(tls.directory, "/") {
		tls.directory = tls.directory + "/"
	}

	// Configure output
	// TODO(vkoskela): Lumberjack file rotation should support hourly and not just daily.
	// TODO(vkoskela): Metrics client should support and pass through size based rotation.
	var maxHistoryInDays int = int(math.Ceil(float64(24) / float64(tls.maxHistory)))
	if maxHistoryInDays < 1 {
		maxHistoryInDays = 1
	}
	tls.output = &logrus.Logger{
		Out: &lumberjack.Logger{
			Filename:   tls.directory + tls.name + tls.extension,
			MaxSize:    math.MaxInt32, // megabytes
			MaxBackups: tls.maxHistory,
			MaxAge:     maxHistoryInDays, // days
		},
		Formatter: gosteno.NewFormatter(),
		Level: logrus.InfoLevel,
	}

	return tls
}

// Option to set the directory. Optional; default is the current working directory of the application.
func TsdLogSinkDirectory(directory string) func(*TsdLogSink) {
	return func(sls *TsdLogSink) {
		sls.directory = directory
	}
}

// Option to set the file name without extension. Optional; default is "query". The file name without extension cannot
// be empty.
func TsdLogSinkName(name string) func(*TsdLogSink) {
	return func(sls *TsdLogSink) {
		sls.name = name
	}
}

// Option to set the file extension. Optional; default is ".log".
func TsdLogSinkExtension(extension string) func(*TsdLogSink) {
	return func(sls *TsdLogSink) {
		sls.extension = extension
	}
}

// Option to set the max history to retain. Optional; default is 24.
func TsdLogSinkMaxHistory(maxHistory int) func(*TsdLogSink) {
	return func(sls *TsdLogSink) {
		sls.maxHistory = maxHistory
	}
}

// Option to set the logger. Optional; default is the logrus.StdLogger instance.
func TsdLogSinkLogger(logger logrus.StdLogger) func(*TsdLogSink) {
	return func(sls *TsdLogSink) {
		sls.logger = logger
	}
}

func (ws *TsdLogSink) record(e Event) {
	// TODO
}
