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

// Interface representing one unit of work's metrics data.
type Event interface {

	// Annotations for unit of work.
	Annotations() map[string]string

	// Timer samples by name for unit of work.
	TimerSamples() map[string][]Quantity

	// Counter samples by name for unit of work.
	CounterSamples() map[string][]Quantity

	// Gauge samples by name for unit of work.
	GaugeSamples() map[string][]Quantity
}
