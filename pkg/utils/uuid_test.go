package utils

import (
	"fmt"
	"testing"
)

func TestSUID(t *testing.T) {
	u := NewSUID()

	fmt.Printf("[%s]\r\n[%s]\r\n", u.StringFull(), u.String())
}
