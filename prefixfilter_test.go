package complete

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPrefixFilteringPredictor_Predict(t *testing.T) {
	t.Parallel()
	initTests()

	t.Run("defaults to empty list", func(t *testing.T) {
		pfp := &PrefixFilteringPredictor{}
		got := pfp.Predict(Args{})
		if len(got) != 0 {
			t.Fail()
		}
	})

	t.Run("passes request to Predictor", func(t *testing.T) {
		args := Args{
			All: []string{"a"},
		}
		want := []string{"b"}
		predictFunc := PredictFunc(func(a Args) []string {
			if !reflect.DeepEqual(a, args) {
				t.Errorf("unexpected args: %v", a)
			}
			return want
		})
		pfp := &PrefixFilteringPredictor{
			Predictor: predictFunc,
		}
		got := pfp.Predict(args)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("unexpected result: %v", got)
		}
	})
}

func TestPrefixFilteringPredictor_FilterPrefix(t *testing.T) {
	t.Parallel()
	initTests()

	t.Run("default PrefixFilterFunc", func(t *testing.T) {
		for _, td := range []struct {
			s      string
			prefix string
			want   bool
		}{
			{
				s:      "ohm",
				prefix: "ohm",
				want:   true,
			},
			{
				s:      "ohm",
				prefix: "",
				want:   true,
			},
			{
				s:      "ohm",
				prefix: "O",
				want:   false,
			},
			{
				s:      "ohm",
				prefix: "q",
				want:   false,
			},
			{
				s:      "öhm",
				prefix: "o",
				want:   false,
			},
			{
				s:      "ohm",
				prefix: "ohmy",
				want:   false,
			},
		} {
			t.Run(fmt.Sprintf("%s %s", td.s, td.prefix), func(t *testing.T) {
				pfp := &PrefixFilteringPredictor{}
				got := pfp.FilterPrefix(td.s, td.prefix)
				if td.want != got {
					t.Errorf("failed %s\ngot: %v\nwant: %v", t.Name(), got, td.want)
				}
			})
		}
	})

	t.Run("CaseInsensitivePrefixFilter", func(t *testing.T) {
		for _, td := range []struct {
			s      string
			prefix string
			want   bool
		}{
			{
				s:      "ohm",
				prefix: "ohm",
				want:   true,
			},
			{
				s:      "ohm",
				prefix: "",
				want:   true,
			},
			{
				s:      "ohm",
				prefix: "O",
				want:   true,
			},
			{
				s:      "ohm",
				prefix: "q",
				want:   false,
			},
			{
				s:      "öhm",
				prefix: "o",
				want:   false,
			},
			{
				s:      "ohm",
				prefix: "ohmy",
				want:   false,
			},
		} {
			t.Run(fmt.Sprintf("%s %s", td.s, td.prefix), func(t *testing.T) {
				pfp := &PrefixFilteringPredictor{
					PrefixFilterFunc: CaseInsensitivePrefixFilter,
				}
				got := pfp.FilterPrefix(td.s, td.prefix)
				if td.want != got {
					t.Errorf("failed %s\ngot: %v\nwant: %v", t.Name(), got, td.want)
				}
			})
		}
	})

	t.Run("PermissivePrefixFilter", func(t *testing.T) {
		pfp := &PrefixFilteringPredictor{
			PrefixFilterFunc: PermissivePrefixFilter,
		}
		got := pfp.FilterPrefix("", "")
		if !got {
			t.Errorf("should have returned true, but didn't")
		}
	})
}
