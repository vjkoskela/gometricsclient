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

// Interface for counter. Instances are initialized to zero on creation. The
// zero-value sample is recorded when the Metrics instance is closed if no
// other actions are taken on this Counter instance.
//
// Modifying the Counter instance's value modifies the single sample value.
// When the associated Metrics instance is closed whatever value the sample
// has is recorded. To create another sample you create a new Counter instance
// with the same name.
//
// Each counter instance is bound to a Metrics instance. After the Metrics
// instance is closed any modifications to the Counter instances bound to that
// Metrics instance will be ignored.
type Counter interface {

	// Increment the counter sample by 1.
	Increment()

	// Decrement the counter sample by 1.
	Decrement()

	// Increment the counter sample by the specified value.
	IncrementByValue(int64)

	// Decrement the counter sample by the specified value.
	DecrementByValue(int64)
}
