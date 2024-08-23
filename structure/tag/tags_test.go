package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetTags(t *testing.T) {
	type TestStruct struct {
		A int    `q:"a"`
		B bool   `q:"b"`
		C string `q:"c"`
	}

	expected := []string{"a", "b", "c"}
	tags := GetTags(TestStruct{}, "q")

	assert.Equal(t, expected, tags)
}
