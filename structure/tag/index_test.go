package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetTagIndex(t *testing.T) {
	type TestStruct struct {
		A int     `quote:"aq"`
		B string  `quote:"bq"`
		C float64 `quote:"cq"`
		D float32 `quote:""`
	}

	s := &TestStruct{}

	tagIndex := GetTagIndex(s, "quote")

	assert.Condition(t, func() bool {
		if tagIndex["aq"] != 0 {
			return false
		}
		if tagIndex["bq"] != 1 {
			return false
		}
		if tagIndex["cq"] != 2 {
			return false
		}
		if _, ok := tagIndex[""]; ok {
			return false
		}
		return true
	})
}
