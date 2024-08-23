package excel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// int
func Test_str2int(t *testing.T) {
	i, err := str2int("1", 0, 0)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, int(1), i) {
		return
	}
	_, err = str2int("1a", 0, 0)
	if !assert.Error(t, err) {
		return
	}
}

// int8
func Test_str2int8(t *testing.T) {
	i, err := str2int8("1", 0, 0)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, int8(1), i) {
		return
	}
	_, err = str2int8("1a", 0, 0)
	if !assert.Error(t, err) {
		return
	}
}

// int16
func Test_str2int16(t *testing.T) {
	i, err := str2int16("1", 0, 0)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, int16(1), i) {
		return
	}
	_, err = str2int16("1a", 0, 0)
	if !assert.Error(t, err) {
		return
	}
}

// int32
func Test_str2int32(t *testing.T) {
	i, err := str2int32("1", 0, 0)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, int32(1), i) {
		return
	}
	_, err = str2int32("1a", 0, 0)
	if !assert.Error(t, err) {
		return
	}
}

// int64
func Test_str2int64(t *testing.T) {
	i, err := str2int64("1", 0, 0)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, int64(1), i) {
		return
	}
	_, err = str2int64("1a", 0, 0)
	if !assert.Error(t, err) {
		return
	}
}

func Test_str2uint8(t *testing.T) {
	// uint8
	i, err := str2uint8("1", 0, 0)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, uint8(1), i) {
		return
	}
	_, err = str2uint8("1a", 0, 0)
	if !assert.Error(t, err) {
		return
	}
}

// uint16
func Test_str2uint16(t *testing.T) {
	i, err := str2uint16("1", 0, 0)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, uint16(1), i) {
		return
	}
	_, err = str2uint16("1a", 0, 0)
	if !assert.Error(t, err) {
		return
	}
}

// uint32
func Test_str2uint32(t *testing.T) {
	i, err := str2uint32("1", 0, 0)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, uint32(1), i) {
		return
	}
	_, err = str2uint32("1a", 0, 0)
	if !assert.Error(t, err) {
		return
	}
}

// uint64
func Test_str2uint64(t *testing.T) {
	i, err := str2uint64("1", 0, 0)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, uint64(1), i) {
		return
	}
	_, err = str2uint64("1a", 0, 0)
	if !assert.Error(t, err) {
		return
	}
}

// float32
func Test_str2float32(t *testing.T) {
	f, err := str2float32("1.1", 0, 0)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, float32(1.1), f) {
		return
	}
	_, err = str2float32("1a", 0, 0)
	if !assert.Error(t, err) {
		return
	}
}

// float64
func Test_str2float64(t *testing.T) {
	f, err := str2float64("1.2", 0, 0)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, float64(1.2), f) {
		return
	}
	_, err = str2float64("1a", 0, 0)
	if !assert.Error(t, err) {
		return
	}
}

// bool
func Test_str2bool(t *testing.T) {
	r, err := str2bool("true", 0, 0)
	if !assert.NoError(t, err) {
		return
	}
	b, ok := r.(bool)
	if !assert.True(t, ok) {
		return
	}
	if !assert.True(t, b) {
		return
	}

	r, err = str2bool("false", 0, 0)
	if !assert.NoError(t, err) {
		return
	}
	b, ok = r.(bool)
	if !assert.True(t, ok) {
		return
	}
	if !assert.False(t, b) {
		return
	}

	_, err = str2bool("hello", 0, 0)
	assert.Error(t, err)
}
