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
	"sync/atomic"
	"time"
	"github.com/Sirupsen/logrus"
)

var (
	_ Timer = (*tsdTimer)(nil)
	_ Quantity = (*tsdTimer)(nil)
)

// Timer implementation.
type tsdTimer struct {
	metrics Metrics
	name string
	isRunning int32
	isAborted int32
	startTime int64
	elapsedTimed int64
	logger logrus.StdLogger
}

func newTsdTimer(m Metrics, n string, l logrus.StdLogger) *tsdTimer {
	if l == nil {
		l = logrus.StandardLogger()
	}
	var t *tsdTimer = &tsdTimer{
		metrics: m,
		name: n,
		logger: l}
	atomic.StoreInt32(&t.isRunning, 1)
	atomic.StoreInt32(&t.isAborted, 0)
	atomic.StoreInt64(&t.startTime, time.Now().UnixNano())
	return t
}

func (t *tsdTimer) Stop() {
	var wasAborted int32 = atomic.LoadInt32(&t.isAborted)
	var wasRunning int32 = atomic.SwapInt32(&t.isRunning, 0)
	if !t.metrics.IsOpen() {
		t.logger.Printf("Timer stopped after metrics instance closed; timer=%s", t.name)
	}
	if wasAborted == 1 {
		t.logger.Printf("Timer stopped after aborted; timer=%s", t.name)
	} else if wasRunning == 0 {
		t.logger.Printf("Timer stopped multiple times; timer=%s", t.name);
	} else {
		atomic.StoreInt64(&t.elapsedTimed, time.Now().UnixNano() - atomic.LoadInt64(&t.startTime))
	}
}

func (t *tsdTimer) Abort() {
	var wasAborted int32 = atomic.SwapInt32(&t.isAborted, 1)
	var wasRunning int32 = atomic.SwapInt32(&t.isRunning, 0)
	if !t.metrics.IsOpen() {
		t.logger.Printf("Timer aborted after metrics instance closed; timer=%s", t.name)
	}
	if wasAborted == 1 {
		t.logger.Printf("Timer aborted multiple times; timer=%s", t.name)
	} else if wasRunning == 0 {
		t.logger.Printf("Timer aborted after closed/stopped; timer=%s", t.name)
	}

}

func (t *tsdTimer) Value() interface{} {
	if atomic.LoadInt32(&t.isRunning) == 0 {
		t.logger.Printf("Timer access before it is stopped; timer=%s", t.name)
	}
	return atomic.LoadInt64(&t.elapsedTimed)
}

func (ws *tsdTimer) Unit() Unit {
	return Nanosecond
}

func (t *tsdTimer) Running() bool {
	return atomic.LoadInt32(&t.isRunning) == 1
}

func (t *tsdTimer) Aborted() bool {
	return atomic.LoadInt32(&t.isAborted) == 1
}
