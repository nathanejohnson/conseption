// This extends the go flag package to allow for a "no default"
// variation of standard flag variables.  In order to accomplish this,
// we have to use pointers to pointers.  If the pp references a nil pointer,
// the flag was not set.  If the pp references a non-nil pointer, the flag
// was set, and **ptr contains the value.  The pp itself returned will never
// be nil, and it is expected that the NoDef*Var methods will never receive
// a nil **.
package nodefflag

import (
	"flag"
	"strconv"
	"time"
)

// implement the Value interface for flags
type ndsf struct {
	sv **string
}

func (s *ndsf) String() string {
	if *s.sv != nil {
		return **s.sv
	}
	return ""
}

func (s *ndsf) Set(val string) error {
	*s.sv = &val
	return nil
}

type ndbf struct {
	bv **bool
}

func (b *ndbf) String() string {
	var ret bool
	if *b.bv != nil {
		ret = **b.bv
	}
	return strconv.FormatBool(ret)
}

func (b *ndbf) Set(val string) error {
	pb, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	*b.bv = &pb
	return nil
}

func (b *ndbf) IsBoolFlag() bool {
	return true
}

type ndif struct {
	iv **int
}

func (i *ndif) String() string {
	if *i.iv != nil {
		return strconv.Itoa(**i.iv)
	}
	return ""
}

func (i *ndif) Set(val string) error {
	pi, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	*i.iv = &pi
	return nil
}

type ndi64f struct {
	iv **int64
}

func (i *ndi64f) String() string {
	if *i.iv != nil {
		return strconv.FormatInt(**i.iv, 10)
	}
	return ""
}

func (i *ndi64f) Set(val string) error {
	pi, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return err
	}
	*i.iv = &pi
	return nil
}

type nduif struct {
	uiv **uint
}

func (ui *nduif) String() string {
	if *ui.uiv != nil {
		return strconv.FormatUint(uint64(**ui.uiv), 10)
	}
	return ""
}

func (ui *nduif) Set(val string) error {
	pui, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return err
	}
	pui2 := uint(pui)
	*ui.uiv = &pui2
	return nil
}

type ndui64f struct {
	uiv **uint64
}

func (ui *ndui64f) String() string {
	if *ui.uiv != nil {
		return strconv.FormatUint(**ui.uiv, 10)
	}
	return ""
}

func (ui *ndui64f) Set(val string) error {
	pui, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return err
	}
	*ui.uiv = &pui
	return nil
}

type ndff struct {
	fv **float64
}

func (f *ndff) String() string {
	if *f.fv != nil {
		return strconv.FormatFloat(**f.fv, 'g', -1, 64)
	}
	return ""
}

func (f *ndff) Set(val string) error {
	pf, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	*f.fv = &pf
	return nil
}

type nddf struct {
	dv **time.Duration
}

func (d *nddf) String() string {
	if *d.dv != nil {
		return (*d.dv).String()
	}
	return ""
}

func (d *nddf) Set(val string) error {
	pd, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	*d.dv = &pd
	return nil
}

// NoDefFlagSet - extends the flag package to add "no default" variants,
// where no defaults are specified.
type NoDefFlagSet struct {
	flag.FlagSet
}

// NewNoDefFlagSet - factory method, initializes the underlying FlagSet
func NewNoDefFlagSet(name string, errorHandling flag.ErrorHandling) *NoDefFlagSet {
	fs := flag.NewFlagSet(name, errorHandling)
	return &NoDefFlagSet{
		FlagSet: *fs,
	}
}

// NoDefString - returns double string pointer, will reference nil
// string pointer if flag was not set, will reference non-nil otherwise.
// This allows you to differentiate between the zero val ("") and not set.
func (ndf *NoDefFlagSet) NoDefString(name, usage string) **string {
	var sv *string
	ndf.NoDefStringVar(&sv, name, usage)
	return &sv
}

// NoDefStringVar - Similar to NoDefString, but you supply the double
// string pointer.
func (ndf *NoDefFlagSet) NoDefStringVar(sv **string, name, usage string) {
	s := &ndsf{sv: sv}
	ndf.Var(s, name, usage)
}

// NoDefBool - returns double bool pointer, will reference
// nil bool pointer if flag was not set, will reference non-nil otherwise.
// This allows you to differentiate between the zero val (false) and not set.
func (ndf *NoDefFlagSet) NoDefBool(name string, usage string) **bool {
	var bv *bool
	ndf.NoDefBoolVar(&bv, name, usage)
	return &bv
}

// NoDefBoolVar - similar to NoDefBool, but you supply the double
// bool pointer.
func (ndf *NoDefFlagSet) NoDefBoolVar(bv **bool, name, usage string) {
	b := &ndbf{bv: bv}
	ndf.Var(b, name, usage)
}

// NoDefInt - returns an int double pointers, will reference
// nil int pointer if flag was not set, will reference non-nil otherwise.
// This allows you to differentiate between the zero val (0) and not set.
func (ndf *NoDefFlagSet) NoDefInt(name, usage string) **int {
	var iv *int
	ndf.NoDefIntVar(&iv, name, usage)
	return &iv
}

// NoDefIntVar - similar to NoDefInt, but you sply the double pointer.
func (ndf *NoDefFlagSet) NoDefIntVar(iv **int, name, usage string) {
	i := &ndif{iv: iv}
	ndf.Var(i, name, usage)
}

// NoDefInt64 - NoDefInt but type is int64
func (ndf *NoDefFlagSet) NoDefInt64(name, usage string) **int64 {
	var iv *int64
	ndf.NoDefInt64Var(&iv, name, usage)
	return &iv
}

// NoDefInt64Var - NoDefIntVar but for int64
func (ndf *NoDefFlagSet) NoDefInt64Var(iv **int64, name, usage string) {
	i := &ndi64f{iv: iv}
	ndf.Var(i, name, usage)
}

// NoDefUint - returns double pointer to a uint.
func (ndf *NoDefFlagSet) NoDefUint(name, usage string) **uint {
	var uiv *uint
	ndf.NoDefUintVar(&uiv, name, usage)
	return &uiv
}

// NoDefUintVar - same as NoDefUint, but you supply the double p.
func (ndf *NoDefFlagSet) NoDefUintVar(uiv **uint, name, usage string) {
	ui := &nduif{uiv: uiv}
	ndf.Var(ui, name, usage)
}

// NoDefUint64 - uint64 version of NoDefUint
func (ndf *NoDefFlagSet) NoDefUint64(name, usage string) **uint64 {
	var uiv *uint64
	ndf.NoDefUint64Var(&uiv, name, usage)
	return &uiv
}

// NoDefUnit64Var - uint64 version of NoDefUintVar
func (ndf *NoDefFlagSet) NoDefUint64Var(uiv **uint64, name, usage string) {
	ui := &ndui64f{uiv: uiv}
	ndf.Var(ui, name, usage)
}

// NoDefFloat64 - returns double pointer to a float64.  Works the same
// as all the other numeric types.
func (ndf *NoDefFlagSet) NoDefFloat64(name, usage string) **float64 {
	var fv *float64
	ndf.NoDefFloat64Var(&fv, name, usage)
	return &fv
}

// NoDefFloat64Var - you supply the pointer, but same as NoDefFloat64
func (ndf *NoDefFlagSet) NoDefFloat64Var(fv **float64, name, usage string) {
	f := &ndff{fv: fv}
	ndf.Var(f, name, usage)
}

// NoDefDuration - duration flag.  returns double pointer, if references
// nil the flag was not set, otherwise it was set.
func (ndf *NoDefFlagSet) NoDefDuration(name, usage string) **time.Duration {
	var dv *time.Duration
	ndf.NoDefDurationVar(&dv, name, usage)
	return &dv
}

// NoDefDurationVar - BYO duration pp version of NoDefDuration
func (ndf *NoDefFlagSet) NoDefDurationVar(dv **time.Duration, name, usage string) {
	d := &nddf{dv: dv}
	ndf.Var(d, name, usage)
}
