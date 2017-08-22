// Package nodefflag extends the go flag package to allow for a "no default"
// variation of standard flag variables.  In order to accomplish this,
// we have to use pointers to pointers.  If the pp references a nil pointer,
// the flag was not set.  If the pp references a non-nil pointer, the flag
// was set, and **ptr contains the value.  The pp itself returned will never
// be nil, and it is expected that the ND*Var methods will never receive
// a nil **.
package nodefflag

import (
	"flag"
	"strconv"
	"time"
)

// implement the Value interface for flags
type ndsf struct{ sv **string }

func (s ndsf) String() string {
	if *s.sv != nil {
		return **s.sv
	}
	return ""
}

func (s ndsf) Set(val string) error {
	*s.sv = &val
	return nil
}

func (s ndsf) Get() interface{} {
	return *s.sv
}

type ndbf struct{ bv **bool }

func (b ndbf) String() string {
	if *b.bv != nil {
		return strconv.FormatBool(**b.bv)
	}
	return ""
}

func (b ndbf) Set(val string) error {
	pb, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	*b.bv = &pb
	return nil
}

func (b ndbf) Get() interface{} {
	return *b.bv
}

func (b ndbf) IsBoolFlag() bool {
	return true
}

type ndif struct{ iv **int }

func (i ndif) String() string {
	if *i.iv != nil {
		return strconv.Itoa(**i.iv)
	}
	return ""
}

func (i ndif) Set(val string) error {
	pi, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	*i.iv = &pi
	return nil
}

func (i ndif) Get() interface{} {
	return *i.iv
}

type ndi64f struct{ iv **int64 }

func (i ndi64f) String() string {
	if *i.iv != nil {
		return strconv.FormatInt(**i.iv, 10)
	}
	return ""
}

func (i ndi64f) Set(val string) error {
	pi, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return err
	}
	*i.iv = &pi
	return nil
}

func (i ndi64f) Get() interface{} {
	return *i.iv
}

type nduif struct{ uiv **uint }

func (ui nduif) String() string {
	if *ui.uiv != nil {
		return strconv.FormatUint(uint64(**ui.uiv), 10)
	}
	return ""
}

func (ui nduif) Set(val string) error {
	pui, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return err
	}
	pui2 := uint(pui)
	*ui.uiv = &pui2
	return nil
}

func (ui nduif) Get() interface{} {
	return *ui.uiv
}

type ndui64f struct{ uiv **uint64 }

func (ui ndui64f) String() string {
	if *ui.uiv != nil {
		return strconv.FormatUint(**ui.uiv, 10)
	}
	return ""
}

func (ui ndui64f) Set(val string) error {
	pui, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return err
	}
	*ui.uiv = &pui
	return nil
}

func (ui ndui64f) Get() interface{} { return *ui.uiv }

type ndff struct{ fv **float64 }

func (f ndff) String() string {
	if *f.fv != nil {
		return strconv.FormatFloat(**f.fv, 'g', -1, 64)
	}
	return ""
}

func (f ndff) Set(val string) error {
	pf, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	*f.fv = &pf
	return nil
}

func (f ndff) Get() interface{} {
	return *f.fv
}

type nddf struct{ dv **time.Duration }

func (d nddf) String() string {
	if *d.dv != nil {
		return (*d.dv).String()
	}
	return ""
}

func (d nddf) Set(val string) error {
	pd, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	*d.dv = &pd
	return nil
}

func (d nddf) Get() interface{} {
	return *d.dv
}

// NDFlagSet - extends the flag package to add "no default" variants,
// where no defaults are specified.
type NDFlagSet struct{ flag.FlagSet }

// NewNDFlagSet - factory method, initializes the underlying FlagSet
func NewNDFlagSet(name string, errorHandling flag.ErrorHandling) *NDFlagSet {
	fs := flag.NewFlagSet(name, errorHandling)
	return &NDFlagSet{
		FlagSet: *fs,
	}
}

// NDString - returns double string pointer, will reference nil
// string pointer if flag was not set, will reference non-nil otherwise.
// This allows you to differentiate between the zero val ("") and not set.
func (ndf *NDFlagSet) NDString(name, usage string) **string {
	var sv *string
	ndf.NDStringVar(&sv, name, usage)
	return &sv
}

// NDStringVar - Similar to NDString, but you supply the double
// string pointer.
func (ndf *NDFlagSet) NDStringVar(sv **string, name, usage string) {
	s := ndsf{sv: sv}
	ndf.Var(s, name, usage)
}

// NDBool - returns double bool pointer, will reference
// nil bool pointer if flag was not set, will reference non-nil otherwise.
// This allows you to differentiate between the zero val (false) and not set.
func (ndf *NDFlagSet) NDBool(name string, usage string) **bool {
	var bv *bool
	ndf.NDBoolVar(&bv, name, usage)
	return &bv
}

// NDBoolVar - similar to NDBool, but you supply the double
// bool pointer.
func (ndf *NDFlagSet) NDBoolVar(bv **bool, name, usage string) {
	b := ndbf{bv: bv}
	ndf.Var(b, name, usage)
}

// NDInt - returns an int double pointers, will reference
// nil int pointer if flag was not set, will reference non-nil otherwise.
// This allows you to differentiate between the zero val (0) and not set.
func (ndf *NDFlagSet) NDInt(name, usage string) **int {
	var iv *int
	ndf.NDIntVar(&iv, name, usage)
	return &iv
}

// NDIntVar - similar to NDInt, but you sply the double pointer.
func (ndf *NDFlagSet) NDIntVar(iv **int, name, usage string) {
	i := ndif{iv: iv}
	ndf.Var(i, name, usage)
}

// NDInt64 - NDInt but type is int64
func (ndf *NDFlagSet) NDInt64(name, usage string) **int64 {
	var iv *int64
	ndf.NDInt64Var(&iv, name, usage)
	return &iv
}

// NDInt64Var - NDIntVar but for int64
func (ndf *NDFlagSet) NDInt64Var(iv **int64, name, usage string) {
	i := ndi64f{iv: iv}
	ndf.Var(i, name, usage)
}

// NDUint - returns double pointer to a uint.
func (ndf *NDFlagSet) NDUint(name, usage string) **uint {
	var uiv *uint
	ndf.NDUintVar(&uiv, name, usage)
	return &uiv
}

// NDUintVar - same as NDUint, but you supply the double p.
func (ndf *NDFlagSet) NDUintVar(uiv **uint, name, usage string) {
	ui := nduif{uiv: uiv}
	ndf.Var(ui, name, usage)
}

// NDUint64 - uint64 version of NDUint
func (ndf *NDFlagSet) NDUint64(name, usage string) **uint64 {
	var uiv *uint64
	ndf.NDUint64Var(&uiv, name, usage)
	return &uiv
}

// NDUint64Var - uint64 version of NDUintVar
func (ndf *NDFlagSet) NDUint64Var(uiv **uint64, name, usage string) {
	ui := ndui64f{uiv: uiv}
	ndf.Var(ui, name, usage)
}

// NDFloat64 - returns double pointer to a float64.  Works the same
// as all the other numeric types.
func (ndf *NDFlagSet) NDFloat64(name, usage string) **float64 {
	var fv *float64
	ndf.NDFloat64Var(&fv, name, usage)
	return &fv
}

// NDFloat64Var - you supply the pointer, but same as NDFloat64
func (ndf *NDFlagSet) NDFloat64Var(fv **float64, name, usage string) {
	f := ndff{fv: fv}
	ndf.Var(f, name, usage)
}

// NDDuration - duration flag.  returns double pointer, if references
// nil the flag was not set, otherwise it was set.
func (ndf *NDFlagSet) NDDuration(name, usage string) **time.Duration {
	var dv *time.Duration
	ndf.NDDurationVar(&dv, name, usage)
	return &dv
}

// NDDurationVar - BYO duration pp version of NDDuration
func (ndf *NDFlagSet) NDDurationVar(dv **time.Duration, name, usage string) {
	d := nddf{dv: dv}
	ndf.Var(d, name, usage)
}
