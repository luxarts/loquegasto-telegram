package controller

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandsController_formatFloat(t *testing.T) {
	n := 123456789.123
	expected := "123.456.789,12"

	actual := formatFloat(n)

	assert.Equal(t, expected, actual)
}
