package complete

import (
	"flag"
	"fmt"
	"strconv"
	"testing"
)

func TestFlags(t *testing.T) {
	t.Parallel()

	var (
		tr boolValue = true
		fl boolValue = false
	)

	fs := flag.NewFlagSet("test", flag.ExitOnError)
	fs.Var(&tr, "foo", "")
	fs.Var(&fl, "bar", "")
	fs.String("foo-bar", "", "")
	cmp := FlagSet(fs)

	Test(t, cmp, "", []string{"-foo", "-bar", "-foo-bar", "-h"})
	Test(t, cmp, "-foo", []string{"-foo", "-foo-bar"})
	Test(t, cmp, "-foo ", []string{"false"})
	Test(t, cmp, "-foo=", []string{"false"})
	Test(t, cmp, "-bar ", []string{"-foo", "-bar", "-foo-bar", "-h"})
	Test(t, cmp, "-bar=", []string{})
}

type boolValue bool

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("bad value %q for bool flag", s)
	}
	*b = boolValue(v)
	return nil
}

func (b *boolValue) Get() interface{} { return bool(*b) }

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }

func (b *boolValue) Predict(_ string) []string {
	// If false, typing the bool flag is expected to turn it on, so there is nothing to complete
	// after the flag.
	if *b == false {
		return nil
	}
	// Otherwise, suggest only to turn it off.
	return []string{"false"}
}
