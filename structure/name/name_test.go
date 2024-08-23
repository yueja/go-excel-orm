package name

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetName(t *testing.T) {
	type testStruct struct {
		A int
	}
	expected := `github.com/yueja/go-excel-orm/structure/name/testStruct`

	s := &testStruct{}
	name := GetFullTypeName(s)

	assert.Equal(t, expected, name)
}
