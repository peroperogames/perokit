package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUuid(t *testing.T) {
	assert.Equal(t, 36, len(NewUuid()))
}
