package compflag

import (
	"flag"
	"testing"
	"time"

	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	t.Parallel()

	t.Run("complete default off", func(t *testing.T) {
		var cmd FlagSet
		_ = cmd.Bool("a", false, "")
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a ", []string{"-a", "-h"})
	})

	t.Run("complete default on", func(t *testing.T) {
		var cmd FlagSet
		_ = cmd.Bool("a", true, "")
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a ", []string{"false"})
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a=", []string{"false"})
	})
}

func TestString(t *testing.T) {
	t.Parallel()

	t.Run("options invalid not checked", func(t *testing.T) {
		var cmd FlagSet
		value := cmd.String("a", "", "", predict.OptValues("1", "2"))
		err := cmd.Parse([]string{"-a", "3"})
		assert.NoError(t, err)
		assert.Equal(t, "3", *value)
	})

	t.Run("options valid checked", func(t *testing.T) {
		var cmd FlagSet
		value := cmd.String("a", "", "", predict.OptValues("1", "2"), predict.OptCheck())
		err := cmd.Parse([]string{"-a", "2"})
		assert.NoError(t, err)
		assert.Equal(t, "2", *value)
	})

	t.Run("options invalid checked", func(t *testing.T) {
		var cmd FlagSet
		_ = cmd.String("a", "", "", predict.OptValues("1", "2"), predict.OptCheck())
		err := cmd.Parse([]string{"-a", "3"})
		assert.Error(t, err)
	})

	t.Run("complete", func(t *testing.T) {
		var cmd FlagSet
		_ = cmd.String("a", "", "", predict.OptValues("1", "2"))
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a ", []string{"1", "2"})
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a=", []string{"1", "2"})
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a 1", []string{"1"})
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a=1", []string{"1"})
	})
}

func TestInt(t *testing.T) {
	t.Parallel()

	t.Run("options invalid not checked", func(t *testing.T) {
		var cmd FlagSet
		value := cmd.Int("a", 0, "", predict.OptValues("1", "2"))
		err := cmd.Parse([]string{"-a", "3"})
		assert.NoError(t, err)
		assert.Equal(t, 3, *value)
	})

	t.Run("options valid checked", func(t *testing.T) {
		var cmd FlagSet
		value := cmd.Int("a", 0, "", predict.OptValues("1", "2"), predict.OptCheck())
		err := cmd.Parse([]string{"-a", "2"})
		assert.NoError(t, err)
		assert.Equal(t, 2, *value)
	})

	t.Run("options invalid checked", func(t *testing.T) {
		var cmd FlagSet
		_ = cmd.Int("a", 0, "", predict.OptValues("1", "2"), predict.OptCheck())
		err := cmd.Parse([]string{"-a", "3"})
		assert.Error(t, err)
	})

	t.Run("options invalid int value", func(t *testing.T) {
		var cmd FlagSet
		_ = cmd.Int("a", 0, "", predict.OptValues("1", "2", "x"), predict.OptCheck())
		err := cmd.Parse([]string{"-a", "x"})
		assert.Error(t, err)
	})

	t.Run("complete", func(t *testing.T) {
		var cmd FlagSet
		_ = cmd.Int("a", 0, "", predict.OptValues("1", "2"))
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a ", []string{"1", "2"})
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a=", []string{"1", "2"})
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a 1", []string{"1"})
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a=1", []string{"1"})
	})
}

func TestDuration(t *testing.T) {
	t.Parallel()

	t.Run("options invalid not checked", func(t *testing.T) {
		var cmd FlagSet
		value := cmd.Duration("a", 0, "", predict.OptValues("1s", "1m"))
		err := cmd.Parse([]string{"-a", "1h"})
		assert.NoError(t, err)
		assert.Equal(t, time.Hour, *value)
	})

	t.Run("options valid checked", func(t *testing.T) {
		var cmd FlagSet
		value := cmd.Duration("a", 0, "", predict.OptValues("1s", "1m"), predict.OptCheck())
		err := cmd.Parse([]string{"-a", "1m"})
		assert.NoError(t, err)
		assert.Equal(t, time.Minute, *value)
	})

	t.Run("options invalid checked", func(t *testing.T) {
		var cmd FlagSet
		_ = cmd.Duration("a", 0, "", predict.OptValues("1s", "1m"), predict.OptCheck())
		err := cmd.Parse([]string{"-a", "1h"})
		assert.Error(t, err)
	})

	t.Run("options invalid duration value", func(t *testing.T) {
		var cmd FlagSet
		_ = cmd.Duration("a", 0, "", predict.OptValues("1h", "1m", "1"), predict.OptCheck())
		err := cmd.Parse([]string{"-a", "1"})
		assert.Error(t, err)
	})

	t.Run("complete", func(t *testing.T) {
		var cmd FlagSet
		_ = cmd.Duration("a", 0, "", predict.OptValues("1s", "1m"))
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a ", []string{"1s", "1m"})
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a=", []string{"1s", "1m"})
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a 1", []string{"1s", "1m"})
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a=1", []string{"1s", "1m"})
		complete.Test(t, complete.FlagSet((*flag.FlagSet)(&cmd)), "-a=1m", []string{"1m"})
	})
}
