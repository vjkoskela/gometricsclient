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
	_ Sink = (*StenoLogSink)(nil)
)

// Sink implementation which emits a json log line wrapped in steno per event to a specified file.
type StenoLogSink struct {
	directory string
	name string
	extension string
	maxHistory int
	logger logrus.StdLogger
	output *logrus.Logger
}

func NewStenoLogSink(options ...func(*StenoLogSink)) *StenoLogSink {
	var sls *StenoLogSink = new(StenoLogSink)

	// Defaults
	sls.logger = logrus.StandardLogger()
	sls.maxHistory = 24
	sls.extension = ".log"
	sls.name = "query"
	sls.directory = "./"

	// Apply options
	for _, option := range options {
		option(sls)
	}

	// Ensure directory ends with slash
	if !strings.HasSuffix(sls.directory, "/") {
		sls.directory = sls.directory + "/"
	}

	// Configure output
	// TODO(vkoskela): Lumberjack file rotation should support hourly and not just daily.
	// TODO(vkoskela): Metrics client should support and pass through size based rotation.
	var maxHistoryInDays int = int(math.Ceil(float64(24) / float64(sls.maxHistory)))
	if maxHistoryInDays < 1 {
		maxHistoryInDays = 1
	}
	sls.output = &logrus.Logger{
		Out: &lumberjack.Logger{
			Filename:   sls.directory + sls.name + sls.extension,
			MaxSize:    math.MaxInt32, // megabytes
			MaxBackups: sls.maxHistory,
			MaxAge:     maxHistoryInDays, // days
		},
		Formatter: gosteno.NewFormatter(),
		Level: logrus.InfoLevel,
	}

	return sls
}

// Option to set the directory. Optional; default is the current working directory of the application.
func StenoLogSinkDirectory(directory string) func(*StenoLogSink) {
	return func(sls *StenoLogSink) {
		sls.directory = directory
	}
}

// Option to set the file name without extension. Optional; default is "query". The file name without extension cannot
// be empty.
func StenoLogSinkName(name string) func(*StenoLogSink) {
	return func(sls *StenoLogSink) {
		sls.name = name
	}
}

// Option to set the file extension. Optional; default is ".log".
func StenoLogSinkExtension(extension string) func(*StenoLogSink) {
	return func(sls *StenoLogSink) {
		sls.extension = extension
	}
}

// Option to set the max history to retain. Optional; default is 24.
func StenoLogSinkMaxHistory(maxHistory int) func(*StenoLogSink) {
	return func(sls *StenoLogSink) {
		sls.maxHistory = maxHistory
	}
}

// Option to set the logger. Optional; default is the logrus.StdLogger instance.
func StenoLogSinkLogger(logger logrus.StdLogger) func(*StenoLogSink) {
	return func(sls *StenoLogSink) {
		sls.logger = logger
	}
}

func (ws *StenoLogSink) record(e Event) {
	// TODO
}
