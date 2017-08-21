// This extends the go flag package to allow for a "no default"
// variation of standard flag variables.  In order to accomplish this,
// we have to use pointers to pointers.  If the pp references a nil pointer,
// the flag was not set.  If the pp references a non-nil, the flag was set,
// and **ptr contains the value.
package nodefflag

import (
	"flag"
	"strconv"
)

// implement the Value interface for flags
type ndsf struct {
	sv *string
}

func (s *ndsf) String() string {
	if s.sv != nil {
		return *s.sv
	}
	return ""
}

func (s *ndsf) Set(val string) error {
	s.sv = &val
	return nil
}

type ndbf struct {
	bv *bool
}

func (b *ndbf) String() string {
	var ret bool
	if b.bv != nil {
		ret = *b.bv
	}
	return strconv.FormatBool(ret)
}

func (b *ndbf) Set(val string) error {
	pb, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	b.bv = &pb
	return nil
}

func (b ndbf) IsBoolFlag() bool {
	return true
}

type NoDefFlagSet struct {
	flag.FlagSet
}

func NewNoDefFlagSet(name string, errorHandling flag.ErrorHandling) *NoDefFlagSet {
	fs := flag.NewFlagSet(name, errorHandling)
	return &NoDefFlagSet{
		FlagSet: *fs,
	}
}

func (ndf *NoDefFlagSet) NoDefString(name, usage string) **string {
	var sv *string
	ndf.NoDefStringVar(&sv, name, usage)
	return &sv
}

func (ndf *NoDefFlagSet) NoDefStringVar(sv **string, name, usage string) {
	s := &ndsf{sv: *sv}
	ndf.Var(s, name, usage)
}

func (ndf *NoDefFlagSet) NoDefBoolVar(bv **bool, name, usage string) {
	b := &ndbf{bv: *bv}
	ndf.Var(b, name, usage)
}

func (ndf *NoDefFlagSet) NoDefBool(name string, usage string) **bool {
	b := &ndbf{}
	ndf.NoDefBoolVar(&b.bv, name, usage)
	return &b.bv
}
