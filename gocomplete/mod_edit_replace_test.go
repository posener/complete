package main

import (
	"reflect"
	"testing"
)

func Test_parseReplace(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want replaceArgs
	}{
		// {"1", "-replace=github.com/wxio/asdf@v1.1.1=gothub.com/wxio/asfas@master", replaceArgs{"github.com/wxio/asdf", "v1.1.1", true, "gothub.com/wxio/asfas", "master"}},
		// {"2", "-replace=github.com/wxio/asdf@v1.1.1=gothub.com/wxio/asfas", replaceArgs{"github.com/wxio/asdf", "v1.1.1", true, "gothub.com/wxio/asfas", ""}},
		// {"2.1", "-replace=github.com/wxio/asdf=gothub.com/wxio/asfas", replaceArgs{"github.com/wxio/asdf", "", true, "gothub.com/wxio/asfas", ""}},
		// {"2.2", "-replace=github.com/wxio/asdf=gothub.com/wxio/asfas@v1.2.1", replaceArgs{"github.com/wxio/asdf", "", true, "gothub.com/wxio/asfas", "v1.2.1"}},
		// {"3", "-replace=github.com/wxio/asdf@v1.1.1", replaceArgs{"github.com/wxio/asdf", "v1.1.1", false, "", ""}},
		// {"4", "-replace=github.com/wxio/asdf", replaceArgs{"github.com/wxio/asdf", "", false, "", ""}},
		// {"5", "-replace=", replaceArgs{"", "", false, "", ""}},
		// {"0", "-replace=github.com/wxio/asdf@v1.1.1gothub.com/wxio/asfas@master", replaceArgs{"github.com/wxio/asdf", "v1.1.1gothub.com/wxio/asfas@master", false, "", ""}},

		{"1", "github.com/wxio/asdf@v1.1.1=gothub.com/wxio/asfas@master", replaceArgs{"github.com/wxio/asdf", true, "v1.1.1", true, "gothub.com/wxio/asfas", true, "master"}},
		{"2", "github.com/wxio/asdf@v1.1.1=gothub.com/wxio/asfas", replaceArgs{"github.com/wxio/asdf", true, "v1.1.1", true, "gothub.com/wxio/asfas", false, ""}},
		{"2.1", "github.com/wxio/asdf=gothub.com/wxio/asfas", replaceArgs{"github.com/wxio/asdf", false, "", true, "gothub.com/wxio/asfas", false, ""}},
		{"2.1.1", "github.com/posener/complete=github.com/posener/complete", replaceArgs{"github.com/posener/complete", false, "", true, "github.com/posener/complete", false, ""}},
		{"2.2", "github.com/wxio/asdf=gothub.com/wxio/asfas@v1.2.1", replaceArgs{"github.com/wxio/asdf", false, "", true, "gothub.com/wxio/asfas", true, "v1.2.1"}},
		{"3", "github.com/wxio/asdf@v1.1.1", replaceArgs{"github.com/wxio/asdf", true, "v1.1.1", false, "", false, ""}},
		{"4", "github.com/wxio/asdf", replaceArgs{"github.com/wxio/asdf", false, "", false, "", false, ""}},
		{"5", "", replaceArgs{"", false, "", false, "", false, ""}},
		{"0", "github.com/wxio/asdf@v1.1.1gothub.com/wxio/asfas@master", replaceArgs{"github.com/wxio/asdf", true, "v1.1.1gothub.com/wxio/asfas@master", false, "", false, ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseReplace(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseReplace() = \n%#+v, want \n%#+v", got, tt.want)
			}
		})
	}
}
