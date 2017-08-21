package nodefflag

import (
	"flag"
	"testing"
	"time"
)

func TestEverything(t *testing.T) {
	fs := NewNoDefFlagSet("nodefflag_test", flag.ExitOnError)
	fs.NoDefBool("test_bool", "bool value")
	fs.NoDefInt("test_int", "int value")
	fs.NoDefInt64("test_int64", "int64 value")
	fs.NoDefUint("test_uint", "uint value")
	fs.NoDefUint64("test_uint64", "uint64 value")
	fs.NoDefString("test_string", "string value")
	fs.NoDefFloat64("test_float64", "float64 value")
	fs.NoDefDuration("test_duration", "time.Duration value")

	visitor := func(f *flag.Flag) {
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

	_ = fs.Set("test_bool", "false")
	_ = fs.Set("test_int", "42")
	_ = fs.Set("test_int64", "-420")
	_ = fs.Set("test_uint", "80")
	_ = fs.Set("test_uint64", "800")
	_ = fs.Set("test_string", "your ad here")
	_ = fs.Set("test_float64", "123.45")
	_ = fs.Set("test_duration", "30s")

	visitor = func(f *flag.Flag) {
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
}
