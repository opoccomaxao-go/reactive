package reactive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat64_BatchUpdate(t *testing.T) {
	f := NewFloat64()
	f.BatchSet(map[string]float64{
		"test":  1,
		"test2": 2,
	})
	f.BatchUpdate(func(m Float64Map) Float64Map {
		m["test4"] = 4
		return Float64Map{
			"test3": 3,
			"test":  0,
		}
	})

	assert.Equal(t, float64(0), f.Get("test4"))
	assert.Equal(t, float64(3), f.Get("test3"))
	assert.Equal(t, float64(2), f.Get("test2"))
	assert.Equal(t, float64(0), f.Get("test"))
}
