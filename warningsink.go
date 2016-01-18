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
	"github.com/Sirupsen/logrus"
)

var (
	_ Sink = (*warningSink)(nil)
)

// Sink implementation which emits a warning each time an Event is to recorded.
type warningSink struct {
	reasons []string
	logger logrus.StdLogger
}

func newWarningSink(r []string, l logrus.StdLogger) *warningSink {
	if (l == nil) {
		l = logrus.StandardLogger()
	}
	return &warningSink{reasons: r, logger: l}
}

func (ws *warningSink) record(e Event) {
	ws.logger.Printf("Unable to record event, metrics disabled; reasons=%v", ws.reasons)
}
