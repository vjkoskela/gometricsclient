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
	"github.com/Sirupsen/logrus"
)

var (
	_ Counter = (*tsdCounter)(nil)
	_ Quantity = (*tsdCounter)(nil)
)

// Counter implementation.
type tsdCounter struct {
	metrics Metrics
	name string
	value int64
	logger logrus.StdLogger
}

func newTsdCounter(m Metrics, n string, l logrus.StdLogger) *tsdCounter {
	if l == nil {
		l = logrus.StandardLogger()
	}
	var c *tsdCounter = &tsdCounter{
		metrics: m,
		name: n,
		logger: l}
	atomic.StoreInt64(&c.value, 0)
	return c
}

func (c *tsdCounter) Increment() {
	c.IncrementByValue(1)
}

func (c *tsdCounter) Decrement() {
	c.IncrementByValue(-1)
}

func (c *tsdCounter) IncrementByValue(d int64) {
	if !c.metrics.IsOpen() {
		c.logger.Printf("Counter manipulated after metrics instance closed; counter=%s", c.name);
	}
	atomic.AddInt64(&c.value, d)
}

func (c *tsdCounter) DecrementByValue(d int64) {
	c.IncrementByValue(-1 * d)
}

func (c *tsdCounter) Value() interface{} {
	return atomic.LoadInt64(&c.value)
}

func (c *tsdCounter) Unit() Unit {
	return nil
}
