package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrorTrapShouldCreateNewInstane(t *testing.T) {
	errt := NewErrorTrap()

	assert.NotNil(t, errt)
	assert.Nil(t, errt.Err())
}

func TestNewErrorTrapShouldTrapFirstErrorOccurence(t *testing.T) {
	err1 := errors.New("pimt: first error")
	err2 := errors.New("pimit: second error")

	errt := NewErrorTrap()

	assert.NotNil(t, errt)
	assert.Nil(t, errt.Err())
	assert.False(t, errt.IsSet())

	errt.Set(err1)

	assert.True(t, errt.IsSet())

	errt.Set(err2)

	assert.True(t, errt.IsSet())
	assert.Equal(t, err1, errt.Err())
}
