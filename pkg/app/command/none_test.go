package command_test

import (
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNone_Run(t *testing.T) {
	assert := assert.New(t)
	assert.NoError((&command.None{}).Run())
}
