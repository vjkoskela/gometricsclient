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

var (
	_ Event = (*tsdEvent)(nil)
)

// Event implementation.
type tsdEvent struct {
	annotations map[string]string
	timerSamples map[string][]Quantity
	counterSamples map[string][]Quantity
	gaugeSamples map[string][]Quantity
}

func newTsdEvent(a map[string]string, ts map[string][]Quantity, cs map[string][]Quantity, gs map[string][]Quantity) *tsdEvent {
	return &tsdEvent{
		annotations: a,
		timerSamples: ts,
		counterSamples: cs,
		gaugeSamples: gs,
	}
}

func (e *tsdEvent) Annotations() map[string]string {
	// TODO(vkoskela): This map should be immutable. [ISSUE-?]
	return e.annotations
}

func (e *tsdEvent) TimerSamples() map[string][]Quantity {
	// TODO(vkoskela): This map should be immutable. [ISSUE-?]
	return e.timerSamples
}

func (e *tsdEvent) CounterSamples() map[string][]Quantity {
	// TODO(vkoskela): This map should be immutable. [ISSUE-?]
	return e.counterSamples
}

func (e *tsdEvent) GaugeSamples() map[string][]Quantity {
	// TODO(vkoskela): This map should be immutable. [ISSUE-?]
	return e.gaugeSamples
}
