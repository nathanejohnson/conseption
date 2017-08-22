package nodefflag

import (
	"flag"
	"testing"
	"time"
)

// These tests were heavily inspired by / ripped off of the
// "flag" package test suite.

func TestNils(t *testing.T) {
	fs := nfs()
	m := make(map[string]*flag.Flag)
	visitor := func(f *flag.Flag) {
		m[f.Name] = f
		g, ok := f.Value.(flag.Getter)
		if !ok {
			t.Errorf("Visit: value does not satisfy Getter: %T", f.Value)
		}
		switch f.Name {
		case "test_bool":
			ok = g.Get() == (*bool)(nil)
		case "test_int":
			ok = g.Get() == (*int)(nil)
		case "test_int64":
			ok = g.Get() == (*int64)(nil)
		case "test_uint":
			ok = g.Get() == (*uint)(nil)
		case "test_uint64":
			ok = g.Get() == (*uint64)(nil)
		case "test_string":
			ok = g.Get() == (*string)(nil)
		case "test_float64":
			ok = g.Get() == (*float64)(nil)
		case "test_duration":
			ok = g.Get() == (*time.Duration)(nil)
		}
		if !ok {
			t.Errorf("Visit: bad value %T(%v) for %s", g.Get(), g.Get(), f.Name)
		}
	}
	fs.VisitAll(visitor)
	if len(m) != 8 {
		t.Error("VisitAll misses some flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
}
func TestSet(t *testing.T) {
	fs := nfs()

	_ = fs.Set("test_bool", "false")
	_ = fs.Set("test_int", "42")
	_ = fs.Set("test_int64", "-420")
	_ = fs.Set("test_uint", "80")
	_ = fs.Set("test_uint64", "800")
	_ = fs.Set("test_string", "your ad here")
	_ = fs.Set("test_float64", "123.45")
	_ = fs.Set("test_duration", "30s")

	m := make(map[string]*flag.Flag)
	visitor := func(f *flag.Flag) {
		m[f.Name] = f
		var ok bool
		g := f.Value.(flag.Getter)
		switch f.Name {
		case "test_bool":
			ok = *(g.Get().(*bool)) == false
		case "test_int":
			ok = *(g.Get().(*int)) == 42
		case "test_int64":
			ok = *(g.Get().(*int64)) == int64(-420)
		case "test_uint":
			ok = *(g.Get().(*uint)) == uint(80)
		case "test_uint64":
			ok = *(g.Get().(*uint64)) == uint64(800)
		case "test_string":
			ok = *(g.Get().(*string)) == "your ad here"
		case "test_float64":
			ok = *(g.Get().(*float64)) == float64(123.45)
		case "test_duration":
			ok = g.Get().(*time.Duration).String() == "30s"
		}
		if !ok {
			t.Errorf("Visit: bad value %T(%v) for %s", g.Get(), g.Get(), f.Name)
		}
	}
	fs.Visit(visitor)
	if len(m) != 8 {
		t.Error("Visit misses some flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}

}

func TestPtrsMatch(t *testing.T) {
	fs := NewNDFlagSet("NDflag_test", flag.ExitOnError)
	var (
		bv    *bool
		iv    *int
		i64v  *int64
		uiv   *uint
		ui64v *uint64
		sv    *string
		fv    *float64
		dv    *time.Duration
	)
	fs.NDBoolVar(&bv, "test_bool", "bool value")
	fs.NDIntVar(&iv, "test_int", "int value")
	fs.NDInt64Var(&i64v, "test_int64", "int64 value")
	fs.NDUintVar(&uiv, "test_uint", "uint value")
	fs.NDUint64Var(&ui64v, "test_uint64", "uint64 value")
	fs.NDStringVar(&sv, "test_string", "string value")
	fs.NDFloat64Var(&fv, "test_float64", "float64 value")
	fs.NDDurationVar(&dv, "test_duration", "time.Duration value")

	_ = fs.Set("test_bool", "true")
	_ = fs.Set("test_int", "42")
	_ = fs.Set("test_int64", "-420")
	_ = fs.Set("test_uint", "80")
	_ = fs.Set("test_uint64", "800")
	_ = fs.Set("test_string", "your ad here")
	_ = fs.Set("test_float64", "123.45")
	_ = fs.Set("test_duration", "30s")
	m := make(map[string]*flag.Flag)

	visitor := func(f *flag.Flag) {
		m[f.Name] = f
		var ok bool
		g := f.Value.(flag.Getter)
		switch f.Name {
		case "test_bool":
			ok = g.Get() == bv && *bv == true
		case "test_int":
			ok = g.Get() == iv && *iv == 42
		case "test_int64":
			ok = g.Get() == i64v && *i64v == int64(-420)
		case "test_uint":
			ok = g.Get() == uiv && *uiv == uint(80)
		case "test_uint64":
			ok = g.Get() == ui64v && *ui64v == uint64(800)
		case "test_string":
			ok = g.Get() == sv && *sv == "your ad here"
		case "test_float64":
			ok = g.Get() == fv && *fv == float64(123.45)
		case "test_duration":
			ok = g.Get() == dv && dv.String() == "30s"
		}
		if !ok {
			t.Errorf("Visit: bad value %T(%#v)%p for %s", g.Get(), g.Get(), g.Get(), f.Name)
		}
	}
	fs.Visit(visitor)
	if len(m) != 8 {
		t.Error("Visit misses some flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
}

func TestPtrsMatchDeux(t *testing.T) {
	fs := NewNDFlagSet("NDflag_test", flag.ExitOnError)

	bv := fs.NDBool("test_bool", "bool value")
	iv := fs.NDInt("test_int", "int value")
	i64v := fs.NDInt64("test_int64", "int64 value")
	uiv := fs.NDUint("test_uint", "uint value")
	ui64v := fs.NDUint64("test_uint64", "uint64 value")
	sv := fs.NDString("test_string", "string value")
	fv := fs.NDFloat64("test_float64", "float64 value")
	dv := fs.NDDuration("test_duration", "time.Duration value")

	_ = fs.Set("test_bool", "true")
	_ = fs.Set("test_int", "42")
	_ = fs.Set("test_int64", "-420")
	_ = fs.Set("test_uint", "80")
	_ = fs.Set("test_uint64", "800")
	_ = fs.Set("test_string", "your ad here")
	_ = fs.Set("test_float64", "123.45")
	_ = fs.Set("test_duration", "30s")
	m := make(map[string]*flag.Flag)

	visitor := func(f *flag.Flag) {
		m[f.Name] = f
		var ok bool
		g := f.Value.(flag.Getter)
		switch f.Name {
		case "test_bool":
			ok = g.Get() == *bv && **bv == true
		case "test_int":
			ok = g.Get() == *iv && **iv == 42
		case "test_int64":
			ok = g.Get() == *i64v && **i64v == int64(-420)
		case "test_uint":
			ok = g.Get() == *uiv && **uiv == uint(80)
		case "test_uint64":
			ok = g.Get() == *ui64v && **ui64v == uint64(800)
		case "test_string":
			ok = g.Get() == *sv && **sv == "your ad here"
		case "test_float64":
			ok = g.Get() == *fv && **fv == float64(123.45)
		case "test_duration":
			ok = g.Get() == *dv && (*dv).String() == "30s"
		}
		if !ok {
			t.Errorf("Visit: bad value %T(%#v)%p for %s", g.Get(), g.Get(), g.Get(), f.Name)
		}
	}
	fs.Visit(visitor)
	if len(m) != 8 {
		t.Error("Visit misses some flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
}

func nfs() *NDFlagSet {
	fs := NewNDFlagSet("NDflag_test", flag.ExitOnError)
	fs.NDBool("test_bool", "bool value")
	fs.NDInt("test_int", "int value")
	fs.NDInt64("test_int64", "int64 value")
	fs.NDUint("test_uint", "uint value")
	fs.NDUint64("test_uint64", "uint64 value")
	fs.NDString("test_string", "string value")
	fs.NDFloat64("test_float64", "float64 value")
	fs.NDDuration("test_duration", "time.Duration value")
	return fs
}
