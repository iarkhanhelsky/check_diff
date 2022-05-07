package command_test

import (
	"bytes"
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion_Run(t *testing.T) {

	buf := &bytes.Buffer{}
	env := command.Env{
		Version:   "1.0.0",
		Commit:    "deadbeef",
		Date:      "2000-10-10",
		OutWriter: buf,
	}

	assert := assert.New(t)
	assert.NoError(command.NewVersion(env).Run())
	assert.Equal("check_diff v1.0.0 (deadbeef, 2000-10-10)\n", buf.String())
}
