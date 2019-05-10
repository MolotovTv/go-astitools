package astistring_test

import (
	"testing"

	"github.com/molotovtv/go-astitools/string"
	"github.com/stretchr/testify/assert"
)

func TestToLength(t *testing.T) {
	assert.Equal(t, "test", astistring.ToLength("test", " ", 4))
	assert.Equal(t, "test", astistring.ToLength("testtest", " ", 4))
	assert.Equal(t, "test  ", astistring.ToLength("test", " ", 6))
}
