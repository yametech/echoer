package command

import "testing"

func TestFormat(t *testing.T) {
	b := NewFormat().
		Header("id", "name", "age").
		Row("1", "xx", "27").
		Row("2", "yy", "32").
		Out()
	if len(b) < 1 {
		t.Fatal("unknown error")
	}
}
