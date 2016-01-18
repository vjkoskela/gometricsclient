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
	"bytes"
)

var (
	_ Unit = (*tsdCompoundUnit)(nil)
	_ CompoundUnit = (*tsdCompoundUnit)(nil)
)

// CompoundUnit implementation.
type tsdCompoundUnit struct {
	numeratorUnits []Unit
	denominatorUnits []Unit
}

func newTsdCompoundUnit(options ...func(*tsdCompoundUnit)) Unit {
	var u *tsdCompoundUnit = new(tsdCompoundUnit)

	// Apply options
	for _, option := range options {
		option(u)
	}

	// Simplify and normalize
	var units map[string]Unit
	var numeratorUnits map[string]int
	var denominatorUnits map[string]int

	// 1. Flatten any nested compound units
	flattenUnits(u.numeratorUnits, units, numeratorUnits, denominatorUnits)
	flattenUnits(u.denominatorUnits, units, denominatorUnits, numeratorUnits)

	// 2. Remove any redundant units (e.g. in both the numerator and denominator)
	reduceCommonUnits(units, numeratorUnits, denominatorUnits)

	// 3. Flatten the unit counts
	u.numeratorUnits = nil
	for n, c := range numeratorUnits {
		nu := units[n]
		for i := 0; i < c; i++ {
			addUnit(u.numeratorUnits, nu)
		}
	}
	u.denominatorUnits = nil
	for n, c := range denominatorUnits {
		nu := units[n]
		for i := 0; i < c; i++ {
			addUnit(u.denominatorUnits, nu)
		}
	}

	// Return a no unit if possible
	if len(u.denominatorUnits) == 0 && len(u.numeratorUnits) == 0 {
		return nil;
	}

	// Return a base unit if possible
	if len(u.denominatorUnits) == 0 && len(u.numeratorUnits) == 1 {
		return u.numeratorUnits[0]
	}

	// Return the simplified compound unit
	return u
}

// Option to set the numerator units. Optional; default is no numerator units.
func tsdCompoundUnitNumeratorUnits(units ...Unit) func(*tsdCompoundUnit) {
	return func(u *tsdCompoundUnit) {
		for _, nu := range units {
			addUnit(u.numeratorUnits, nu)
		}
	}
}

// Option to add a numerator unit. Optional; default is no numerator units.
func tsdCompoundUnitAddNumeratorUnit(nu Unit) func(*tsdCompoundUnit) {
	return func(u *tsdCompoundUnit) {
		addUnit(u.numeratorUnits, nu)
	}
}

// Option to set the denominator units. Optional; default is no denominator units.
func tsdCompoundUnitDenominatorUnits(units ...Unit) func(*tsdCompoundUnit) {
	return func(u *tsdCompoundUnit) {
		for _, du := range units {
			addUnit(u.denominatorUnits, du)
		}
	}
}

// Option to add a denominator unit. Optional; default is no denominator units.
func tsdCompoundUnitAddDenominatorUnit(du Unit) func(*tsdCompoundUnit) {
	return func(u *tsdCompoundUnit) {
		addUnit(u.numeratorUnits, du)
	}
}

func (u *tsdCompoundUnit) Name() string {
	var nameBuffer *bytes.Buffer
	var numeratorParenthesis bool = len(u.numeratorUnits) > 1 && len(u.denominatorUnits) > 0
	var denominatorParenthesis bool = len(u.denominatorUnits) > 1
	if len(u.numeratorUnits) > 0 {
		if numeratorParenthesis {
			nameBuffer.WriteString("(")
		}
		for _, nu := range u.numeratorUnits {
			nameBuffer.WriteString(nu.Name())
			nameBuffer.WriteString("*")
		}
		nameBuffer.Truncate(nameBuffer.Len() - 1)
		if numeratorParenthesis {
			nameBuffer.WriteString(")")
		}
	} else {
		nameBuffer.WriteString("1")
	}
	if len(u.denominatorUnits) > 0 {
		if denominatorParenthesis {
			nameBuffer.WriteString("(")
		}
		for _, du := range u.denominatorUnits {
			nameBuffer.WriteString(du.Name())
			nameBuffer.WriteString("*")
		}
		nameBuffer.Truncate(nameBuffer.Len() - 1)
		if denominatorParenthesis {
			nameBuffer.WriteString(")")
		}
	}
	return nameBuffer.String()
}

func (u *tsdCompoundUnit) NumeratorUnits() []Unit {
	// TODO(vkoskela): This array should be immutable. [ISSUE-?]
	return u.numeratorUnits
}

func (u *tsdCompoundUnit) DenominatorUnits() []Unit {
	// TODO(vkoskela): This array should be immutable. [ISSUE-?]
	return u.denominatorUnits
}

func flattenUnits(units []Unit, names map[string]Unit, nUnits map[string]int, dUnits map[string]int) {
	for _, u := range units {
		switch ut := u.(type) {
		default:
			names[u.Name()] = u
			nUnits[u.Name()] = nUnits[u.Name()] + 1
		case CompoundUnit:
			flattenUnits(ut.NumeratorUnits(), names, nUnits, dUnits)
			flattenUnits(ut.DenominatorUnits(), names, dUnits, nUnits)
		}
	}
}

func reduceCommonUnits(names map[string]Unit, nUnits map[string]int, dUnits map[string]int) {
	for n, _ := range names {
		nuc := nUnits[n]
		duc := dUnits[n]
		if nuc > 0 && duc > 0 {
			if nuc > duc {
				delete(dUnits, n)
				nUnits[n] = nuc - duc
			} else if duc > nuc {
				delete(nUnits, n)
				dUnits[n] = duc - nuc
			} else { // nuc == duc
				delete(nUnits, n)
				delete(dUnits, n)
			}
		}
	}
}

func addUnit(a []Unit, u Unit) []Unit {
	n := len(a)
	if n == cap(a) {
		newA := make([]Unit, len(a), 2 * len(a) + 1)
		copy(newA, a)
		a = newA
	}
	a = a[0 : n + 1]
	a[n] = u
	return a
}
