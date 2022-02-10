package wego

import (
	"testing"
)

func TestSplitString(t *testing.T) {
	data := ""
	key, val , has := SplitString(data, ";")
	t.Logf("str=%s, key=%s, val=%s, has=%v", data, key, val, has)

	data = "aaa"
	key, val , has = SplitString(data, ";")
	t.Logf("str=%s, key=%s, val=%s, has=%v", data, key, val, has)

	data = "aaa;"
	key, val , has = SplitString(data, ";")
	t.Logf("str=%s, key=%s, val=%s, has=%v", data, key, val, has)

	data = ";aaa"
	key, val , has = SplitString(data, ";")
	t.Logf("str=%s, key=%s, val=%s, has=%v", data, key, val, has)

	data = " ;aaa"
	key, val , has = SplitString(data, ";")
	t.Logf("str=%s, key=%s, val=%s, has=%v", data, key, val, has)

	data = "aaa; "
	key, val , has = SplitString(data, ";")
	t.Logf("str=%s, key=%s, val=%s, has=%v", data, key, val, has)

	data = "aaa;bbb"
	key, val , has = SplitString(data, ";")
	t.Logf("str=%s, key=%s, val=%s, has=%v", data, key, val, has)

	data = "aaa ;bbb"
	key, val , has = SplitString(data, ";")
	t.Logf("str=%s, key=%s, val=%s, has=%v", data, key, val, has)

	data = "aaa; bbb"
	key, val , has = SplitString(data, ";")
	t.Logf("str=%s, key=%s, val=%s, has=%v", data, key, val, has)

	data = "aaa;bbb;ccc"
	key, val , has = SplitString(data, ";")
	t.Logf("str=%s, key=%s, val=%s, has=%v", data, key, val, has)
}
