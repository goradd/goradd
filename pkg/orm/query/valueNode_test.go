package query_test

import (
	"github.com/goradd/goradd/pkg/orm/query"
	"github.com/stretchr/testify/assert"
	"testing"
)

type pretendInt int16
type pretendString string
type pretendBool bool
type pretendFloat float32

func TestValueSetting(t *testing.T) {
	var p pretendInt = 2
	var s pretendString = "b"
	var u uint8 = 3
	var b pretendBool = true
	var f pretendFloat = 3.1

	assert.True(t, query.Value(p).Equals(query.Value(2)))
	assert.True(t, query.Value(2).Equals(query.Value(2)))
	assert.True(t, query.Value(2.2).Equals(query.Value(2.2)))
	assert.True(t, query.Value(u).Equals(query.Value(u)))
	assert.True(t, query.Value("a").Equals(query.Value("a")))
	assert.True(t, query.Value(s).Equals(query.Value("b")))
	assert.True(t, query.Value(f).Equals(query.Value(float32(3.1))))
	assert.True(t, query.Value(b).Equals(query.Value(true)))

	assert.False(t, query.Value(u).Equals(query.Value(2.2)))

}

func TestValueArray(t *testing.T) {
	a := []int{1,2}
	b := []int{1,3}
	c := []string{"a", "b"}
	d := 5
	e := []string{"c"}

	assert.True(t, query.Value(a).Equals(query.Value(a)), "Array node equal to self.")
	assert.False(t, query.Value(a).Equals(query.Value(b)))
	assert.False(t, query.Value(a).Equals(query.Value(c)))
	assert.False(t, query.Value(a).Equals(query.Value(d)))
	assert.False(t, query.Value(a).Equals(query.Value(e)))
	assert.False(t, query.Value(d).Equals(query.Value(a)))
}