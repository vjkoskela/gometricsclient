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
	_ Quantity = (*tsdQuantity)(nil)
)

// Quantity implementation.
type tsdQuantity struct {
	value interface{}
	unit Unit
}

func newTsdQuantityFromLong(v int64, u Unit) *tsdQuantity {
	return &tsdQuantity{
		value: v,
		unit: u,
	}
}

func newTsdQuantityFromDouble(v float64, u Unit) *tsdQuantity {
	return &tsdQuantity{
		value: v,
		unit: u,
	}
}

func (q *tsdQuantity) Value() interface{} {
	return q.value
}

func (q *tsdQuantity) Unit() Unit {
	return q.unit
}
