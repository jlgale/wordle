package wordle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func mk(s string) Letters {
	return NewLetters([]byte(s))
}

var all Letters = mk("abcdefghijklmnopqrstuvwxyz")

func TestAdd(t *testing.T) {
	a := mk("abc")
	b := mk("xyz")
	aa := mk("abcd")
	e := mk("")
	assert.Equal(t, mk("abcxyz"), a.Add(b))
	assert.Equal(t, mk(""), e.Add(e))
	assert.Equal(t, mk("abc"), a.Add(a))
	assert.Equal(t, mk("abcd"), a.Add(aa))
}

func TestRemove(t *testing.T) {
	a := mk("abcd")
	e := mk("")
	assert.Equal(t, a, a.Remove(e))
	assert.Equal(t, e, a.Remove(a))
	assert.Equal(t, mk("ad"), a.Remove(mk("bc")))
}

func TestIntersect(t *testing.T) {
	a := mk("abcd")
	b := mk("efg")
	assert.Equal(t, a, a.Intersect(a))
	assert.Equal(t, mk(""), a.Intersect(b))
}

func TestString(t *testing.T) {
	assert.Equal(t, "[a-d]", mk("abcd").String())
	assert.Equal(t, "[acd]", mk("acd").String())
	assert.Equal(t, "[]", mk("").String())
	assert.Equal(t, "[a-z]", all.String())
	assert.Equal(t, "[a-cp-t]", mk("abc").Add(mk("p")).Add(mk("qrst")).String())
}

func TestContains(t *testing.T) {
	assert.Equal(t, true, mk("tuv").Contains('t'))
	assert.Equal(t, true, mk("tuv").Contains('u'))
	assert.Equal(t, false, mk("abcefg").Contains('d'))
	assert.Equal(t, false, mk("").Contains('t'))

	assert.Equal(t, true, mk("").Empty())
	assert.Equal(t, false, mk("a").Empty())
}
