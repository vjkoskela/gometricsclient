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
	"time"
)

// Interface for logging metrics: timers, counters and gauges. Clients should
// create one instance of an implementing class for each unit of work. At the
// end of the unit of work the client should invoke close() on that instance.
// After the close() method is invoked the Metrics instance cannot be used to
// record further metrics and should be discarded.
type Metrics interface {

	// Create and initialize a counter sample. It is valid to create multiple Counter instances with the same name, even
	// concurrently, each will record a unique sample for the counter of the specified name.
	CreateCounter(string) Counter

	// Increment the specified counter by 1. All counters are initialized to zero. Creates a sample if one does not
	// exist. To create a new sample invoke resetCounter.
	IncrementCounter(string)

	// Increment the specified counter by the specified amount. All counters are initialized to zero. Creates a sample
	// if one does not exist. To create a new sample invoke resetCounter.
	IncrementCounterByValue(string, int64)

	// Decrement the specified counter by 1. All counters are initialized to zero. Creates a sample if one does not
	// exist. To create a new sample invoke resetCounter.
	DecrementCounter(string)

	// Decrement the specified counter by the specified amount. All counters are initialized to zero. Creates a sample
	// if one does not exist. To create a new sample invoke resetCounter.
	DecrementCounterByValue(string, int64)

	// Create a new sample for the counter with value zero. This most commonly used to record a zero-count for a
	// particular counter. If clients wish to record set count metrics then all counters should be reset before
	// conditionally invoking increment and/or decrement.
	ResetCounter(string)

	// Create and start a timer. It is valid to create multiple Timer instances with the same name, even concurrently,
	// each will record a unique sample for the timer of the specified name.
	CreateTimer(string) Timer

	// Start measurement of a sample for the specified timer. Use createTimer to make multiple concurrent measurements
	// for the same timer.
	StartTimer(string)

	// Stop measurement of the current sample for the specified timer. Use createTimer to make multiple concurrent
	// measurements of the same timer.
	StopTimer(string)

	// Set the timer to the specified value. This is most commonly used to record timers from external sources that are
	// not directly integrated with the metrics client.
	SetTimer(string, int64, Unit)

	// Set the specified gauge reading as a double.
	SetGaugeDouble(string, float64)

	// Set the specified gauge reading as a double with a well-known unit.
	SetGaugeDoubleWithUnit(string, float64, Unit)

	// Set the specified gauge reading as a long.
	SetGaugeLong(string, int64)

	// Set the specified gauge reading as a long with a well-known unit.
	SetGaugeLongWithUnit(string, int64, Unit)

	// Add an attribute that describes the captured metrics or context.
	AddAnnotation(string, string)

	// Add attributes that describe the captured metrics or context.
	AddAnnotations(map[string]string)

	// Is this Metrics instance open.
	IsOpen() bool

	// Close this Metrics instance. This should complete publication of metrics to the underlying data store(s). Once
	// the metrics object is closed, no further metrics can be recorded.
	Close()

	// Time this Metrics instance was opened. If this instance has not been opened the returned time is zero which may
	// be determined using the time.Time.IsZero function.
	GetOpenTime() time.Time

	// Time this Metrics instance was closed. If this instance has not been closed the returned time is zero which may
	// be determined using the time.Time.IsZero function.
	GetCloseTime() time.Time
}
