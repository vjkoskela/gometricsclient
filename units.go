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
	// ** Time **

	// Nanoseconds.
	Nanosecond Unit = newTsdUnitWithScale(second, nano)

	// Microseconds.
	Microsecond Unit = newTsdUnitWithScale(second, micro)

	// Milliseconds.
	Millisecond Unit = newTsdUnitWithScale(second, milli)

	// Seconds.
	Second Unit = newTsdUnit(second)

	// Minutes.
	Minute Unit = newTsdUnit(minute)

	// Hours.
	Hour Unit = newTsdUnit(hour)

	// Days.
	Day Unit = newTsdUnit(day)

	// ** Frequency **

	// TODO: Define all other supported units.
)