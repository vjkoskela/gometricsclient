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
	"sync"
	"sync/atomic"
	"time"
	"github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"
)

const (
	finalTimestampKey string = "_end"
	initialTimestampKey string = "_start"
	idKey string = "_id"
	hostKey string = "_host"
	serviceKey string = "_service"
	clusterKey string = "_cluster"
)

var (
	_ Metrics = (*tsdMetrics)(nil)
)

/*
Metrics implementation that publishes metrics as time series data (TSD).

For more information about the semantics of this class and its methods please refer to the Metrics interface
documentation. To create an instance of this class use TsdMetricsFactory.

This class is thread safe; however, it makes no effort to order operations on the same data. For example, it is safe for
two threads to start and stop separate timers but if the threads start and stop the same timer than it is up to the
caller to ensure that start is called before stop.

The library does attempt to detect incorrect usage, for example modifying metrics after closing and starting but never
stopping a timer; however, in a multithreaded environment it is not guaranteed that these warnings are emitted. It is up
to clients to ensure that multithreaded use of the same TsdMetrics instance is correct.
*/
type tsdMetrics struct {
	sinks               []Sink
	service             string
	host                string
	cluster             string
	logger              logrus.StdLogger
	isOpen              int32
	initialTime         time.Time
	finalTme            time.Time
	annotations         map[string]string
	mutexAnnotations    sync.Mutex
	guageSamples        map[string][]Quantity
	mutexGuageSamples   sync.Mutex
	timerSamples        map[string][]Quantity
	timers              map[string]*tsdTimer
	mutexTimerSamples   sync.Mutex
	counterSamples      map[string][]Quantity
	counters            map[string]*tsdCounter
	mutexCounterSamples sync.Mutex
}

func newTsdMetrics(s []Sink, sn string, cn string, hn string, l logrus.StdLogger) *tsdMetrics {
	if l == nil {
		l = logrus.StandardLogger()
	}

	return &tsdMetrics{
		sinks: s,
		service: sn,
		host: hn,
		cluster: cn,
		logger: l,
		isOpen: 1,
		initialTime: time.Now(),
	}
}

func (m *tsdMetrics) CreateCounter(n string) Counter {
	if !m.assertIsOpen() {
		return nil
	}
	m.mutexCounterSamples.Lock()
	defer m.mutexCounterSamples.Unlock()
	samples, samplesExist := m.counterSamples[n]
	if !samplesExist {
		samples = make([]Quantity, 0, 10)
	}
	var counter *tsdCounter = newTsdCounter(m, n, m.logger)
	m.timerSamples[n] = addQuantity(samples, counter)
	return counter
}

func (m *tsdMetrics) IncrementCounter(n string) {
	m.IncrementCounterByValue(n, 1)
}

func (m *tsdMetrics) IncrementCounterByValue(n string, v int64) {
	if !m.assertIsOpen() {
		return
	}
	m.mutexCounterSamples.Lock()
	defer m.mutexCounterSamples.Unlock()
	var counter *tsdCounter
	if counter, counterExist := m.counters[n]; !counterExist {
		counter = newTsdCounter(m, n, m.logger)
		samples, samplesExist := m.counterSamples[n]
		if !samplesExist {
			samples = make([]Quantity, 0, 10)
		}
		m.counterSamples[n] = addQuantity(samples, counter)
		m.counters[n] = counter
	}
	counter.IncrementByValue(v)
}

func (m *tsdMetrics) DecrementCounter(n string) {
	m.IncrementCounterByValue(n, -1)
}

func (m *tsdMetrics) DecrementCounterByValue(n string, v int64) {
	m.IncrementCounterByValue(n, -1 * v)
}

func (m *tsdMetrics) ResetCounter(n string) {
	if !m.assertIsOpen() {
		return
	}
	m.mutexCounterSamples.Lock()
	defer m.mutexCounterSamples.Unlock()
	var counter *tsdCounter = newTsdCounter(m, n, m.logger)
	samples, samplesExist := m.counterSamples[n]
	if !samplesExist {
		samples = make([]Quantity, 0, 10)
	}
	m.counterSamples[n] = addQuantity(samples, counter)
	m.counters[n] = counter
}

func (m *tsdMetrics) CreateTimer(n string) Timer {
	if !m.assertIsOpen() {
		return nil
	}
	m.mutexTimerSamples.Lock()
	defer m.mutexTimerSamples.Unlock()
	samples, samplesExist := m.timerSamples[n]
	if !samplesExist {
		samples = make([]Quantity, 0, 10)
	}
	var timer *tsdTimer = newTsdTimer(m, n, m.logger)
	m.timerSamples[n] = addQuantity(samples, timer)
	return timer
}

func (m *tsdMetrics) StartTimer(n string) {
	if !m.assertIsOpen() {
		return
	}
	m.mutexTimerSamples.Lock()
	defer m.mutexTimerSamples.Unlock()
	samples, samplesExist := m.timerSamples[n]
	if !samplesExist {
		samples = make([]Quantity, 0, 10)
	}
	if _, timerExist := m.timers[n]; timerExist {
		// To measure multiple samples for the same timer within the same metrics instance concurrently use CreateTimer
		m.logger.Printf("Cannot start timer because timer already started; timerName=%s", n)
		return
	}
	var timer *tsdTimer = newTsdTimer(m, n, m.logger)
	m.timers[n] = timer
	m.timerSamples[n] = addQuantity(samples, timer)
}

func (m *tsdMetrics) StopTimer(n string) {
	if !m.assertIsOpen() {
		return
	}
	m.mutexTimerSamples.Lock()
	defer m.mutexTimerSamples.Unlock()
	timer, timerExist := m.timers[n]
	if !timerExist {
		m.logger.Printf("Cannot stop timer because timer was not started; timerName=%s", n)
		return
	}
	timer.Stop()
}

func (m *tsdMetrics) SetTimer(n string, v int64, u Unit) {
	if !m.assertIsOpen() {
		return
	}
	m.mutexTimerSamples.Lock()
	defer m.mutexTimerSamples.Unlock()
	samples, samplesExist := m.timerSamples[n]
	if !samplesExist {
		samples = make([]Quantity, 0, 10)
	}
	m.timerSamples[n] = addQuantity(samples, newTsdQuantityFromLong(v, u))
}

func (m *tsdMetrics) SetGaugeDouble(n string, v float64) {
	if !m.assertIsOpen() {
		return
	}
	m.mutexGuageSamples.Lock()
	defer m.mutexGuageSamples.Unlock()
	samples, exist := m.guageSamples[n]
	if !exist {
		samples = make([]Quantity, 0, 10)
	}
	m.guageSamples[n] = addQuantity(samples, newTsdQuantityFromDouble(v, nil))
}

func (m *tsdMetrics) SetGaugeDoubleWithUnit(n string, v float64, u Unit) {
	if !m.assertIsOpen() {
		return
	}
	m.mutexGuageSamples.Lock()
	defer m.mutexGuageSamples.Unlock()
	samples, exist := m.guageSamples[n]
	if !exist {
		samples = make([]Quantity, 0, 10)
	}
	m.guageSamples[n] = addQuantity(samples, newTsdQuantityFromDouble(v, u))
}

func (m *tsdMetrics) SetGaugeLong(n string, v int64) {
	if !m.assertIsOpen() {
		return
	}
	m.mutexGuageSamples.Lock()
	defer m.mutexGuageSamples.Unlock()
	samples, exist := m.guageSamples[n]
	if !exist {
		samples = make([]Quantity, 0, 10)
	}
	m.guageSamples[n] = addQuantity(samples, newTsdQuantityFromLong(v, nil))
}

func (m *tsdMetrics) SetGaugeLongWithUnit(n string, v int64, u Unit) {
	if !m.assertIsOpen() {
		return
	}
	m.mutexGuageSamples.Lock()
	defer m.mutexGuageSamples.Unlock()
	samples, exist := m.guageSamples[n]
	if !exist {
		samples = make([]Quantity, 0, 10)
	}
	m.guageSamples[n] = addQuantity(samples, newTsdQuantityFromLong(v, u))
}

func (m *tsdMetrics) AddAnnotation(k string, v string) {
	if !m.assertIsOpen() {
		return
	}
	m.mutexAnnotations.Lock()
	defer m.mutexAnnotations.Unlock()
	m.annotations[k] = v
}

func (m *tsdMetrics) AddAnnotations(a map[string]string) {
	if !m.assertIsOpen() {
		return
	}
	m.mutexAnnotations.Lock()
	defer m.mutexAnnotations.Unlock()
	for k, v := range a {
		m.annotations[k] = v
	}
}

func (m *tsdMetrics) IsOpen() bool {
	return atomic.LoadInt32(&m.isOpen) == 1
}

func (m *tsdMetrics) Close() {
	if atomic.SwapInt32(&m.isOpen, 0) != 1 {
		m.logger.Printf("Metrics object was already closed");
		return
	}
	m.mutexAnnotations.Lock()
	defer m.mutexAnnotations.Unlock()
	m.mutexTimerSamples.Lock()
	defer m.mutexTimerSamples.Unlock()
	m.mutexCounterSamples.Lock()
	defer m.mutexCounterSamples.Unlock()
	m.mutexGuageSamples.Lock()
	defer m.mutexGuageSamples.Unlock()

	m.finalTme = time.Now()

	// TODO(vkoskela): Check for user entries on these reserved keys before overwriting. [ISSUE-?]
	m.annotations[idKey] = uuid.New()
	m.annotations[hostKey] = m.host
	m.annotations[serviceKey] = m.service
	m.annotations[clusterKey] = m.cluster
	m.annotations[initialTimestampKey] = m.initialTime.Format(time.RFC3339Nano)
	m.annotations[finalTimestampKey] = m.finalTme.Format(time.RFC3339Nano)

	var event Event = newTsdEvent(m.annotations, m.timerSamples, m.counterSamples, m.guageSamples)
	for _, s := range m.sinks {
		s.record(event)
	}
}

func (m *tsdMetrics) GetOpenTime() time.Time {
	return m.initialTime
}

func (m *tsdMetrics) GetCloseTime() time.Time {
	return m.finalTme
}

func (m *tsdMetrics) assertIsOpen() bool {
	var isOpen bool = m.IsOpen()
	if !isOpen {
		m.logger.Printf("Metrics object was closed during an operation; you may have a race condition");
	}
	return isOpen;
}

func addQuantity(a []Quantity, v Quantity) []Quantity {
	n := len(a)
	if n == cap(a) {
		newA := make([]Quantity, len(a), 2 * len(a) + 1)
		copy(newA, a)
		a = newA
	}
	a = a[0 : n + 1]
	a[n] = v
	return a
}