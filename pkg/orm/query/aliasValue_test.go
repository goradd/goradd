package query

import (
	"fmt"
	"github.com/goradd/goradd/pkg/datetime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAliasValue(t *testing.T) {
	tests := []struct {
		name  interface{}
		value interface{}
		want  interface{}
	}{
		{"abc", NewAliasValue("abc").String(), "abc"},
		{5, NewAliasValue("5").Int(), 5},
		{1.23, NewAliasValue("1.23").Float(), 1.23},
		{true, NewAliasValue("true").Bool(), true},
		{false, NewAliasValue("false").Bool(), false},
		{"IsNil", NewAliasValue(nil).IsNil(), true},
		{"IsNull", NewAliasValue(nil).IsNull(), true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.name), func(t *testing.T) {
			if got := tt.value; got != tt.want {
				t.Errorf("NewAliasValue(%#v), got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestAliasValue_DateTime(t *testing.T) {
	assert.True(t, NewAliasValue("2001-05-04").DateTime().Equal(datetime.NewDateTime("2001-05-04T0:00:00Z")))
}
