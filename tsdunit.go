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
	_ Unit = (*tsdUnit)(nil)
)

// Unit implementation.
type tsdUnit struct {
	baseUnit baseUnit
	baseScale baseScale
}

func newTsdUnit(u baseUnit) *tsdUnit {
	return &tsdUnit{
		baseUnit: u,
	}
}

func newTsdUnitWithScale(u baseUnit, s baseScale) *tsdUnit {
	return &tsdUnit{
		baseUnit: u,
		baseScale: s,
	}
}

func (u *tsdUnit) Name() string {
	return u.baseScale.Name() + u.baseUnit.Name()
}

func (u *tsdUnit) BaseUnit() baseUnit {
	return u.baseUnit
}

func (u *tsdUnit) BaseScale() baseScale {
	return u.baseScale
}
