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
package main

import (
	"github.com/vjkoskela/gometricsclient"
)

// This file contains code samples from documentation contained in the project. Any updates to the code below should
// be reflected in the documentation and vice versa.

func sampleNewDefaultTsdMetricsFactory() {
	_ *gometricsclient.MetricsFactory = gometricsclient.NewDefaultTsdMetricsFactory(
		"MyService",
		"MyService-US-Prod",
		"/usr/local/var/my-app/logs")
}

func sampleNewTsdMetricsFactory() {
	_ *gometricsclient.MetricsFactory = gometricsclient.NewTsdMetricsFactory(
		"MyService",
		"MyService-US-Prod",
		gometricsclient.TsdMetricsFactorySinks([]gometricsclient.Sink{gometricsclient.NewStenoLogSink()}))
}

func sampleNewTsdMetricsFactoryWithStenoLogSinkDirectory() {
	_ *gometricsclient.MetricsFactory = gometricsclient.NewTsdMetricsFactory(
		"MyService",
		"MyService-US-Prod",
		gometricsclient.TsdMetricsFactorySinks([]gometricsclient.Sink{gometricsclient.NewStenoLogSink(
			gometricsclient.StenoLogSinkDirectory("/usr/local/var/my-app/logs"))}))
}

func sampleNewTsdMetricsFactoryWithStenoLogSinkCustomized() {
	_ *gometricsclient.MetricsFactory = gometricsclient.NewTsdMetricsFactory(
		"MyService",
		"MyService-US-Prod",
		gometricsclient.TsdMetricsFactorySinks([]gometricsclient.Sink{gometricsclient.NewStenoLogSink(
			gometricsclient.StenoLogSinkDirectory("/usr/local/var/my-app/logs"),
			gometricsclient.StenoLogSinkName("tsd"),
			gometricsclient.StenoLogSinkExtension(".txt"))}))
}
