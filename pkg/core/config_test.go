package core

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDefaultConfig(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(Config{OutputFormat: "stdout"}, NewDefaultConfig())
}
