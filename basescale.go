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

// Base scale.
type baseScale int

const (
	yocto baseScale = iota + 1
	zepto
	atto
	femto
	pico
	nano
	micro
	milli
	centi
	deci
	deca
	hecto
	kilo
	mega
	giga
	tera
	peta
	exa
	zetta
	yotta
	kibi
	mebi
	gibi
	tebi
	pebi
	exbi
	zebi
	yobi
)

var baseScales = [...]string {
	"yocto",
	"zepto",
	"atto",
	"femto",
	"pico",
	"nano",
	"micro",
	"milli",
	"centi",
	"deci",
	"deca",
	"hecto",
	"kilo",
	"mega",
	"giga",
	"tera",
	"peta",
	"exa",
	"zetta",
	"yotta",
	"kibi",
	"mebi",
	"gibi",
	"tebi",
	"pebi",
	"exbi",
	"zebi",
	"yobi",
}

func (bs baseScale) Name() string {
	return baseScales[bs - 1]
}
