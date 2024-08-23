package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetTag2Value(t *testing.T) {
	type TestStruct struct {
		A int     `quote:"aq"`
		B string  `quote:"bq"`
		C float64 `quote:"cq"`
		D float32 `quote:""`
	}

	s := &TestStruct{
		A: 1,
		B: "2",
		C: 3,
		D: 4,
	}

	tag2Value := GetTag2Value(s, "quote")

	assert.Condition(t, func() bool {
		if tag2Value["aq"] != int(1) {
			return false
		}
		if tag2Value["bq"] != "2" {
			return false
		}
		if tag2Value["cq"] != float64(3) {
			return false
		}
		if tag2Value[""] != nil {
			return false
		}
		return true
	})
}
