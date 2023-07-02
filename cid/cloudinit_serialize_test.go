package cid

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_strParse(t *testing.T) {
	tests := []struct {
		name string
		args string
		want any
	}{
		{"1", "tRuE", true},
		{"2", "FalsE", false},
		{"3", ".", "."},
		{"4", "[]", []any{}},
		{"5", "[.]", []any{"."}},
		{"6", "[50, 50]", []any{int64(50), int64(50)}},
		{"7", "[[10, 20]]", []any{[]any{int64(10), int64(20)}}},
		{"8", "[aaaa]", []any{"aaaa"}},
		{"9", "[[]]", []any{[]any{}}},
		{"10", "[[]", []any{"["}},
		{"11", "[]]", []any{"]"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strParse(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				spew.Dump(got, tt.want)
				t.Errorf("strParse() = %v, want %v", got, tt.want)
			}
		})
	}
}
